package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	timestampRe      = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}(?:\s+|$)`)
	commandMarkerRe  = regexp.MustCompile(`^CUSTOM KUBX KUBECTL COMMAND:\s*(.*?)\s*$`)
	separatorRe      = regexp.MustCompile(`^-{20,}$|^={20,}$`)
	kubeconfigHintRe = regexp.MustCompile(`^Found kubeconfig in .+ secret$`)
)

type commandBlock struct {
	Command string
	Output  string
}

type commandResult struct {
	Command string `json:"command"`
	Output  string `json:"output"`
}

func stripTimestamp(line string) string {
	return timestampRe.ReplaceAllString(strings.TrimRight(line, "\n"), "")
}

func trimBlankLines(lines []string) []string {
	start, end := 0, len(lines)
	for start < end && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	for end > start && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}
	return lines[start:end]
}

// parseConsoleOutput extracts one output block per echoed command from the
// Jenkins consoleText: markers open a block, separators close finished ones,
// timestamps/kubeconfig hints are noise and get dropped.
func parseConsoleOutput(rawConsole string) []commandBlock {
	var (
		blocks         []commandBlock
		currentCommand string
		currentLines   []string
		haveCurrent    bool
	)

	flush := func() {
		if !haveCurrent {
			currentLines = nil
			return
		}
		blocks = append(blocks, commandBlock{
			Command: currentCommand,
			Output:  strings.Join(trimBlankLines(currentLines), "\n"),
		})
		currentCommand = ""
		currentLines = nil
		haveCurrent = false
	}

	for _, rawLine := range strings.Split(rawConsole, "\n") {
		line := stripTimestamp(rawLine)
		stripped := strings.TrimSpace(line)

		if marker := commandMarkerRe.FindStringSubmatch(stripped); marker != nil {
			flush()
			currentCommand = strings.TrimSpace(marker[1])
			currentLines = nil
			haveCurrent = true
			continue
		}
		if !haveCurrent {
			continue
		}
		if stripped == "CUSTOM KUBX KUBECTL:" {
			continue
		}
		if separatorRe.MatchString(stripped) {
			if len(currentLines) > 0 {
				flush()
			}
			continue
		}
		if kubeconfigHintRe.MatchString(stripped) {
			continue
		}
		currentLines = append(currentLines, strings.TrimRightFunc(line, unicode.IsSpace))
	}
	flush()
	return blocks
}

// pairResultsWithRequestedCommands aligns parsed console blocks with the
// requested commands. Any mismatch is reported as a warning; the caller turns
// non-empty warnings into an error.
func pairResultsWithRequestedCommands(requested []string, parsed []commandBlock) ([]commandResult, []string, error) {
	var (
		results  []commandResult
		warnings []string
	)
	if len(requested) != len(parsed) {
		warnings = append(warnings, fmt.Sprintf(
			"requested command count does not match parsed Jenkins console blocks (%d requested, %d parsed)",
			len(requested), len(parsed)))
	}

	for i, requestedCommand := range requested {
		var block *commandBlock
		if i < len(parsed) {
			block = &parsed[i]
		}
		parsedCommand := requestedCommand
		output := ""
		if block != nil {
			parsedCommand = block.Command
			output = block.Output
			want, err := normalizeCommand(requestedCommand)
			if err != nil {
				return nil, nil, err
			}
			got, err := normalizeCommand(parsedCommand)
			if err != nil {
				return nil, nil, err
			}
			if want != got {
				warnings = append(warnings, fmt.Sprintf(
					"command mismatch at position %d: requested %q, parsed %q",
					i+1, requestedCommand, parsedCommand))
			}
		}
		results = append(results, commandResult{Command: requestedCommand, Output: output})
	}

	for i := len(requested); i < len(parsed); i++ {
		warnings = append(warnings, fmt.Sprintf("extra command block found in Jenkins console: %q", parsed[i].Command))
		results = append(results, commandResult{Command: parsed[i].Command, Output: parsed[i].Output})
	}
	return results, warnings, nil
}
