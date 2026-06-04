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
	if opts.Baseline != "toon" || opts.CSVMode != "off" {
		t.Fatalf("unexpected default baseline/csv: %+v", opts)
	}
}

func TestParseCustomAndError(t *testing.T) {
	opts, err := Parse([]string{"-records", "1", "-iters", "2", "-warmup", "3", "-seed", "4", "-indent", "5", "-sample-limit", "6", "-no-progress", "-baseline", "json", "-csv", "file", "-csv-file", "x.csv"})
	if err != nil {
		t.Fatalf("Parse custom error: %v", err)
	}
	if opts.Records != 1 || opts.Iters != 2 || opts.Warmup != 3 || opts.Seed != 4 || opts.IndentSize != 5 || opts.SampleLimit != 6 {
		t.Fatalf("unexpected parsed opts: %+v", opts)
	}
	if !opts.NoProgress || opts.Baseline != "json" || opts.CSVMode != "file" || opts.CSVFile != "x.csv" {
		t.Fatalf("unexpected parsed flags: %+v", opts)
	}

	if _, err := Parse([]string{"-records", "bad"}); err == nil {
		t.Fatalf("expected parse error")
	}
}
