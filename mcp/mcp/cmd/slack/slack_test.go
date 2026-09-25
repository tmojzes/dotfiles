package main

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"
)

// fakeReplySource scripts one fetchReplies response per poll; the last
// response repeats once the script is exhausted.
type fakeReplySource struct {
	polls [][]slackMessage
	calls int
}

func (f *fakeReplySource) fetchReplies(context.Context, string) ([]slackMessage, error) {
	if f.calls >= len(f.polls) {
		return f.polls[len(f.polls)-1], nil
	}
	result := f.polls[f.calls]
	f.calls++
	return result, nil
}

func fastConfig() serverConfig {
	return serverConfig{
		pollInterval:   time.Millisecond,
		quietWindow:    5 * time.Millisecond,
		replyTimeout:   time.Second,
		requestTimeout: time.Second,
	}
}

const (
	testThreadRootTS = "1000.000001" // the ':thread:debug satellite' root
	testPostTS       = "1000.000100" // our own post in the thread
)

func TestCollectRepliesQuietWindow(t *testing.T) {
	source := &fakeReplySource{polls: [][]slackMessage{
		{{Ts: testThreadRootTS, Text: ":thread:debug satellite"}}, // first poll: root only
		{
			{Ts: testThreadRootTS, Text: "ignored"},
			{Ts: testPostTS, Text: "ignored"},
			{Ts: "1000.000200", Text: "first"}, // first reply arrives
		},
		{
			{Ts: testThreadRootTS, Text: "ignored"},
			{Ts: testPostTS, Text: "ignored"},
			{Ts: "1000.000200", Text: "first"}, // no change: quiet window starts
		},
		{
			{Ts: testThreadRootTS, Text: "ignored"},
			{Ts: testPostTS, Text: "ignored"},
			{Ts: "1000.000200", Text: "first"}, // still quiet: done
		},
	}}

	replies, timedOut, err := collectReplies(context.Background(), source, testThreadRootTS, testPostTS, fastConfig())
	if err != nil {
		t.Fatalf("collectReplies: %v", err)
	}
	if timedOut {
		t.Fatal("collectReplies timed out before the quiet window")
	}
	want := []slackMessage{{Ts: "1000.000200", Text: "first"}}
	if !reflect.DeepEqual(replies, want) {
		t.Fatalf("replies = %#v, want %#v", replies, want)
	}
}

func TestCollectRepliesIgnoresOlderThreadReplies(t *testing.T) {
	// The debug thread holds replies to earlier commands; only messages
	// newer than our own post may be collected.
	source := &fakeReplySource{polls: [][]slackMessage{
		{
			{Ts: "1000.000002", Text: "old reply to an earlier command"}, // before our post: ignored
			{Ts: "1000.000050", Text: "another old reply"},               // before our post: ignored
			{Ts: testPostTS, Text: "<@bot> cluster x"},                   // our own post: ignored
			{Ts: "1000.000300", Text: "fresh bot reply"},                 // after our post: collected
		},
		{
			{Ts: "1000.000002", Text: "old reply to an earlier command"},
			{Ts: "1000.000050", Text: "another old reply"},
			{Ts: testPostTS, Text: "<@bot> cluster x"},
			{Ts: "1000.000300", Text: "fresh bot reply"},
		},
	}}

	replies, timedOut, err := collectReplies(context.Background(), source, testThreadRootTS, testPostTS, fastConfig())
	if err != nil {
		t.Fatalf("collectReplies: %v", err)
	}
	if timedOut {
		t.Fatal("collectReplies timed out before the quiet window")
	}
	want := []slackMessage{{Ts: "1000.000300", Text: "fresh bot reply"}}
	if !reflect.DeepEqual(replies, want) {
		t.Fatalf("replies = %#v, want %#v", replies, want)
	}
}

func TestCollectRepliesDedupesAcrossPolls(t *testing.T) {
	source := &fakeReplySource{polls: [][]slackMessage{
		{
			{Ts: testThreadRootTS, Text: "root"},
			{Ts: "1000.000200", Text: "a"},
			{Ts: "1000.000300", Text: "b"},
		},
		{
			{Ts: testThreadRootTS, Text: "root"},
			{Ts: "1000.000200", Text: "a"}, // already seen
			{Ts: "1000.000300", Text: "b"}, // already seen
			{Ts: "1000.000400", Text: "c"}, // new: resets the quiet window
		},
		{
			{Ts: testThreadRootTS, Text: "root"},
			{Ts: "1000.000200", Text: "a"},
			{Ts: "1000.000300", Text: "b"},
			{Ts: "1000.000400", Text: "c"},
		},
	}}

	replies, timedOut, err := collectReplies(context.Background(), source, testThreadRootTS, testPostTS, fastConfig())
	if err != nil {
		t.Fatalf("collectReplies: %v", err)
	}
	if timedOut {
		t.Fatal("collectReplies timed out before the quiet window")
	}
	want := []slackMessage{
		{Ts: "1000.000200", Text: "a"},
		{Ts: "1000.000300", Text: "b"},
		{Ts: "1000.000400", Text: "c"},
	}
	if !reflect.DeepEqual(replies, want) {
		t.Fatalf("replies = %#v, want %#v", replies, want)
	}
}

func TestCollectRepliesTimeoutPartial(t *testing.T) {
	source := &fakeReplySource{polls: [][]slackMessage{
		{
			{Ts: testThreadRootTS, Text: "root"},
			{Ts: "1000.000200", Text: "only reply"}, // repeated forever, never goes quiet
		},
	}}
	cfg := fastConfig()
	cfg.quietWindow = time.Hour // can never elapse
	cfg.replyTimeout = 20 * time.Millisecond

	replies, timedOut, err := collectReplies(context.Background(), source, testThreadRootTS, testPostTS, cfg)
	if err != nil {
		t.Fatalf("collectReplies: %v", err)
	}
	if !timedOut {
		t.Fatal("collectReplies should have timed out")
	}
	want := []slackMessage{{Ts: "1000.000200", Text: "only reply"}}
	if !reflect.DeepEqual(replies, want) {
		t.Fatalf("replies = %#v, want %#v", replies, want)
	}
}

func TestCollectRepliesTimeoutNoReplies(t *testing.T) {
	source := &fakeReplySource{polls: [][]slackMessage{
		{{Ts: testThreadRootTS, Text: "root"}}, // no bot reply, ever
	}}
	cfg := fastConfig()
	cfg.replyTimeout = 20 * time.Millisecond

	replies, _, err := collectReplies(context.Background(), source, testThreadRootTS, testPostTS, cfg)
	if err == nil || !strings.Contains(err.Error(), "no bot replies") {
		t.Fatalf("err = %v, want a no-bot-replies error", err)
	}
	if replies != nil {
		t.Fatalf("replies = %#v, want nil", replies)
	}
}

func TestFindLatestThreadRoot(t *testing.T) {
	messages := []slackMessage{
		{Ts: "1790000000.000001", Text: ":thread:debug satellite"},     // older thread with the same title
		{Ts: "1790000100.000001", Text: "<@U08AUQQRPKJ> cluster x"},    // unrelated root
		{Ts: "1790000200.000001", Text: "  :thread:debug satellite  "}, // newest thread (whitespace trimmed)
		{Ts: "1790000150.000001", Text: ":thread:debug satellite"},     // middle thread
	}

	got, err := findLatestThreadRoot(messages, debugThreadTitle)
	if err != nil {
		t.Fatalf("findLatestThreadRoot: %v", err)
	}
	if want := "1790000200.000001"; got != want {
		t.Fatalf("findLatestThreadRoot = %q, want %q", got, want)
	}
}

func TestFindLatestThreadRootMissing(t *testing.T) {
	messages := []slackMessage{{Ts: "1000.000001", Text: "armada-xo started; Region: x"}}

	_, err := findLatestThreadRoot(messages, debugThreadTitle)
	if err == nil || !strings.Contains(err.Error(), `no thread titled ":thread:debug satellite"`) {
		t.Fatalf("err = %v, want a missing-thread error", err)
	}
}

func TestRenderRepliesSortsByTsAndSkipsEmptyText(t *testing.T) {
	replies := []slackMessage{
		{Ts: "1790348854.476249", Text: "Looking for..."}, // Slack returned this one second: sort by ts, not arrival
		{Ts: "1790348854.378499", Text: "Retrieving..."},
		{Ts: "1790348855.316989", Text: "  "}, // whitespace only: skipped
		{Ts: "1790348855.316990", Text: "Region: br-sao ```output```"},
	}

	got := renderReplies(replies, "1790348853.213609")
	want := "Retrieving...\nLooking for...\nRegion: br-sao ```output```\n\n" +
		"Thread: https://ibm.enterprise.slack.com/archives/C53AJ95TP/p1790348853213609"
	if got != want {
		t.Fatalf("renderReplies = %q, want %q", got, want)
	}
}

func TestRenderRepliesAllEmpty(t *testing.T) {
	got := renderReplies([]slackMessage{{Ts: "1000.000002", Text: "  "}}, "1000.000001")
	want := "The bot replied with no readable text. " +
		"Thread: https://ibm.enterprise.slack.com/archives/C53AJ95TP/p1000000001"
	if got != want {
		t.Fatalf("renderReplies = %q, want %q", got, want)
	}
}

func TestRunSlackCommandValidation(t *testing.T) {
	ctx := context.Background()
	cfg := fastConfig()

	cases := []struct {
		name    string
		input   slackCommandInput
		wantErr string
	}{
		{
			name:    "empty command",
			input:   slackCommandInput{Command: "   "},
			wantErr: "command must be a non-empty string",
		},
		{
			name:    "unknown bot",
			input:   slackCommandInput{Command: "cluster x", Bot: "U9999999999"},
			wantErr: "bot must be U08AUQQRPKJ (production, default) or U08DHQA1FG8 (stage)",
		},
		{
			name:    "missing credentials",
			input:   slackCommandInput{Command: "cluster x"},
			wantErr: "SLACK_XOXC and SLACK_XOXD must be set",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("SLACK_XOXC", "")
			t.Setenv("SLACK_XOXD", "")
			_, err := runSlackCommand(ctx, cfg, tc.input)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want it to contain %q", err, tc.wantErr)
			}
		})
	}
}

func TestThreadPermalink(t *testing.T) {
	got := threadPermalink("1790348853.213609")
	want := "https://ibm.enterprise.slack.com/archives/C53AJ95TP/p1790348853213609"
	if got != want {
		t.Fatalf("threadPermalink = %q, want %q", got, want)
	}
}

func TestSlackAPIError(t *testing.T) {
	err := slackAPIError("invalid_auth")
	if err == nil || !strings.Contains(err.Error(), "re-run the slack-token skill") {
		t.Fatalf("err = %v, want a refresh hint", err)
	}
	err = slackAPIError("not_in_channel")
	if err == nil || !strings.Contains(err.Error(), "not a member of #armada-xo") {
		t.Fatalf("err = %v, want a channel membership hint", err)
	}
	err = slackAPIError("")
	if err == nil || !strings.Contains(err.Error(), "empty error code") {
		t.Fatalf("err = %v, want an unknown-error message", err)
	}
}
