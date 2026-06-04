package usecase

import (
	"reflect"

	"toon-showcase/internal/core/entity"
)

type shapeCase struct {
	Name    string
	Prepare func() (any, any)
}

type probeSliceStruct struct {
	Value []entity.OrderSummary `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value>item" toon:"value"`
}

type probeSliceMapString struct {
	Value []map[string]string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeMapToSliceStruct struct {
	Value map[string][]entity.OrderSummary `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeSliceAny struct {
	Value []any `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeSlicePointer struct {
	Value []*entity.OrderSummary `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value>item" toon:"value"`
}

type probeNestedSliceMap struct {
	Value [][]map[string]string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeMapNonStringKey struct {
	Value map[int]string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

func runShapeMatrix(codecs []Codec) entity.ShapeMatrix {
	tests := []shapeCase{
		{
			Name: "slice_struct",
			Prepare: func() (any, any) {
				in := probeSliceStruct{Value: []entity.OrderSummary{{OrderID: "A", CustomerID: "C1", Locale: "en-US", Priority: 1, TotalUSD: 11.2, EmojiStatus: "🚀"}}}
				var out probeSliceStruct
				return in, &out
			},
		},
		{
			Name: "slice_map_string",
			Prepare: func() (any, any) {
				in := probeSliceMapString{Value: []map[string]string{{"a": "1", "b": "2"}, {"c": "3"}}}
				var out probeSliceMapString
				return in, &out
			},
		},
		{
			Name: "map_string_to_slice_struct",
			Prepare: func() (any, any) {
				in := probeMapToSliceStruct{Value: map[string][]entity.OrderSummary{"group": {{OrderID: "A", CustomerID: "X", Locale: "ja-JP", Priority: 2, TotalUSD: 33.4, EmojiStatus: "✨"}}}}
				var out probeMapToSliceStruct
				return in, &out
			},
		},
		{
			Name: "slice_any_mixed",
			Prepare: func() (any, any) {
				in := probeSliceAny{Value: []any{1, "a", true}}
				var out probeSliceAny
				return in, &out
			},
		},
		{
			Name: "slice_pointer_with_nil",
			Prepare: func() (any, any) {
				v := &entity.OrderSummary{OrderID: "X", CustomerID: "C", Locale: "zh-CN", Priority: 1, TotalUSD: 9.9, EmojiStatus: "🙂"}
				in := probeSlicePointer{Value: []*entity.OrderSummary{v, nil}}
				var out probeSlicePointer
				return in, &out
			},
		},
		{
			Name: "nested_slice_map",
			Prepare: func() (any, any) {
				in := probeNestedSliceMap{Value: [][]map[string]string{{{"a": "1"}}, {{"b": "2"}}}}
				var out probeNestedSliceMap
				return in, &out
			},
		},
		{
			Name: "map_non_string_key",
			Prepare: func() (any, any) {
				in := probeMapNonStringKey{Value: map[int]string{1: "one", 2: "two"}}
				var out probeMapNonStringKey
				return in, &out
			},
		},
	}

	resultsByFormat := make([][]entity.ShapeResult, len(codecs))
	passCounts := make([]int, len(codecs))
	formatNames := make([]string, len(codecs))
	testNames := make([]string, 0, len(tests))

	for _, t := range tests {
		testNames = append(testNames, t.Name)
	}

	for fi, c := range codecs {
		formatNames[fi] = c.Name()
		rows := make([]entity.ShapeResult, 0, len(tests))
		for _, tc := range tests {
			res := evalShapeCase(c, tc)
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				passCounts[fi]++
			}
			rows = append(rows, res)
		}
		resultsByFormat[fi] = rows
	}

	return entity.ShapeMatrix{Tests: testNames, FormatNames: formatNames, ResultsByFormat: resultsByFormat, PassCounts: passCounts}
}

func evalShapeCase(codec Codec, tc shapeCase) entity.ShapeResult {
	in, outPtr := tc.Prepare()
	raw, err := codec.MarshalAny(in)
	if err != nil {
		return entity.ShapeResult{EncodeOK: false, DecodeOK: false, RoundTripOK: false, Detail: err.Error()}
	}
	if err := codec.UnmarshalAny(raw, outPtr); err != nil {
		return entity.ShapeResult{EncodeOK: true, DecodeOK: false, RoundTripOK: false, Detail: err.Error()}
	}
	outVal := reflect.ValueOf(outPtr).Elem().Interface()
	if !normalizedEqual(in, outVal) {
		return entity.ShapeResult{EncodeOK: true, DecodeOK: true, RoundTripOK: false, Detail: "decoded value differs"}
	}
	return entity.ShapeResult{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"}
}
