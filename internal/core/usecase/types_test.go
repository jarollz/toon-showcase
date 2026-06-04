package usecase

import "testing"

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
	valid := Options{Records: 1, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "toon", CSVMode: "off"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid options should pass: %v", err)
	}

	tests := []Options{
		{Records: 0, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "toon", CSVMode: "off"},
		{Records: 1, Iters: 0, Warmup: 0, IndentSize: 1, Baseline: "toon", CSVMode: "off"},
		{Records: 1, Iters: 1, Warmup: -1, IndentSize: 1, Baseline: "toon", CSVMode: "off"},
		{Records: 1, Iters: 1, Warmup: 0, IndentSize: 0, Baseline: "toon", CSVMode: "off"},
		{Records: 1, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "toon", CSVMode: "bad"},
		{Records: 1, Iters: 1, Warmup: 0, IndentSize: 1, Baseline: "bad", CSVMode: "off"},
	}

	for i, tc := range tests {
		if err := tc.Validate(); err == nil {
			t.Fatalf("case %d expected validation error", i)
		}
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
