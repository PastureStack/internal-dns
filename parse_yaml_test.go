package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAnswersYAMLIncludesPTRAndTXT(t *testing.T) {
	path := filepath.Join(t.TempDir(), "answers.yaml")
	content := []byte(`default:
  authorative: [pasture.internal]
  ptr:
    10.0.0.8:
      answer: api.pasture.internal.
  txt:
    api.pasture.internal.:
      answer: [ready]
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	answers, err := ParseAnswers(path)
	if err != nil {
		t.Fatal(err)
	}
	defaults := answers[DEFAULT_KEY]
	if len(defaults.Ptr) != 1 || defaults.Ptr["8.0.0.10.in-addr.arpa."].Answer != "api.pasture.internal." {
		t.Fatalf("PTR records = %#v", defaults.Ptr)
	}
	if got := defaults.Txt["api.pasture.internal."].Answer; len(got) != 1 || got[0] != "ready" {
		t.Fatalf("TXT record = %#v", got)
	}
}

func TestParseAnswersAcceptsAuthoritativeAliases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "canonical YAML key",
			content: `default:
  authoritative: [pasture.internal]
`,
			want: "pasture.internal",
		},
		{
			name: "legacy misspelled YAML key",
			content: `default:
  authorative: [legacy.internal]
`,
			want: "legacy.internal",
		},
		{
			name:    "canonical JSON key",
			content: `{"default":{"authoritative":["pasture.internal"]}}`,
			want:    "pasture.internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "answers.yaml")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatal(err)
			}
			answers, err := ParseAnswers(path)
			if err != nil {
				t.Fatal(err)
			}
			got := answers[DEFAULT_KEY].Authoritative
			if len(got) != 1 || got[0] != tt.want {
				t.Fatalf("authoritative suffixes = %#v, want %q", got, tt.want)
			}
		})
	}
}

func TestClientAnswersMarshalUsesCanonicalAuthoritativeKey(t *testing.T) {
	data, err := json.Marshal(ClientAnswers{Authoritative: []string{"pasture.internal"}})
	if err != nil {
		t.Fatal(err)
	}
	encoded := string(data)
	if !strings.Contains(encoded, `"authoritative"`) || strings.Contains(encoded, `"authorative"`) {
		t.Fatalf("encoded client answers use a non-canonical key: %s", encoded)
	}
}
