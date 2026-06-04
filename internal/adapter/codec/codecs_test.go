package codec

import (
	"encoding/xml"
	"testing"

	"toon-showcase/internal/core/usecase"
)

func TestBuildCodecsRoundTrip(t *testing.T) {
	codecs := Build(2)
	if len(codecs) != 8 {
		t.Fatalf("codec count = %d, want 8", len(codecs))
	}

	data := usecase.GenerateDataset(3, 20260604)
	for _, c := range codecs {
		raw, err := c.EncodeData(data)
		if err != nil {
			t.Fatalf("%s encode error: %v", c.Name(), err)
		}
		if len(raw) == 0 {
			t.Fatalf("%s produced empty payload", c.Name())
		}
		out, err := c.DecodeData(raw)
		if err != nil {
			t.Fatalf("%s decode error: %v", c.Name(), err)
		}
		if len(out.OrderTable) != len(data.OrderTable) {
			t.Fatalf("%s decode mismatch orderTable size", c.Name())
		}

		blob, err := c.MarshalAny(data)
		if err != nil {
			t.Fatalf("%s MarshalAny error: %v", c.Name(), err)
		}
		var outData = usecase.GenerateDataset(0, 1)
		if err := c.UnmarshalAny(blob, &outData); err != nil {
			t.Fatalf("%s UnmarshalAny error: %v", c.Name(), err)
		}
		if outData.Version == "" {
			t.Fatalf("%s MarshalAny/UnmarshalAny mismatch", c.Name())
		}
	}
}

func TestXMLHelpers(t *testing.T) {
	data := usecase.GenerateDataset(2, 1)
	raw, err := marshalXMLBenchmark(data)
	if err != nil {
		t.Fatalf("marshalXMLBenchmark error: %v", err)
	}
	out, err := unmarshalXMLBenchmark(raw)
	if err != nil {
		t.Fatalf("unmarshalXMLBenchmark error: %v", err)
	}
	if len(out.OrdersByID) != len(data.OrdersByID) {
		t.Fatalf("orders len mismatch")
	}

	if _, err := unmarshalXMLBenchmark([]byte("<bad xml")); err == nil {
		t.Fatalf("expected xml parse error")
	}

	if got := orderedOrders(nil); got != nil {
		t.Fatalf("orderedOrders nil should return nil")
	}
}

func TestXMLBenchmarkDataShape(t *testing.T) {
	x := xmlBenchmarkData{Version: "1"}
	raw, err := xml.Marshal(x)
	if err != nil {
		t.Fatalf("marshal xmlBenchmarkData error: %v", err)
	}
	if len(raw) == 0 {
		t.Fatalf("empty xml payload")
	}
}
