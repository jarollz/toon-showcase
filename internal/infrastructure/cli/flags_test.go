package cli

import "testing"

func TestParseDefaults(t *testing.T) {
	opts, err := Parse(nil)
	if err != nil {
		t.Fatalf("Parse defaults error: %v", err)
	}
	if opts.Records != 120 || opts.Iters != 500 || opts.Warmup != 50 {
		t.Fatalf("unexpected defaults: %+v", opts)
	}
	if opts.Baseline != "toon" {
		t.Fatalf("unexpected default baseline: %+v", opts)
	}
	if opts.CSVDir != "" {
		t.Fatalf("csv dir default should be empty")
	}
	if opts.EncodedExamplesDir != "" {
		t.Fatalf("encoded examples dir default should be empty")
	}
	if opts.MarkdownFile != "" {
		t.Fatalf("markdown file default should be empty")
	}
}

func TestParseCustomAndError(t *testing.T) {
	opts, err := Parse([]string{"-records", "1", "-iters", "2", "-warmup", "3", "-seed", "4", "-indent", "5", "-sample-limit", "6", "-no-progress", "-baseline", "json", "-csv-dir", "stdout", "-md-file", "report.md", "-encoded-examples-dir", "x-output"})
	if err != nil {
		t.Fatalf("Parse custom error: %v", err)
	}
	if opts.Records != 1 || opts.Iters != 2 || opts.Warmup != 3 || opts.Seed != 4 || opts.IndentSize != 5 || opts.SampleLimit != 6 {
		t.Fatalf("unexpected parsed opts: %+v", opts)
	}
	if !opts.NoProgress || opts.Baseline != "json" || opts.CSVDir != "stdout" {
		t.Fatalf("unexpected parsed flags: %+v", opts)
	}
	if opts.EncodedExamplesDir != "x-output" || opts.MarkdownFile != "report.md" {
		t.Fatalf("unexpected new path flags: %+v", opts)
	}

	if _, err := Parse([]string{"-records", "bad"}); err == nil {
		t.Fatalf("expected parse error")
	}
	if _, err := Parse([]string{"-csv", "stdout"}); err == nil {
		t.Fatalf("expected removed -csv flag parse error")
	}
}
