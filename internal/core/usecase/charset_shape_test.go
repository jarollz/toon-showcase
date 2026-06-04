package usecase

import (
	"encoding/json"
	"errors"
	"testing"

	"toon-showcase/internal/core/entity"
)

func TestRunCharsetMatrix(t *testing.T) {
	ok := successAnyCodec()
	ok.name = "OK"

	badMarshal := successAnyCodec()
	badMarshal.name = "BadMarshal"
	badMarshal.marshalAny = func(any) ([]byte, error) { return nil, errors.New("m") }

	badUnmarshal := successAnyCodec()
	badUnmarshal.name = "BadUnmarshal"
	badUnmarshal.unmarshalAny = func([]byte, any) error { return errors.New("u") }

	mismatch := successAnyCodec()
	mismatch.name = "Mismatch"
	mismatch.unmarshalAny = func(_ []byte, v any) error {
		p := anyToCharsetProbe(v)
		p.Owner = "different"
		p.Meta = []entity.TextEntry{{Key: "probe", Value: "different"}}
		return nil
	}

	report := runCharsetMatrix([]Codec{ok, badMarshal, badUnmarshal, mismatch})
	if len(report.Tests) == 0 || len(report.FormatNames) != 4 {
		t.Fatalf("unexpected report size")
	}
	if report.PassCounts[0] != len(report.Tests) {
		t.Fatalf("ok codec should pass all")
	}
	if report.PassCounts[1] != 0 || report.PassCounts[2] != 0 || report.PassCounts[3] != 0 {
		t.Fatalf("non-ok codecs should fail")
	}
}

func TestRunShapeMatrix(t *testing.T) {
	ok := successAnyCodec()
	ok.name = "OK"

	badMarshal := successAnyCodec()
	badMarshal.name = "BadMarshal"
	badMarshal.marshalAny = func(any) ([]byte, error) { return nil, errors.New("m") }

	badUnmarshal := successAnyCodec()
	badUnmarshal.name = "BadUnmarshal"
	badUnmarshal.unmarshalAny = func([]byte, any) error { return errors.New("u") }

	report := runShapeMatrix([]Codec{ok, badMarshal, badUnmarshal})
	if len(report.Tests) == 0 || len(report.FormatNames) != 3 {
		t.Fatalf("unexpected report size")
	}
	if report.PassCounts[0] != len(report.Tests) {
		t.Fatalf("ok codec should pass all shape tests")
	}
	if report.PassCounts[1] != 0 || report.PassCounts[2] != 0 {
		t.Fatalf("error codecs should fail all")
	}

	wantArtifacts := len(report.Tests) * len(report.FormatNames)
	if len(report.EncodedArtifacts) != wantArtifacts {
		t.Fatalf("artifacts len = %d, want %d", len(report.EncodedArtifacts), wantArtifacts)
	}

	okArtifact, found := findShapeArtifact(report, "OK", "slice_struct")
	if !found {
		t.Fatalf("missing OK/slice_struct artifact")
	}
	if !okArtifact.Result.EncodeOK || len(okArtifact.Encoded) == 0 {
		t.Fatalf("expected encoded bytes for successful artifact")
	}

	errArtifact, found := findShapeArtifact(report, "BadMarshal", "slice_struct")
	if !found {
		t.Fatalf("missing BadMarshal/slice_struct artifact")
	}
	if errArtifact.Result.EncodeOK || len(errArtifact.Encoded) != 0 {
		t.Fatalf("expected failed artifact with no encoded bytes")
	}
}

func TestEvalShapeCaseMismatch(t *testing.T) {
	c := successAnyCodec()
	c.unmarshalAny = func(_ []byte, out any) error {
		switch p := out.(type) {
		case *probeSliceStruct:
			p.Value = []entity.OrderSummary{{OrderID: "DIFF"}}
		}
		return nil
	}

	tc := shapeCase{
		Name: "custom",
		Prepare: func() (any, any) {
			in := probeSliceStruct{Value: []entity.OrderSummary{{OrderID: "A"}}}
			var out probeSliceStruct
			return in, &out
		},
	}

	eval := evalShapeCase(c, tc)
	res := eval.Result
	if res.RoundTripOK {
		t.Fatalf("expected mismatch")
	}
}

func successAnyCodec() fakeCodec {
	c := successCodec()
	c.marshalAny = json.Marshal
	c.unmarshalAny = func(raw []byte, out any) error {
		return json.Unmarshal(raw, out)
	}
	return c
}

func anyToCharsetProbe(v any) *charsetProbe {
	p, _ := v.(*charsetProbe)
	if p == nil {
		return &charsetProbe{}
	}
	return p
}

func findShapeArtifact(report entity.ShapeMatrix, formatName, caseName string) (entity.ShapeEncodedArtifact, bool) {
	for _, artifact := range report.EncodedArtifacts {
		if artifact.FormatName == formatName && artifact.CaseName == caseName {
			return artifact, true
		}
	}
	return entity.ShapeEncodedArtifact{}, false
}
