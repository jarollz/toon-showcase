package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRunSuccessNoCSV(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "5", "-iters", "3", "-warmup", "1", "-no-progress", "-seed", "42"}, stdout, stderr, func() time.Time {
		return time.Unix(100, 0)
	})

	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "Multi-Format Marshal/Unmarshal Showcase") {
		t.Fatalf("stdout missing benchmark section")
	}
	if !strings.Contains(out, "Character Set Support Comparison") {
		t.Fatalf("stdout missing charset section")
	}
}

func TestRunInvalidArgs(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "nope"}, stdout, stderr, time.Now)

	if code != 2 {
		t.Fatalf("run code = %d, want 2", code)
	}
	if stderr.Len() == 0 {
		t.Fatalf("expected stderr output")
	}
}

func TestRunCSVFileError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{
		"-records", "3",
		"-iters", "2",
		"-warmup", "0",
		"-seed", "42",
		"-no-progress",
		"-csv", "file",
		"-csv-file", "/not/a/real/dir/benchmark.csv",
	}, stdout, stderr, time.Now)

	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "failed to write benchmark csv") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}
