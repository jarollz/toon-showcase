package main

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestRunOutputExampleDefaultDir(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join("output", "examples_20260604_235958_+0800"))
	})

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-output-example=true"}, stdout, stderr, func() time.Time {
		return time.Date(2026, 6, 4, 23, 59, 58, 0, time.FixedZone("UTC+8", 8*60*60))
	})

	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Encoded examples output dir: output/examples_20260604_235958_+0800") {
		t.Fatalf("stdout missing output examples dir line: %s", stdout.String())
	}
	if _, err := os.Stat(filepath.Join("output", "examples_20260604_235958_+0800", "shape-encoded")); err != nil {
		t.Fatalf("expected shape-encoded folder: %v", err)
	}
}

func TestRunOutputExampleWriteError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-output-example=true", "-output-example-dir", "/dev/null/bad"}, stdout, stderr, time.Now)

	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "failed to write encoded examples") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}
