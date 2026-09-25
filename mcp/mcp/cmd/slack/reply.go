package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

// findLatestThreadRoot returns the ts of the newest message whose text
// matches the thread title exactly. Several threads share the title over
// time; the latest one is the one new commands are posted under. Slack ts
// values ("seconds.microseconds") sort chronologically as strings.
func findLatestThreadRoot(messages []slackMessage, title string) (string, error) {
	var latest string
	for _, message := range messages {
		if strings.TrimSpace(message.Text) != title {
			continue
		}
		if latest == "" || strings.Compare(message.Ts, latest) > 0 {
			latest = message.Ts
		}
	}
	if latest == "" {
		return "", fmt.Errorf("no thread titled %q found in the recent history of #armada-xo; create the thread in Slack first", title)
	}
	return latest, nil
}

// collectReplies polls the thread rooted at threadTS for replies newer than
// afterTS (the timestamp of our own post — the thread also holds replies to
// earlier commands, which must not be collected) until the source goes
// quiet — no new messages for the configured quiet window, because the
// armada-xo bots reply with a burst of several messages — or the overall
// reply timeout elapses. The poller re-fetches the whole thread every time;
// already-seen replies are filtered out by their ts.
//
// On timeout the replies collected so far are returned with timedOut=true;
// with none at all it is an error, since the command outcome is unknown.
func collectReplies(ctx context.Context, source replySource, threadTS, afterTS string, cfg serverConfig) (replies []slackMessage, timedOut bool, err error) {
	seen := map[string]bool{}
	var lastNew time.Time
	deadline := time.Now().Add(cfg.replyTimeout)

	for {
		messages, err := source.fetchReplies(ctx, threadTS)
		if err != nil {
			return nil, false, err
		}
		added := false
		for _, message := range messages {
			if message.Ts == "" || strings.Compare(message.Ts, afterTS) <= 0 || seen[message.Ts] {
				continue
			}
			seen[message.Ts] = true
			replies = append(replies, message)
			added = true
		}

		now := time.Now()
		if added {
			lastNew = now
		} else if len(replies) > 0 && now.Sub(lastNew) >= cfg.quietWindow {
			return replies, false, nil
		}
		if now.After(deadline) {
			if len(replies) > 0 {
				return replies, true, nil
			}
			return nil, false, errors.New("no bot replies arrived before the reply timeout; check the thread manually")
		}

		select {
		case <-ctx.Done():
			return nil, false, ctx.Err()
		case <-time.After(cfg.pollInterval):
		}
	}
}

// renderReplies formats the collected thread replies as the tool text: each
// message on its own line, in ts order (Slack's ts is the message id and
// sorts chronologically), plus a permalink to the thread for manual
// follow-up. Replies without readable text are skipped.
func renderReplies(replies []slackMessage, threadTS string) string {
	sorted := slices.Clone(replies)
	slices.SortFunc(sorted, func(a, b slackMessage) int {
		return strings.Compare(a.Ts, b.Ts)
	})

	var chunks []string
	for _, reply := range sorted {
		if text := strings.TrimSpace(reply.Text); text != "" {
			chunks = append(chunks, text)
		}
	}

	permalink := threadPermalink(threadTS)
	body := strings.Join(chunks, "\n")
	if body == "" {
		return "The bot replied with no readable text. Thread: " + permalink
	}
	return body + "\n\nThread: " + permalink
}

// threadPermalink builds a web-client permalink for a message in
// #armada-xo. A Slack permalink is /archives/<channel>/p<ts minus the dot>.
func threadPermalink(threadTS string) string {
	return "https://ibm.enterprise.slack.com/archives/" + armadaChannelID + "/p" + strings.ReplaceAll(threadTS, ".", "")
}
