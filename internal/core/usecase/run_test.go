package usecase

import (
	"errors"
	"testing"
	"time"

	"toon-showcase/internal/core/entity"
)

type fakeCodec struct {
	key          string
	name         string
	binary       bool
	encodeData   func(entity.BenchmarkData) ([]byte, error)
	decodeData   func([]byte) (entity.BenchmarkData, error)
	marshalAny   func(any) ([]byte, error)
	unmarshalAny func([]byte, any) error
}

func (f fakeCodec) Key() string { return f.key }
func (f fakeCodec) Name() string { return f.name }
func (f fakeCodec) Binary() bool { return f.binary }
func (f fakeCodec) EncodeData(d entity.BenchmarkData) ([]byte, error) { return f.encodeData(d) }
func (f fakeCodec) DecodeData(b []byte) (entity.BenchmarkData, error) { return f.decodeData(b) }
func (f fakeCodec) MarshalAny(v any) ([]byte, error) { return f.marshalAny(v) }
func (f fakeCodec) UnmarshalAny(b []byte, v any) error { return f.unmarshalAny(b, v) }

type fakeProgress struct {
	setStageCount int
	advanceCount  int
	started       bool
	finished      bool
}

func (p *fakeProgress) SetStage(string) { p.setStageCount++ }
func (p *fakeProgress) Advance(n int64) {
	if n > 0 {
		p.advanceCount += int(n)
	}
}
func (p *fakeProgress) Start() { p.started = true }
func (p *fakeProgress) Finish() { p.finished = true }

func successCodec() fakeCodec {
	return fakeCodec{
		key:    "toon",
		name:   "TOON",
		binary: false,
		encodeData: func(entity.BenchmarkData) ([]byte, error) {
			return []byte(`{"ok":true}`), nil
		},
		decodeData: func([]byte) (entity.BenchmarkData, error) {
			return entity.BenchmarkData{}, nil
		},
		marshalAny: func(any) ([]byte, error) {
			return []byte("x"), nil
		},
		unmarshalAny: func([]byte, any) error {
			return nil
		},
	}
}

func TestRunnerRunSuccess(t *testing.T) {
	prog := &fakeProgress{}
	r := Runner{
		Codecs: []Codec{successCodec()},
		Now:    func() time.Time { return time.Unix(123, 0) },
		NewProgress: func(total int64) Progress {
			if total <= 0 {
				t.Fatalf("expected positive total")
			}
			return prog
		},
	}

	out, err := r.Run(Options{Records: 1, Iters: 1, Warmup: 0, Seed: -1, IndentSize: 2, SampleLimit: 10, Baseline: "toon", CSVMode: "off"})
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}
	if out.ResolvedSeed != 123 {
		t.Fatalf("resolved seed = %d, want 123", out.ResolvedSeed)
	}
	if len(out.Results) != 1 {
		t.Fatalf("results len = %d, want 1", len(out.Results))
	}
	if out.Baseline.Key != "toon" {
		t.Fatalf("baseline key = %s, want toon", out.Baseline.Key)
	}
	if !prog.started || !prog.finished {
		t.Fatalf("progress lifecycle not complete")
	}
}

func TestRunnerRunNoProgress(t *testing.T) {
	r := Runner{
		Codecs: []Codec{successCodec()},
		Now:    time.Now,
		NewProgress: func(total int64) Progress {
			t.Fatalf("progress should not be created when NoProgress=true")
			return nil
		},
	}
	_, err := r.Run(Options{Records: 1, Iters: 1, Warmup: 0, Seed: 1, IndentSize: 2, SampleLimit: 10, NoProgress: true, Baseline: "toon", CSVMode: "off"})
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}
}

func TestRunnerRunErrors(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		r := Runner{Codecs: []Codec{successCodec()}, Now: time.Now}
		_, err := r.Run(Options{Records: 0, Iters: 1, Warmup: 0, Seed: 1, IndentSize: 2, SampleLimit: 10, Baseline: "toon", CSVMode: "off"})
		if err == nil {
			t.Fatalf("expected validation error")
		}
	})

	t.Run("no codecs", func(t *testing.T) {
		r := Runner{Now: time.Now}
		_, err := r.Run(Options{Records: 1, Iters: 1, Warmup: 0, Seed: 1, IndentSize: 2, SampleLimit: 10, Baseline: "toon", CSVMode: "off"})
		if err == nil {
			t.Fatalf("expected no codecs error")
		}
	})

	t.Run("nil clock", func(t *testing.T) {
		r := Runner{Codecs: []Codec{successCodec()}}
		_, err := r.Run(Options{Records: 1, Iters: 1, Warmup: 0, Seed: -1, IndentSize: 2, SampleLimit: 10, Baseline: "toon", CSVMode: "off"})
		if err == nil {
			t.Fatalf("expected clock error")
		}
	})

	t.Run("benchmark failed", func(t *testing.T) {
		bad := successCodec()
		bad.encodeData = func(entity.BenchmarkData) ([]byte, error) {
			return nil, errors.New("boom")
		}
		p := &fakeProgress{}
		r := Runner{
			Codecs: []Codec{bad},
			Now:    time.Now,
			NewProgress: func(int64) Progress {
				return p
			},
		}
		_, err := r.Run(Options{Records: 1, Iters: 1, Warmup: 0, Seed: 1, IndentSize: 2, SampleLimit: 10, Baseline: "toon", CSVMode: "off"})
		if err == nil {
			t.Fatalf("expected benchmark error")
		}
		if !p.finished {
			t.Fatalf("progress should be finished on error")
		}
	})

	t.Run("baseline not found", func(t *testing.T) {
		c := successCodec()
		c.key = "json-compact"
		r := Runner{Codecs: []Codec{c}, Now: time.Now}
		_, err := r.Run(Options{Records: 1, Iters: 1, Warmup: 0, Seed: 1, IndentSize: 2, SampleLimit: 10, NoProgress: true, Baseline: "toon", CSVMode: "off"})
		if err == nil {
			t.Fatalf("expected baseline missing error")
		}
	})
}

func TestTotalBenchmarkSteps(t *testing.T) {
	if got := totalBenchmarkSteps(0, 10, 1); got != 1 {
		t.Fatalf("totalBenchmarkSteps zero formats = %d, want 1", got)
	}
	if got := totalBenchmarkSteps(1, -5, -5); got != 1 {
		t.Fatalf("totalBenchmarkSteps non-positive total = %d, want 1", got)
	}
	if got := totalBenchmarkSteps(2, 3, 1); got != int64(2*(1*2+3*2+1)) {
		t.Fatalf("unexpected total steps: %d", got)
	}
}
