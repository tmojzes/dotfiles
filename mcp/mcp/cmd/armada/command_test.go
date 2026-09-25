package main

import (
	"reflect"
	"testing"
)

func TestNormalizeCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
		wantErr bool
	}{
		{"plain subcommand", "get nodes", "get nodes", false},
		{"kubectl prefix stripped", "kubectl get pods -A", "get pods -A", false},
		{"oc prefix stripped", "oc get svc -n default", "get svc -n default", false},
		{"describe", "kubectl describe pod api-0 -n kube-system", "describe pod api-0 -n kube-system", false},
		{"logs", "kubectl logs deploy/api", "logs deploy/api", false},
		{"top", "top nodes", "top nodes", false},
		{"api-resources", "kubectl api-resources", "api-resources", false},
		{"version", "kubectl version", "version", false},
		{"verb only", "kubectl get", "get", false},
		{"surrounding whitespace", "  get nodes  ", "get nodes", false},
		{"quotes preserved", `get pods -l "app=api"`, `get pods -l "app=api"`, false},
		{"disallowed verb", "kubectl delete pod foo", "", true},
		{"apply is not readonly", "kubectl apply -f manifest.yaml", "", true},
		{"exec is not readonly", "kubectl exec -it pod -- sh", "", true},
		{"pipe fragment", "get pods | grep crash", "", true},
		{"redirect fragment", "get pods > /tmp/out", "", true},
		{"command substitution", "get pods $(get pods -o name)", "", true},
		{"newline fragment", "get pods\nget nodes", "", true},
		{"watch flag", "kubectl get pods -w", "", true},
		{"long watch flag", "kubectl get pods --watch", "", true},
		{"follow flag", "kubectl logs api-0 -f", "", true},
		{"no subcommand", "kubectl", "", true},
		{"empty", "   ", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeCommand(tt.command)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("normalizeCommand(%q) = %q, want error", tt.command, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeCommand(%q) unexpected error: %v", tt.command, err)
			}
			if got != tt.want {
				t.Fatalf("normalizeCommand(%q) = %q, want %q", tt.command, got, tt.want)
			}
		})
	}
}

func TestNormalizeCommands(t *testing.T) {
	raw := []string{
		"kubectl get nodes; get ns",
		"kubectl top pods\nget pv",
		"  oc  describe node  ",
	}
	want := []string{"get nodes", "get ns", "top pods", "get pv", "describe node"}
	got, err := normalizeCommands(raw)
	if err != nil {
		t.Fatalf("normalizeCommands: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeCommands = %#v, want %#v", got, want)
	}
}

func TestNormalizeCommandsRejections(t *testing.T) {
	for name, raw := range map[string][]string{
		"nil list":      nil,
		"empty list":    {},
		"empty entries": {";\n"},
	} {
		if _, err := normalizeCommands(raw); err == nil {
			t.Errorf("%s: normalizeCommands(%#v) should fail", name, raw)
		}
	}
}

func TestShlexSplit(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []string
		wantErr bool
	}{
		{"plain", "a b c", []string{"a", "b", "c"}, false},
		{"double quoted", `a "b c" d`, []string{"a", "b c", "d"}, false},
		{"single quoted", `a 'b c' d`, []string{"a", "b c", "d"}, false},
		{"quote inside word", `a b'c'd`, []string{"a", "bcd"}, false},
		{"empty double quotes", `a "" b`, []string{"a", "", "b"}, false},
		{"escaped space", `a\ b`, []string{"a b"}, false},
		{"escaped quote in double quotes", `"a\"b"`, []string{`a"b`}, false},
		{"escaped backslash in double quotes", `"a\\b"`, []string{`a\b`}, false},
		{"unclosed single quote", `'unclosed`, nil, true},
		{"unclosed double quote", `"unclosed`, nil, true},
		{"trailing backslash", `trailing\`, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := shlexSplit(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("shlexSplit(%q) = %#v, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("shlexSplit(%q) unexpected error: %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("shlexSplit(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}
