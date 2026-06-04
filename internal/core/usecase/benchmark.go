package usecase

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"toon-showcase/internal/core/entity"
)

func runBenchmark(codec Codec, data entity.BenchmarkData, iters, warmup, sampleLimit int, progress Progress) (entity.BenchmarkResult, error) {
	for i := 0; i < warmup; i++ {
		raw, err := codec.EncodeData(data)
		if err != nil {
			return entity.BenchmarkResult{}, fmt.Errorf("warmup marshal: %w", err)
		}
		if progress != nil {
			progress.Advance(1)
		}
		if _, err := codec.DecodeData(raw); err != nil {
			return entity.BenchmarkResult{}, fmt.Errorf("warmup unmarshal: %w; encoded sample=%s", err, formatSample(raw, codec.Binary(), 4000))
		}
		if progress != nil {
			progress.Advance(1)
		}
	}

	marshalStart := time.Now()
	var encoded []byte
	for i := 0; i < iters; i++ {
		raw, err := codec.EncodeData(data)
		if err != nil {
			return entity.BenchmarkResult{}, fmt.Errorf("marshal iteration %d: %w", i, err)
		}
		encoded = raw
		if progress != nil {
			progress.Advance(1)
		}
	}
	marshalTotal := time.Since(marshalStart)

	unmarshalStart := time.Now()
	for i := 0; i < iters; i++ {
		if _, err := codec.DecodeData(encoded); err != nil {
			return entity.BenchmarkResult{}, fmt.Errorf("unmarshal iteration %d: %w", i, err)
		}
		if progress != nil {
			progress.Advance(1)
		}
	}
	unmarshalTotal := time.Since(unmarshalStart)

	out, err := codec.DecodeData(encoded)
	if err != nil {
		return entity.BenchmarkResult{}, fmt.Errorf("final decode for roundtrip: %w", err)
	}
	if progress != nil {
		progress.Advance(1)
	}

	encodedChars := utf8.RuneCount(encoded)
	if codec.Binary() {
		encodedChars = -1
	}

	res := entity.BenchmarkResult{
		Key:            codec.Key(),
		Name:           codec.Name(),
		Binary:         codec.Binary(),
		MarshalTotal:   marshalTotal,
		MarshalAvg:     time.Duration(int64(marshalTotal) / int64(iters)),
		UnmarshalTotal: unmarshalTotal,
		UnmarshalAvg:   time.Duration(int64(unmarshalTotal) / int64(iters)),
		EncodedBytes:   len(encoded),
		EncodedChars:   encodedChars,
		Sample:         formatSample(encoded, codec.Binary(), sampleLimit),
	}

	if datasetsEqual(data, out) {
		res.RoundTripOK = true
		res.RoundTripDetails = "ok"
	} else {
		res.RoundTripOK = false
		res.RoundTripDetails = "mismatch with original dataset"
	}

	return res, nil
}

func datasetsEqual(a, b entity.BenchmarkData) bool {
	ja, err := json.Marshal(a)
	if err != nil {
		return false
	}
	jb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(ja) == string(jb)
}

func normalizedEqual(a, b any) bool {
	ja, err := json.Marshal(a)
	if err != nil {
		return false
	}
	jb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(ja) == string(jb)
}

func formatSample(raw []byte, binary bool, limit int) string {
	if limit <= 0 {
		return ""
	}
	if binary {
		hexText := hex.EncodeToString(raw)
		if len(hexText) <= limit {
			return hexText
		}
		return hexText[:limit] + "..."
	}
	text := string(raw)
	r := []rune(text)
	if len(r) <= limit {
		return strings.ReplaceAll(text, "\n", "\\n")
	}
	return strings.ReplaceAll(string(r[:limit]), "\n", "\\n") + "..."
}
