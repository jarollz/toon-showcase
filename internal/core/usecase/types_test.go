package usecase

import (
	"testing"
	"time"
)

func TestNormalizeBaselineKey(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "toon", want: "toon"},
		{in: " JSON ", want: "json-compact"},
		{in: "yml", want: "yaml"},
		{in: "xmlpretty", want: "xml-pretty"},
		{in: "msgpack", want: "messagepack"},
		{in: "bad", wantErr: true},
	}

	for _, tc := range tests {
		got, err := NormalizeBaselineKey(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("NormalizeBaselineKey(%q) expected error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("NormalizeBaselineKey(%q) unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("NormalizeBaselineKey(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestGenerateDatasetDeterministic(t *testing.T) {
	a := GenerateDataset(10, 20260604)
	b := GenerateDataset(10, 20260604)

	if !datasetsEqual(a, b) {
		t.Fatalf("dataset generation should be deterministic for same seed")
	}
}

func TestOptionsValidate(t *testing.T) {
	valid := Options{Records: 1, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "toon"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid options should pass: %v", err)
	}

	tests := []Options{
		{Records: 0, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "toon"},
		{Records: 1, Iters: 0, Warmup: 0, IndentSize: 1, Baseline: "toon"},
		{Records: 1, Iters: 1, Warmup: -1, IndentSize: 1, Baseline: "toon"},
		{Records: 1, Iters: 1, Warmup: 0, IndentSize: 0, Baseline: "toon"},
		{Records: 1, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "bad"},
	}

	for i, tc := range tests {
		if err := tc.Validate(); err == nil {
			t.Fatalf("case %d expected validation error", i)
		}
	}

	if err := (Options{Records: 1, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "toon", EncodedExamplesDir: "   "}).Validate(); err != nil {
		t.Fatalf("encoded examples dir should allow empty/trimmed value: %v", err)
	}
}

func TestResolveEncodedExamplesDir(t *testing.T) {
	now := time.Date(2026, 6, 4, 23, 59, 58, 0, time.FixedZone("UTC+8", 8*60*60))

	got := ResolveEncodedExamplesDir("", now)
	want := ""
	if got != want {
		t.Fatalf("resolved dir = %q, want %q", got, want)
	}

	got = ResolveEncodedExamplesDir(" custom/path ", now)
	if got != "custom/path" {
		t.Fatalf("custom dir should be trimmed, got %q", got)
	}
}

func TestGenerateDatasetInvariants(t *testing.T) {
	d := GenerateDataset(5, 99)
	if len(d.OrdersByID) != 5 {
		t.Fatalf("orders len=%d want 5", len(d.OrdersByID))
	}
	if len(d.OrderTable) != 5 {
		t.Fatalf("order table len=%d want 5", len(d.OrderTable))
	}
	if len(d.Meta) == 0 {
		t.Fatalf("expected metadata")
	}
}
