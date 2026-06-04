package usecase

import "toon-showcase/internal/core/entity"

type charsetProbe struct {
	Owner string             `json:"owner" yaml:"owner" toml:"owner" msgpack:"owner" xml:"owner" toon:"owner"`
	Meta  []entity.TextEntry `json:"meta" yaml:"meta" toml:"meta" msgpack:"meta" xml:"meta>entry" toon:"meta"`
}

func runCharsetMatrix(codecs []Codec) entity.CharsetMatrix {
	tests := []entity.CharsetCase{
		{Name: "ASCII", Value: "hello-world_123"},
		{Name: "English", Value: "Quick brown fox"},
		{Name: "Japanese", Value: "こんにちは世界"},
		{Name: "Chinese", Value: "你好，世界"},
		{Name: "Emoji", Value: "🚀✨🙂"},
		{Name: "Mixed UTF-8", Value: "hello 你好 こんにちは 🚀"},
		{Name: "Newline/Tab", Value: "line1\nline2\tcolumn"},
		{Name: "Quote/Backslash", Value: "say \"hi\" \\ done"},
		{Name: "Delimiter chars", Value: "comma, pipe| colon:"},
		{Name: "Control U+0001", Value: string([]byte{0x01})},
		{Name: "Control U+0000", Value: string([]byte{0x00})},
	}

	resultsByFormat := make([][]entity.CharsetResult, len(codecs))
	passCounts := make([]int, len(codecs))
	formatNames := make([]string, len(codecs))

	for fi, c := range codecs {
		formatNames[fi] = c.Name()
		rows := make([]entity.CharsetResult, 0, len(tests))
		for _, tc := range tests {
			res := evalCharsetCase(c, tc)
			if res.MarshalOK && res.RoundTripOK {
				passCounts[fi]++
			}
			rows = append(rows, res)
		}
		resultsByFormat[fi] = rows
	}

	return entity.CharsetMatrix{Tests: tests, FormatNames: formatNames, ResultsByFormat: resultsByFormat, PassCounts: passCounts}
}

func evalCharsetCase(codec Codec, tc entity.CharsetCase) entity.CharsetResult {
	probe := charsetProbe{Owner: tc.Value, Meta: []entity.TextEntry{{Key: "probe", Value: tc.Value}}}

	raw, err := codec.MarshalAny(probe)
	if err != nil {
		return entity.CharsetResult{MarshalOK: false, RoundTripOK: false, Detail: err.Error()}
	}
	var out charsetProbe
	if err := codec.UnmarshalAny(raw, &out); err != nil {
		return entity.CharsetResult{MarshalOK: true, RoundTripOK: false, Detail: err.Error()}
	}
	if out.Owner != probe.Owner || len(out.Meta) == 0 || out.Meta[0].Value != probe.Meta[0].Value {
		return entity.CharsetResult{MarshalOK: true, RoundTripOK: false, Detail: "decoded value differs"}
	}
	return entity.CharsetResult{MarshalOK: true, RoundTripOK: true, Detail: "ok"}
}
