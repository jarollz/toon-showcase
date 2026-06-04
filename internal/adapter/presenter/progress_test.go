package presenter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestProgressBarLifecycle(t *testing.T) {
	pAny := NewProgressBar(10)
	p, ok := pAny.(*progressBar)
	if !ok {
		t.Fatalf("NewProgressBar returned unexpected type")
	}

	p.SetStage("stage-1")
	p.Advance(0)
	p.Advance(2)
	p.Start()
	p.Start()
	time.Sleep(5 * time.Millisecond)
	p.Finish()

	if !p.started {
		t.Fatalf("progress should be started")
	}
}

func TestProgressBarNoStartFinish(t *testing.T) {
	pAny := NewProgressBar(0)
	p := pAny.(*progressBar)
	p.Finish()
	p.render(true)
}

func TestProgressBarRenderClearsStaleTextAndFramesWithNewlines(t *testing.T) {
	buf := &bytes.Buffer{}
	pAny := NewProgressBarWithWriter(10, buf)
	p := pAny.(*progressBar)

	p.SetStage("very long stage")
	p.render(false)
	p.SetStage("short")
	p.render(false)
	p.render(true)

	out := buf.String()
	if !strings.HasPrefix(out, "\n") {
		t.Fatalf("progress output should start with newline, got: %q", out)
	}
	if !strings.Contains(out, "short ") {
		t.Fatalf("expected whitespace cleanup for shorter render, got: %q", out)
	}
	if !strings.HasSuffix(out, "\n\n") {
		t.Fatalf("progress output should end with blank line, got: %q", out)
	}
}
