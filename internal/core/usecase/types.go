package usecase

import (
	"fmt"
	"strings"
	"time"

	"toon-showcase/internal/core/entity"
)

type Options struct {
	Records          int
	Iters            int
	Warmup           int
	Seed             int64
	IndentSize       int
	SampleLimit      int
	NoProgress       bool
	Baseline         string
	CSVMode          string
	CSVFile          string
	OutputExample    bool
	OutputExampleDir string
}

func (o Options) Validate() error {
	if o.Records <= 0 || o.Iters <= 0 || o.Warmup < 0 || o.IndentSize <= 0 {
		return fmt.Errorf("records and iters must be > 0; warmup >= 0; indent > 0")
	}
	if o.CSVMode != "off" && o.CSVMode != "stdout" && o.CSVMode != "file" {
		return fmt.Errorf("csv mode must be one of: off, stdout, file")
	}
	_, err := NormalizeBaselineKey(o.Baseline)
	if err != nil {
		return err
	}
	return nil
}

func ResolveOutputExampleDir(raw string, now time.Time) string {
	dir := strings.TrimSpace(raw)
	if dir != "" {
		return dir
	}
	return fmt.Sprintf("output/examples_%s_%s_%s", now.Format("20060102"), now.Format("150405"), now.Format("-0700"))
}

type RunOutput struct {
	Results                  []entity.BenchmarkResult
	Baseline                 entity.BenchmarkResult
	Charset                  entity.CharsetMatrix
	Shape                    entity.ShapeMatrix
	ResolvedSeed             int64
	BaselineKey              string
	ResolvedOutputExampleDir string
}

type Codec interface {
	Key() string
	Name() string
	Binary() bool
	EncodeData(entity.BenchmarkData) ([]byte, error)
	DecodeData([]byte) (entity.BenchmarkData, error)
	MarshalAny(any) ([]byte, error)
	UnmarshalAny([]byte, any) error
}

type Progress interface {
	SetStage(string)
	Advance(int64)
	Start()
	Finish()
}

type Runner struct {
	Codecs      []Codec
	Now         func() time.Time
	NewProgress func(total int64) Progress
}

func NormalizeBaselineKey(raw string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(raw))
	aliases := map[string]string{
		"toon":         "toon",
		"json-compact": "json-compact",
		"json_compact": "json-compact",
		"jsoncompact":  "json-compact",
		"json":         "json-compact",
		"json-pretty":  "json-pretty",
		"json_pretty":  "json-pretty",
		"jsonpretty":   "json-pretty",
		"pretty-json":  "json-pretty",
		"yaml":         "yaml",
		"yml":          "yaml",
		"toml":         "toml",
		"xml":          "xml-compact",
		"xml-compact":  "xml-compact",
		"xml_compact":  "xml-compact",
		"xmlcompact":   "xml-compact",
		"xml-pretty":   "xml-pretty",
		"xml_pretty":   "xml-pretty",
		"xmlpretty":    "xml-pretty",
		"pretty-xml":   "xml-pretty",
		"messagepack":  "messagepack",
		"msgpack":      "messagepack",
		"msg-pack":     "messagepack",
		"message-pack": "messagepack",
		"message_pack": "messagepack",
	}
	if normalized, ok := aliases[key]; ok {
		return normalized, nil
	}
	return "", fmt.Errorf("invalid baseline %q; allowed values: toon|json-compact|json-pretty|yaml|toml|xml-compact|xml-pretty|messagepack", raw)
}
