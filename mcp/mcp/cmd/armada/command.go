package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var allowedCommands = map[string]bool{
	"get":           true,
	"describe":      true,
	"logs":          true,
	"top":           true,
	"api-resources": true,
	"version":       true,
}

var disallowedStreamingFlags = map[string]bool{
	"-w": true, "--watch": true, "-f": true, "--follow": true,
}

// disallowedShellFragments are rejected on the raw command string before any
// tokenization, so quoting cannot smuggle them past the check.
var disallowedShellFragments = []string{"&&", "||", "|", ">", "<", "`", "$(", "${", ";", "\n", "\r"}

var (
	commandPrefixRe = regexp.MustCompile(`^(?:kubectl|oc)(?:\s+|$)`)
	commandSplitRe  = regexp.MustCompile(`[;\n]+`)
)

// splitCommandInput splits one raw commands entry on ';' and newlines.
func splitCommandInput(raw any) ([]string, error) {
	command, ok := raw.(string)
	if !ok {
		return nil, errors.New("each command must be a string")
	}
	var parts []string
	for _, part := range commandSplitRe.Split(command, -1) {
		if part = strings.TrimSpace(part); part != "" {
			parts = append(parts, part)
		}
	}
	if len(parts) == 0 {
		return nil, errors.New("command cannot be empty")
	}
	return parts, nil
}

// normalizeCommand validates one command and strips an optional kubectl/oc
// prefix. Only readonly verbs survive; the returned command is what the
// Jenkins job receives.
func normalizeCommand(command string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", errors.New("command cannot be empty")
	}
	for _, fragment := range disallowedShellFragments {
		if strings.Contains(command, fragment) {
			return "", fmt.Errorf("command contains a disallowed shell fragment: %q", command)
		}
	}

	tokens, err := shlexSplit(command)
	if err != nil {
		return "", fmt.Errorf("command has invalid shell quoting: %w", err)
	}
	if len(tokens) == 0 {
		return "", errors.New("command cannot be empty")
	}

	tokenOffset := 0
	if tokens[0] == "kubectl" || tokens[0] == "oc" {
		tokenOffset = 1
	}
	if tokenOffset >= len(tokens) {
		return "", errors.New("command must include a kubectl subcommand")
	}

	if verb := tokens[tokenOffset]; !allowedCommands[verb] {
		return "", errors.New("only readonly kubectl commands are allowed: api-resources, describe, get, logs, top, version")
	}

	for _, token := range tokens[tokenOffset+1:] {
		if disallowedStreamingFlags[token] {
			return "", errors.New("streaming/watch flags are not allowed for Jenkins-executed commands")
		}
	}

	if tokenOffset == 0 {
		return command, nil
	}
	if loc := commandPrefixRe.FindStringIndex(command); loc != nil {
		normalized := strings.TrimSpace(command[loc[1]:])
		if normalized == "" {
			return "", errors.New("command became empty after removing the kubectl prefix")
		}
		return normalized, nil
	}
	return command, nil
}

// normalizeCommands validates the commands argument: a list of strings, each
// possibly holding several commands separated by ';' or newlines.
func normalizeCommands(raw any) ([]string, error) {
	items, ok := raw.([]any)
	if !ok {
		return nil, errors.New("commands must be a list of command strings")
	}
	var normalized []string
	for _, item := range items {
		parts, err := splitCommandInput(item)
		if err != nil {
			return nil, err
		}
		for _, part := range parts {
			command, err := normalizeCommand(part)
			if err != nil {
				return nil, err
			}
			normalized = append(normalized, command)
		}
	}
	if len(normalized) == 0 {
		return nil, errors.New("at least one command is required")
	}
	return normalized, nil
}

// shlexSplit splits a command the way Python's shlex.split does: whitespace
// separated tokens with single quotes, double quotes, and backslash escapes.
func shlexSplit(s string) ([]string, error) {
	var (
		tokens   []string
		current  strings.Builder
		inToken  bool
		inSingle bool
		inDouble bool
	)

	appendToken := func() {
		if inToken {
			tokens = append(tokens, current.String())
			current.Reset()
			inToken = false
		}
	}

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case inSingle:
			if r == '\'' {
				inSingle = false
			} else {
				current.WriteRune(r)
			}
		case inDouble:
			switch r {
			case '"':
				inDouble = false
			case '\\':
				if i+1 < len(runes) && (runes[i+1] == '"' || runes[i+1] == '\\') {
					i++
					current.WriteRune(runes[i])
				} else {
					current.WriteRune(r)
				}
			default:
				current.WriteRune(r)
			}
		default:
			switch r {
			case '\'':
				inSingle = true
				inToken = true
			case '"':
				inDouble = true
				inToken = true
			case '\\':
				if i+1 >= len(runes) {
					return nil, errors.New("no escaped character at end of input")
				}
				i++
				current.WriteRune(runes[i])
				inToken = true
			case ' ', '\t', '\r', '\n':
				appendToken()
			default:
				current.WriteRune(r)
				inToken = true
			}
		}
	}
	if inSingle || inDouble {
		return nil, errors.New("no closing quotation")
	}
	appendToken()
	return tokens, nil
}
