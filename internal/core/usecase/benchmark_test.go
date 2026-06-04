package usecase

import (
	"errors"
	"testing"

	"toon-showcase/internal/core/entity"
)

func TestRunBenchmarkErrors(t *testing.T) {
	data := entity.BenchmarkData{}

	t.Run("warmup marshal", func(t *testing.T) {
		c := successCodec()
		c.encodeData = func(entity.BenchmarkData) ([]byte, error) { return nil, errors.New("e1") }
		_, err := runBenchmark(c, data, 1, 1, 5, nil)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("warmup unmarshal", func(t *testing.T) {
		c := successCodec()
		c.decodeData = func([]byte) (entity.BenchmarkData, error) { return entity.BenchmarkData{}, errors.New("e2") }
		_, err := runBenchmark(c, data, 1, 1, 5, nil)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("marshal iteration", func(t *testing.T) {
		count := 0
		c := successCodec()
		c.encodeData = func(entity.BenchmarkData) ([]byte, error) {
			count++
			if count > 0 {
				return nil, errors.New("e3")
			}
			return []byte("x"), nil
		}
		_, err := runBenchmark(c, data, 1, 0, 5, nil)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("unmarshal iteration", func(t *testing.T) {
		count := 0
		c := successCodec()
		c.decodeData = func([]byte) (entity.BenchmarkData, error) {
			count++
			if count == 1 {
				return entity.BenchmarkData{}, nil
			}
			return entity.BenchmarkData{}, errors.New("e4")
		}
		_, err := runBenchmark(c, data, 2, 0, 5, nil)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("final decode", func(t *testing.T) {
		count := 0
		c := successCodec()
		c.decodeData = func([]byte) (entity.BenchmarkData, error) {
			count++
			if count <= 1 {
				return entity.BenchmarkData{}, nil
			}
			return entity.BenchmarkData{}, errors.New("e5")
		}
		_, err := runBenchmark(c, data, 1, 0, 5, nil)
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestRunBenchmarkSuccessAndMismatch(t *testing.T) {
	p := &fakeProgress{}
	c := successCodec()
	res, err := runBenchmark(c, entity.BenchmarkData{}, 2, 1, 6, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Sample == "" {
		t.Fatalf("expected sample")
	}
	if !res.RoundTripOK {
		t.Fatalf("expected roundtrip ok")
	}
	if p.advanceCount == 0 {
		t.Fatalf("expected progress advances")
	}

	mismatch := successCodec()
	mismatch.decodeData = func([]byte) (entity.BenchmarkData, error) {
		return entity.BenchmarkData{Version: "x"}, nil
	}
	res, err = runBenchmark(mismatch, entity.BenchmarkData{}, 1, 0, 4, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.RoundTripOK {
		t.Fatalf("expected mismatch")
	}
}

func TestFormatSampleAndHelpers(t *testing.T) {
	if got := formatSample([]byte("abc"), false, 0); got != "" {
		t.Fatalf("limit zero expected empty")
	}
	if got := formatSample([]byte("abcdef"), true, 4); got != "6162..." {
		t.Fatalf("binary sample = %q", got)
	}
	if got := formatSample([]byte("a\nb"), false, 10); got != "a\\nb" {
		t.Fatalf("text sample = %q", got)
	}
	if got := formatSample([]byte("abcdef"), false, 3); got != "abc..." {
		t.Fatalf("truncated text sample = %q", got)
	}

	if datasetsEqual(entity.BenchmarkData{Version: "a"}, entity.BenchmarkData{Version: "b"}) {
		t.Fatalf("datasets should differ")
	}
	if normalizedEqual(struct{ A int }{A: 1}, struct{ A int }{A: 1}) == false {
		t.Fatalf("normalized equal should be true")
	}
}
