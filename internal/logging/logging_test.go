package logging

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestInfoProducesSingleRecord(t *testing.T) {
	var output bytes.Buffer
	originalWriter := log.Writer()
	originalFlags := log.Flags()
	defer func() {
		log.SetOutput(originalWriter)
		log.SetFlags(originalFlags)
	}()
	log.SetOutput(&output)
	log.SetFlags(0)

	Infof("source=%s", "first\r\nforged\nthird")
	if got := strings.TrimSuffix(output.String(), "\n"); got != "source=first forged third" {
		t.Fatalf("unexpected log record: %q", got)
	}
}
