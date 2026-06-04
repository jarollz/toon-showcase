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

func TestRunCSVDirError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{
		"-records", "3",
		"-iters", "2",
		"-warmup", "0",
		"-seed", "42",
		"-no-progress",
		"-csv-dir", "/not/a/real/dir/benchmark.csv",
	}, stdout, stderr, time.Now)

	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "invalid csv dir") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunCSVDirStdout(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-csv-dir", "stdout"}, stdout, stderr, time.Now)
	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "CSV Benchmark Output") {
		t.Fatalf("stdout missing csv output: %s", stdout.String())
	}
}

func TestRunEncodedExamplesDir(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	dir := t.TempDir()
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Join(dir, "shape-encoded"))
	})

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-encoded-examples-dir", dir}, stdout, stderr, func() time.Time {
		return time.Date(2026, 6, 4, 23, 59, 58, 0, time.FixedZone("UTC+8", 8*60*60))
	})

	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Encoded examples output dir: "+dir) {
		t.Fatalf("stdout missing output examples dir line: %s", stdout.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "shape-encoded")); err != nil {
		t.Fatalf("expected shape-encoded folder: %v", err)
	}
}

func TestRunEncodedExamplesDirWriteError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-encoded-examples-dir", "/dev/null/bad"}, stdout, stderr, time.Now)

	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "invalid encoded examples dir") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunMarkdownFile(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mdFile := filepath.Join(t.TempDir(), "report.md")

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-md-file", mdFile}, stdout, stderr, time.Now)
	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	b, err := os.ReadFile(mdFile)
	if err != nil {
		t.Fatalf("failed read markdown file: %v", err)
	}
	if !strings.Contains(string(b), "| Format |") {
		t.Fatalf("markdown output missing table: %s", string(b))
	}
	if !strings.Contains(string(b), "## Character Set Support Comparison") {
		t.Fatalf("markdown output missing charset section: %s", string(b))
	}
	if !strings.Contains(string(b), "## Shape Compatibility Matrix") {
		t.Fatalf("markdown output missing shape section: %s", string(b))
	}
}

func TestRunMarkdownStdout(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-md-file", "stdout"}, stdout, stderr, time.Now)
	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "\x1b[") {
		t.Fatalf("expected ansi markdown output, got: %s", stdout.String())
	}
	if strings.Contains(stdout.String(), "Multi-Format Marshal/Unmarshal Showcase") {
		t.Fatalf("stdout markdown mode should not include plain text report header")
	}
}

func TestRunCSVDirWritesFiles(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	csvDir := filepath.Join(t.TempDir(), "csv")

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-csv-dir", csvDir}, stdout, stderr, time.Now)
	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(csvDir, "benchmark_report.csv")); err != nil {
		t.Fatalf("missing benchmark csv: %v", err)
	}
	if _, err := os.Stat(filepath.Join(csvDir, "benchmark_report.charset.csv")); err != nil {
		t.Fatalf("missing charset csv: %v", err)
	}
	if _, err := os.Stat(filepath.Join(csvDir, "benchmark_report.shape.csv")); err != nil {
		t.Fatalf("missing shape csv: %v", err)
	}
}

func TestRunCSVDirExistingFileError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	filePath := filepath.Join(t.TempDir(), "x.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("setup file failed: %v", err)
	}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-csv-dir", filePath}, stdout, stderr, time.Now)
	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "invalid csv dir") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunEncodedExamplesDirExistingFileError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	filePath := filepath.Join(t.TempDir(), "x.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("setup file failed: %v", err)
	}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-encoded-examples-dir", filePath}, stdout, stderr, time.Now)
	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "invalid encoded examples dir") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunAllOutputsFileModes(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	base := t.TempDir()
	csvDir := filepath.Join(base, "csv")
	encodedDir := filepath.Join(base, "encoded")
	mdFile := filepath.Join(base, "md", "report.md")

	code := run([]string{
		"-records", "3",
		"-iters", "2",
		"-warmup", "0",
		"-seed", "42",
		"-no-progress",
		"-csv-dir", csvDir,
		"-encoded-examples-dir", encodedDir,
		"-md-file", mdFile,
	}, stdout, stderr, time.Now)
	if code != 0 {
		t.Fatalf("run code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(mdFile); err != nil {
		t.Fatalf("missing markdown file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(encodedDir, "shape-encoded")); err != nil {
		t.Fatalf("missing encoded examples folder: %v", err)
	}
	if _, err := os.Stat(filepath.Join(csvDir, "benchmark_report.csv")); err != nil {
		t.Fatalf("missing benchmark csv: %v", err)
	}
}

func TestRunMarkdownWriteError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mdPath := t.TempDir()

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-md-file", mdPath}, stdout, stderr, time.Now)
	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "failed to write markdown file") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunMarkdownParentDirError(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-md-file", "/dev/null/report.md"}, stdout, stderr, time.Now)
	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "invalid md-file path") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunInvalidBaseline(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	code := run([]string{"-records", "3", "-iters", "2", "-warmup", "0", "-seed", "42", "-no-progress", "-baseline", "bad"}, stdout, stderr, time.Now)
	if code != 1 {
		t.Fatalf("run code = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "invalid baseline") {
		t.Fatalf("unexpected stderr: %s", stderr.String())
	}
}
