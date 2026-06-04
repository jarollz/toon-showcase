package main

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	toon "github.com/toon-format/toon-go"
	"github.com/vmihailenco/msgpack/v5"
	"go.yaml.in/yaml/v3"
)

type TextEntry struct {
	Key   string `json:"key" yaml:"key" toml:"key" msgpack:"key" xml:"key" toon:"key"`
	Value string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type Address struct {
	Line1      string `json:"line1" yaml:"line1" toml:"line1" msgpack:"line1" xml:"line1" toon:"line1"`
	City       string `json:"city" yaml:"city" toml:"city" msgpack:"city" xml:"city" toon:"city"`
	Region     string `json:"region" yaml:"region" toml:"region" msgpack:"region" xml:"region" toon:"region"`
	PostalCode string `json:"postalCode" yaml:"postalCode" toml:"postalCode" msgpack:"postalCode" xml:"postalCode" toon:"postalCode"`
	Country    string `json:"country" yaml:"country" toml:"country" msgpack:"country" xml:"country" toon:"country"`
}

type Customer struct {
	ID          string   `json:"id" yaml:"id" toml:"id" msgpack:"id" xml:"id" toon:"id"`
	DisplayName string   `json:"displayName" yaml:"displayName" toml:"displayName" msgpack:"displayName" xml:"displayName" toon:"displayName"`
	Email       string   `json:"email" yaml:"email" toml:"email" msgpack:"email" xml:"email" toon:"email"`
	Locale      string   `json:"locale" yaml:"locale" toml:"locale" msgpack:"locale" xml:"locale" toon:"locale"`
	Address     Address  `json:"address" yaml:"address" toml:"address" msgpack:"address" xml:"address" toon:"address"`
	Tags        []string `json:"tags" yaml:"tags" toml:"tags" msgpack:"tags" xml:"tags>tag" toon:"tags"`
}

type LineItem struct {
	SKU      string   `json:"sku" yaml:"sku" toml:"sku" msgpack:"sku" xml:"sku" toon:"sku"`
	Title    string   `json:"title" yaml:"title" toml:"title" msgpack:"title" xml:"title" toon:"title"`
	Quantity int      `json:"quantity" yaml:"quantity" toml:"quantity" msgpack:"quantity" xml:"quantity" toon:"quantity"`
	UnitUSD  float64  `json:"unitUsd" yaml:"unitUsd" toml:"unitUsd" msgpack:"unitUsd" xml:"unitUsd" toon:"unitUsd"`
	Labels   []string `json:"labels" yaml:"labels" toml:"labels" msgpack:"labels" xml:"labels>label" toon:"labels"`
}

type Shipment struct {
	Carrier    string `json:"carrier" yaml:"carrier" toml:"carrier" msgpack:"carrier" xml:"carrier" toon:"carrier"`
	TrackingID string `json:"trackingId" yaml:"trackingId" toml:"trackingId" msgpack:"trackingId" xml:"trackingId" toon:"trackingId"`
	ShippedAt  string `json:"shippedAt" yaml:"shippedAt" toml:"shippedAt" msgpack:"shippedAt" xml:"shippedAt" toon:"shippedAt"`
	Delivered  bool   `json:"delivered" yaml:"delivered" toml:"delivered" msgpack:"delivered" xml:"delivered" toon:"delivered"`
}

type AuditEvent struct {
	At      string `json:"at" yaml:"at" toml:"at" msgpack:"at" xml:"at" toon:"at"`
	Actor   string `json:"actor" yaml:"actor" toml:"actor" msgpack:"actor" xml:"actor" toon:"actor"`
	Action  string `json:"action" yaml:"action" toml:"action" msgpack:"action" xml:"action" toon:"action"`
	Message string `json:"message" yaml:"message" toml:"message" msgpack:"message" xml:"message" toon:"message"`
}

type Discount struct {
	Code       string  `json:"code" yaml:"code" toml:"code" msgpack:"code" xml:"code" toon:"code"`
	Percentage float64 `json:"percentage" yaml:"percentage" toml:"percentage" msgpack:"percentage" xml:"percentage" toon:"percentage"`
}

type LocalizedNote struct {
	Lang string `json:"lang" yaml:"lang" toml:"lang" msgpack:"lang" xml:"lang" toon:"lang"`
	Text string `json:"text" yaml:"text" toml:"text" msgpack:"text" xml:"text" toon:"text"`
}

type Order struct {
	OrderID          string          `json:"orderId" yaml:"orderId" toml:"orderId" msgpack:"orderId" xml:"orderId" toon:"orderId"`
	CreatedAt        string          `json:"createdAt" yaml:"createdAt" toml:"createdAt" msgpack:"createdAt" xml:"createdAt" toon:"createdAt"`
	Customer         Customer        `json:"customer" yaml:"customer" toml:"customer" msgpack:"customer" xml:"customer" toon:"customer"`
	Items            []LineItem      `json:"items" yaml:"items" toml:"items" msgpack:"items" xml:"items>item" toon:"items"`
	Shipments        []Shipment      `json:"shipments" yaml:"shipments" toml:"shipments" msgpack:"shipments" xml:"shipments>shipment" toon:"shipments"`
	Discount         *Discount       `json:"discount,omitempty" yaml:"discount,omitempty" toml:"discount,omitempty" msgpack:"discount,omitempty" xml:"discount,omitempty" toon:"discount,omitempty"`
	IsGift           bool            `json:"isGift" yaml:"isGift" toml:"isGift" msgpack:"isGift" xml:"isGift" toon:"isGift"`
	Priority         int             `json:"priority" yaml:"priority" toml:"priority" msgpack:"priority" xml:"priority" toon:"priority"`
	TotalUSD         float64         `json:"totalUsd" yaml:"totalUsd" toml:"totalUsd" msgpack:"totalUsd" xml:"totalUsd" toon:"totalUsd"`
	LocalizedNotes   []LocalizedNote `json:"localizedNotes" yaml:"localizedNotes" toml:"localizedNotes" msgpack:"localizedNotes" xml:"localizedNotes>note" toon:"localizedNotes"`
	EmojiStatus      string          `json:"emojiStatus" yaml:"emojiStatus" toml:"emojiStatus" msgpack:"emojiStatus" xml:"emojiStatus" toon:"emojiStatus"`
	SpecialInstr     *string         `json:"specialInstr,omitempty" yaml:"specialInstr,omitempty" toml:"specialInstr,omitempty" msgpack:"specialInstr,omitempty" xml:"specialInstr,omitempty" toon:"specialInstr,omitempty"`
	AuditTrail       []AuditEvent    `json:"auditTrail" yaml:"auditTrail" toml:"auditTrail" msgpack:"auditTrail" xml:"auditTrail>event" toon:"auditTrail"`
	FulfillmentCodes []string        `json:"fulfillmentCodes" yaml:"fulfillmentCodes" toml:"fulfillmentCodes" msgpack:"fulfillmentCodes" xml:"fulfillmentCodes>code" toon:"fulfillmentCodes"`
}

type OrderSummary struct {
	OrderID     string  `json:"orderId" yaml:"orderId" toml:"orderId" msgpack:"orderId" xml:"orderId" toon:"orderId"`
	CustomerID  string  `json:"customerId" yaml:"customerId" toml:"customerId" msgpack:"customerId" xml:"customerId" toon:"customerId"`
	Locale      string  `json:"locale" yaml:"locale" toml:"locale" msgpack:"locale" xml:"locale" toon:"locale"`
	Priority    int     `json:"priority" yaml:"priority" toml:"priority" msgpack:"priority" xml:"priority" toon:"priority"`
	TotalUSD    float64 `json:"totalUsd" yaml:"totalUsd" toml:"totalUsd" msgpack:"totalUsd" xml:"totalUsd" toon:"totalUsd"`
	EmojiStatus string  `json:"emojiStatus" yaml:"emojiStatus" toml:"emojiStatus" msgpack:"emojiStatus" xml:"emojiStatus" toon:"emojiStatus"`
}

type BenchmarkData struct {
	XMLName       xml.Name         `json:"-" yaml:"-" toml:"-" msgpack:"-" xml:"showcase" toon:"-"`
	Version       string           `json:"version" yaml:"version" toml:"version" msgpack:"version" xml:"version" toon:"version"`
	GeneratedAt   string           `json:"generatedAt" yaml:"generatedAt" toml:"generatedAt" msgpack:"generatedAt" xml:"generatedAt" toon:"generatedAt"`
	Owner         string           `json:"owner" yaml:"owner" toml:"owner" msgpack:"owner" xml:"owner" toon:"owner"`
	DefaultLocale string           `json:"defaultLocale" yaml:"defaultLocale" toml:"defaultLocale" msgpack:"defaultLocale" xml:"defaultLocale" toon:"defaultLocale"`
	Meta          []TextEntry      `json:"meta" yaml:"meta" toml:"meta" msgpack:"meta" xml:"meta>entry" toon:"meta"`
	OrdersByID    map[string]Order `json:"ordersById" yaml:"ordersById" toml:"ordersById" msgpack:"ordersById" xml:"-" toon:"ordersById"`
	OrderTable    []OrderSummary   `json:"orderTable" yaml:"orderTable" toml:"orderTable" msgpack:"orderTable" xml:"orderTable>summary" toon:"orderTable"`
}

type xmlBenchmarkData struct {
	XMLName       xml.Name       `xml:"showcase"`
	Version       string         `xml:"version"`
	GeneratedAt   string         `xml:"generatedAt"`
	Owner         string         `xml:"owner"`
	DefaultLocale string         `xml:"defaultLocale"`
	Meta          []TextEntry    `xml:"meta>entry"`
	Orders        []Order        `xml:"orders>order"`
	OrderTable    []OrderSummary `xml:"orderTable>summary"`
}

type formatDef struct {
	Key          string
	Name         string
	Binary       bool
	EncodeData   func(BenchmarkData) ([]byte, error)
	DecodeData   func([]byte) (BenchmarkData, error)
	MarshalAny   func(any) ([]byte, error)
	UnmarshalAny func([]byte, any) error
}

type benchmarkResult struct {
	Key              string
	Name             string
	Binary           bool
	MarshalTotal     time.Duration
	MarshalAvg       time.Duration
	UnmarshalTotal   time.Duration
	UnmarshalAvg     time.Duration
	EncodedBytes     int
	EncodedChars     int
	RoundTripOK      bool
	RoundTripDetails string
	Sample           string
}

type charsetCase struct {
	Name  string
	Value string
}

type charsetResult struct {
	MarshalOK   bool
	RoundTripOK bool
	Detail      string
}

type charsetMatrix struct {
	Tests           []charsetCase
	FormatNames     []string
	ResultsByFormat [][]charsetResult
	PassCounts      []int
}

type shapeCase struct {
	Name    string
	Prepare func() (any, any)
}

type shapeResult struct {
	EncodeOK    bool
	DecodeOK    bool
	RoundTripOK bool
	Detail      string
}

type shapeMatrix struct {
	Tests           []shapeCase
	FormatNames     []string
	ResultsByFormat [][]shapeResult
	PassCounts      []int
}

type probeSliceStruct struct {
	Value []OrderSummary `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value>item" toon:"value"`
}

type probeSliceMapString struct {
	Value []map[string]string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeMapToSliceStruct struct {
	Value map[string][]OrderSummary `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeSliceAny struct {
	Value []any `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeSlicePointer struct {
	Value []*OrderSummary `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value>item" toon:"value"`
}

type probeNestedSliceMap struct {
	Value [][]map[string]string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type probeMapNonStringKey struct {
	Value map[int]string `json:"value" yaml:"value" toml:"value" msgpack:"value" xml:"value" toon:"value"`
}

type charsetProbe struct {
	Owner string      `json:"owner" yaml:"owner" toml:"owner" msgpack:"owner" xml:"owner" toon:"owner"`
	Meta  []TextEntry `json:"meta" yaml:"meta" toml:"meta" msgpack:"meta" xml:"meta>entry" toon:"meta"`
}

type progressBar struct {
	total   int64
	done    int64
	stageMu sync.RWMutex
	stage   string
	stopCh  chan struct{}
	stopped chan struct{}
	started bool
	startMu sync.Mutex
}

func main() {
	records := flag.Int("records", 120, "number of orders in dataset")
	iters := flag.Int("iters", 500, "benchmark iterations for each marshal and unmarshal")
	warmup := flag.Int("warmup", 50, "warmup iterations before measuring")
	seed := flag.Int64("seed", -1, "seed for dataset generation; use -1 for current unix timestamp")
	indentSize := flag.Int("indent", 2, "explicit indent size for JSON pretty, TOON, YAML, and TOML")
	sampleLimit := flag.Int("sample-limit", 180, "max characters shown for encoded sample")
	noProgress := flag.Bool("no-progress", false, "disable animated progress bar")
	baseline := flag.String("baseline", "toon", "baseline format for ratio/delta: toon|json-compact|json-pretty|yaml|toml|xml-compact|xml-pretty|messagepack")
	csvMode := flag.String("csv", "off", "csv output mode: off|stdout|file")
	csvFile := flag.String("csv-file", "benchmark_report.csv", "csv output path when -csv=file")
	flag.Parse()

	if *records <= 0 || *iters <= 0 || *warmup < 0 || *indentSize <= 0 {
		panic("records and iters must be > 0; warmup >= 0; indent > 0")
	}
	if *csvMode != "off" && *csvMode != "stdout" && *csvMode != "file" {
		panic("csv mode must be one of: off, stdout, file")
	}

	resolvedSeed := *seed
	if resolvedSeed == -1 {
		resolvedSeed = time.Now().Unix()
	}

	dataset := generateDataset(*records, resolvedSeed)
	formats := buildFormats(*indentSize)
	resolvedBaselineKey, err := normalizeBaselineKey(*baseline)
	if err != nil {
		panic(err.Error())
	}

	var progress *progressBar
	if !*noProgress {
		progress = newProgressBar(totalBenchmarkSteps(len(formats), *iters, *warmup))
		progress.SetStage("initializing benchmark")
		progress.Start()
	}

	results := make([]benchmarkResult, 0, len(formats))
	for _, f := range formats {
		if progress != nil {
			progress.SetStage("benchmark " + f.Name)
		}
		res, err := runBenchmark(f, dataset, *iters, *warmup, *sampleLimit, progress)
		if err != nil {
			if progress != nil {
				progress.Finish()
			}
			panic(fmt.Sprintf("%s failed: %v", f.Name, err))
		}
		results = append(results, res)
	}
	if progress != nil {
		progress.Finish()
	}
	baselineRes, err := findBaselineResult(results, resolvedBaselineKey)
	if err != nil {
		panic(err.Error())
	}

	printBenchmarkReport(results, baselineRes, *records, *iters, *warmup, resolvedSeed, *indentSize)
	charsetReport := runCharsetMatrix(formats)
	printCharsetMatrix(charsetReport)
	shapeReport := runShapeMatrix(formats)
	printShapeMatrix(shapeReport)

	if *csvMode != "off" {
		benchmarkCSV, err := buildBenchmarkCSV(results, baselineRes, *records, *iters, *warmup, resolvedSeed, *indentSize)
		if err != nil {
			panic(fmt.Sprintf("failed to build benchmark csv: %v", err))
		}
		charsetCSV, err := buildCharsetCSV(charsetReport)
		if err != nil {
			panic(fmt.Sprintf("failed to build charset csv: %v", err))
		}
		shapeCSV, err := buildShapeCSV(shapeReport)
		if err != nil {
			panic(fmt.Sprintf("failed to build shape csv: %v", err))
		}
		emitCSVs(*csvMode, *csvFile, benchmarkCSV, charsetCSV, shapeCSV)
	}
}

func totalBenchmarkSteps(formatCount, iters, warmup int) int64 {
	if formatCount <= 0 {
		return 1
	}
	perFormat := int64(warmup*2 + iters*2 + 1)
	total := int64(formatCount) * perFormat
	if total <= 0 {
		return 1
	}
	return total
}

func newProgressBar(total int64) *progressBar {
	if total <= 0 {
		total = 1
	}
	return &progressBar{
		total:   total,
		stopCh:  make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

func (p *progressBar) SetStage(stage string) {
	p.stageMu.Lock()
	p.stage = stage
	p.stageMu.Unlock()
}

func (p *progressBar) Advance(n int64) {
	if n <= 0 {
		return
	}
	atomic.AddInt64(&p.done, n)
}

func (p *progressBar) Start() {
	p.startMu.Lock()
	defer p.startMu.Unlock()
	if p.started {
		return
	}
	p.started = true
	go p.renderLoop()
}

func (p *progressBar) Finish() {
	p.startMu.Lock()
	started := p.started
	p.startMu.Unlock()
	if !started {
		return
	}
	atomic.StoreInt64(&p.done, p.total)
	close(p.stopCh)
	<-p.stopped
}

func (p *progressBar) renderLoop() {
	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()
	defer close(p.stopped)

	for {
		select {
		case <-ticker.C:
			p.render(false)
		case <-p.stopCh:
			p.render(true)
			return
		}
	}
}

func (p *progressBar) render(final bool) {
	done := atomic.LoadInt64(&p.done)
	total := p.total
	if done > total {
		done = total
	}
	percent := (float64(done) / float64(total)) * 100
	if total == 0 {
		percent = 100
	}
	if final {
		percent = 100
		done = total
	}

	p.stageMu.RLock()
	stage := p.stage
	p.stageMu.RUnlock()

	const barWidth = 30
	filled := int((float64(done) / float64(total)) * barWidth)
	if final {
		filled = barWidth
	}
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)

	if stage == "" {
		stage = "processing"
	}

	fmt.Fprintf(os.Stderr, "\rProgress [%s] %6.2f%% (%d/%d) %s", bar, percent, done, total, stage)
	if final {
		fmt.Fprintln(os.Stderr)
	}
}

func normalizeBaselineKey(raw string) (string, error) {
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

func findBaselineResult(results []benchmarkResult, baselineKey string) (benchmarkResult, error) {
	for _, res := range results {
		if res.Key == baselineKey {
			return res, nil
		}
	}
	return benchmarkResult{}, fmt.Errorf("baseline key %q not found in benchmark results", baselineKey)
}

func buildFormats(indentSize int) []formatDef {
	indent := strings.Repeat(" ", indentSize)

	marshalYAML := func(v any) ([]byte, error) {
		var buf bytes.Buffer
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(indentSize)
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
		if err := enc.Close(); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	marshalTOML := func(v any) ([]byte, error) {
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.Indent = indent
		if err := enc.Encode(v); err != nil {
			return nil, err
		}
		return bytes.TrimRight(buf.Bytes(), "\n"), nil
	}

	xmlIndent := indent
	marshalXMLPretty := func(data BenchmarkData) ([]byte, error) {
		x := xmlBenchmarkData{
			Version:       data.Version,
			GeneratedAt:   data.GeneratedAt,
			Owner:         data.Owner,
			DefaultLocale: data.DefaultLocale,
			Meta:          data.Meta,
			OrderTable:    data.OrderTable,
			Orders:        orderedOrders(data.OrdersByID),
		}
		return xml.MarshalIndent(x, "", xmlIndent)
	}

	return []formatDef{
		{
			Key:    "json-compact",
			Name:   "JSON compact",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return json.Marshal(data)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				var out BenchmarkData
				err := json.Unmarshal(raw, &out)
				return out, err
			},
			MarshalAny: json.Marshal,
			UnmarshalAny: func(raw []byte, out any) error {
				return json.Unmarshal(raw, out)
			},
		},
		{
			Key:    "json-pretty",
			Name:   "JSON pretty",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return json.MarshalIndent(data, "", indent)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				var out BenchmarkData
				err := json.Unmarshal(raw, &out)
				return out, err
			},
			MarshalAny: func(v any) ([]byte, error) { return json.MarshalIndent(v, "", indent) },
			UnmarshalAny: func(raw []byte, out any) error {
				return json.Unmarshal(raw, out)
			},
		},
		{
			Key:    "toon",
			Name:   "TOON",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return toon.Marshal(data, toon.WithIndent(indentSize))
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				var out BenchmarkData
				err := toon.Unmarshal(raw, &out)
				return out, err
			},
			MarshalAny: func(v any) ([]byte, error) { return toon.Marshal(v, toon.WithIndent(indentSize)) },
			UnmarshalAny: func(raw []byte, out any) error {
				return toon.Unmarshal(raw, out)
			},
		},
		{
			Key:    "yaml",
			Name:   "YAML",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return marshalYAML(data)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				var out BenchmarkData
				err := yaml.Unmarshal(raw, &out)
				return out, err
			},
			MarshalAny: marshalYAML,
			UnmarshalAny: func(raw []byte, out any) error {
				return yaml.Unmarshal(raw, out)
			},
		},
		{
			Key:    "toml",
			Name:   "TOML",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return marshalTOML(data)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				var out BenchmarkData
				err := toml.Unmarshal(raw, &out)
				return out, err
			},
			MarshalAny: marshalTOML,
			UnmarshalAny: func(raw []byte, out any) error {
				return toml.Unmarshal(raw, out)
			},
		},
		{
			Key:    "xml-compact",
			Name:   "XML compact",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return marshalXMLBenchmark(data)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				out, err := unmarshalXMLBenchmark(raw)
				return out, err
			},
			MarshalAny: xml.Marshal,
			UnmarshalAny: func(raw []byte, out any) error {
				return xml.Unmarshal(raw, out)
			},
		},
		{
			Key:    "xml-pretty",
			Name:   "XML pretty",
			Binary: false,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return marshalXMLPretty(data)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				out, err := unmarshalXMLBenchmark(raw)
				return out, err
			},
			MarshalAny: func(v any) ([]byte, error) {
				x, ok := v.(BenchmarkData)
				if ok {
					return marshalXMLPretty(x)
				}
				return xml.MarshalIndent(v, "", xmlIndent)
			},
			UnmarshalAny: func(raw []byte, out any) error {
				return xml.Unmarshal(raw, out)
			},
		},
		{
			Key:    "messagepack",
			Name:   "MessagePack",
			Binary: true,
			EncodeData: func(data BenchmarkData) ([]byte, error) {
				return msgpack.Marshal(data)
			},
			DecodeData: func(raw []byte) (BenchmarkData, error) {
				var out BenchmarkData
				err := msgpack.Unmarshal(raw, &out)
				return out, err
			},
			MarshalAny: msgpack.Marshal,
			UnmarshalAny: func(raw []byte, out any) error {
				return msgpack.Unmarshal(raw, out)
			},
		},
	}
}

func runBenchmark(def formatDef, data BenchmarkData, iters, warmup, sampleLimit int, progress *progressBar) (benchmarkResult, error) {
	for i := 0; i < warmup; i++ {
		raw, err := def.EncodeData(data)
		if err != nil {
			return benchmarkResult{}, fmt.Errorf("warmup marshal: %w", err)
		}
		if progress != nil {
			progress.Advance(1)
		}
		if _, err := def.DecodeData(raw); err != nil {
			return benchmarkResult{}, fmt.Errorf("warmup unmarshal: %w; encoded sample=%s", err, formatSample(raw, def.Binary, 4000))
		}
		if progress != nil {
			progress.Advance(1)
		}
	}

	marshalStart := time.Now()
	var encoded []byte
	for i := 0; i < iters; i++ {
		raw, err := def.EncodeData(data)
		if err != nil {
			return benchmarkResult{}, fmt.Errorf("marshal iteration %d: %w", i, err)
		}
		encoded = raw
		if progress != nil {
			progress.Advance(1)
		}
	}
	marshalTotal := time.Since(marshalStart)

	unmarshalStart := time.Now()
	for i := 0; i < iters; i++ {
		if _, err := def.DecodeData(encoded); err != nil {
			return benchmarkResult{}, fmt.Errorf("unmarshal iteration %d: %w", i, err)
		}
		if progress != nil {
			progress.Advance(1)
		}
	}
	unmarshalTotal := time.Since(unmarshalStart)

	out, err := def.DecodeData(encoded)
	if err != nil {
		return benchmarkResult{}, fmt.Errorf("final decode for roundtrip: %w", err)
	}
	if progress != nil {
		progress.Advance(1)
	}

	encodedChars := utf8.RuneCount(encoded)
	if def.Binary {
		encodedChars = -1
	}

	res := benchmarkResult{
		Key:            def.Key,
		Name:           def.Name,
		Binary:         def.Binary,
		MarshalTotal:   marshalTotal,
		MarshalAvg:     time.Duration(int64(marshalTotal) / int64(iters)),
		UnmarshalTotal: unmarshalTotal,
		UnmarshalAvg:   time.Duration(int64(unmarshalTotal) / int64(iters)),
		EncodedBytes:   len(encoded),
		EncodedChars:   encodedChars,
		Sample:         formatSample(encoded, def.Binary, sampleLimit),
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

func marshalXMLBenchmark(data BenchmarkData) ([]byte, error) {
	x := xmlBenchmarkData{
		Version:       data.Version,
		GeneratedAt:   data.GeneratedAt,
		Owner:         data.Owner,
		DefaultLocale: data.DefaultLocale,
		Meta:          data.Meta,
		OrderTable:    data.OrderTable,
		Orders:        orderedOrders(data.OrdersByID),
	}
	return xml.Marshal(x)
}

func unmarshalXMLBenchmark(raw []byte) (BenchmarkData, error) {
	var x xmlBenchmarkData
	if err := xml.Unmarshal(raw, &x); err != nil {
		return BenchmarkData{}, err
	}
	ordersByID := make(map[string]Order, len(x.Orders))
	for _, order := range x.Orders {
		ordersByID[order.OrderID] = order
	}
	return BenchmarkData{
		Version:       x.Version,
		GeneratedAt:   x.GeneratedAt,
		Owner:         x.Owner,
		DefaultLocale: x.DefaultLocale,
		Meta:          x.Meta,
		OrdersByID:    ordersByID,
		OrderTable:    x.OrderTable,
	}, nil
}

func orderedOrders(m map[string]Order) []Order {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	orders := make([]Order, 0, len(keys))
	for _, k := range keys {
		orders = append(orders, m[k])
	}
	return orders
}

func datasetsEqual(a, b BenchmarkData) bool {
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

func generateDataset(records int, seed int64) BenchmarkData {
	r := rand.New(rand.NewSource(seed))
	base := time.Date(2026, time.June, 4, 9, 30, 0, 0, time.UTC)

	owners := []string{"Platform Team", "Backend Guild", "開発チーム", "数据平台组"}
	cities := []string{"Tokyo", "上海", "Singapore", "San Francisco"}
	regions := []string{"JP-13", "CN-SH", "SG-01", "US-CA"}
	locales := []string{"en-US", "ja-JP", "zh-CN"}
	carriers := []string{"DHL", "FedEx", "Yamato", "顺丰"}
	emojiStatus := []string{"🚀 shipped", "📦 packing", "✅ delivered", "🛠️ processing"}
	englishTitles := []string{"Notebook", "Wireless Mouse", "Mechanical Keyboard"}
	japaneseTitles := []string{"抹茶セット", "和風ランプ", "工芸マグ"}
	chineseTitles := []string{"竹制茶盘", "丝绸围巾", "手工陶杯"}

	ordersByID := make(map[string]Order, records)
	orderTable := make([]OrderSummary, 0, records)

	for i := 0; i < records; i++ {
		loc := locales[i%len(locales)]
		city := cities[i%len(cities)]
		region := regions[i%len(regions)]

		item1Qty := 1 + r.Intn(3)
		item2Qty := 1 + r.Intn(2)
		item3Qty := 1 + r.Intn(4)

		item1Price := 12.5 + float64(r.Intn(300))/10
		item2Price := 15.5 + float64(r.Intn(250))/10
		item3Price := 8.9 + float64(r.Intn(200))/10

		items := []LineItem{
			{SKU: fmt.Sprintf("ENG-%03d", i), Title: englishTitles[i%len(englishTitles)], Quantity: item1Qty, UnitUSD: item1Price, Labels: []string{"english", "office", "stable"}},
			{SKU: fmt.Sprintf("JPN-%03d", i), Title: japaneseTitles[i%len(japaneseTitles)], Quantity: item2Qty, UnitUSD: item2Price, Labels: []string{"日本語", "craft", "gift"}},
			{SKU: fmt.Sprintf("CHN-%03d", i), Title: chineseTitles[i%len(chineseTitles)], Quantity: item3Qty, UnitUSD: item3Price, Labels: []string{"中文", "home", "limited"}},
		}

		total := float64(item1Qty)*item1Price + float64(item2Qty)*item2Price + float64(item3Qty)*item3Price

		created := base.Add(time.Duration(i) * time.Hour)
		createdText := created.Format(time.RFC3339Nano)
		shippedText := created.Add(3 * time.Hour).Format(time.RFC3339Nano)

		localizedNotes := []LocalizedNote{
			{Lang: "en", Text: "Ship fast; include invoice copy."},
			{Lang: "ja", Text: "丁寧に梱包してください。"},
			{Lang: "zh", Text: "请小心包装并附上发票。"},
			{Lang: "emoji", Text: "🚀✨🙂"},
		}

		audit := []AuditEvent{
			{At: createdText, Actor: "system", Action: "created", Message: "Order created via API"},
			{At: created.Add(20 * time.Minute).Format(time.RFC3339Nano), Actor: "qa.bot", Action: "validated", Message: "UTF-8 fields validated: 你好 / こんにちは / hello"},
			{At: shippedText, Actor: "ops", Action: "shipped", Message: "Tracking assigned 📦"},
		}

		var discount *Discount
		if i%3 == 0 {
			discount = &Discount{Code: fmt.Sprintf("SAVE-%02d", i%20), Percentage: 5 + float64(i%10)}
		}

		var special *string
		if i%2 == 0 {
			msg := fmt.Sprintf("Gift wrap: yes; note=%s / 注文 / 订单", emojiStatus[i%len(emojiStatus)])
			special = &msg
		}

		order := Order{
			OrderID:   fmt.Sprintf("ORD-%06d", i+1),
			CreatedAt: createdText,
			Customer: Customer{
				ID:          fmt.Sprintf("CUS-%05d", i+1000),
				DisplayName: fmt.Sprintf("User %d 你好 こんにちは", i+1),
				Email:       fmt.Sprintf("user%03d@example.com", i+1),
				Locale:      loc,
				Address: Address{
					Line1:      fmt.Sprintf("%d Market Street", 100+i),
					City:       city,
					Region:     region,
					PostalCode: fmt.Sprintf("%05d", 10000+i),
					Country:    "Global",
				},
				Tags: []string{"vip", "newsletter", "emoji-🙂"},
			},
			Items: items,
			Shipments: []Shipment{
				{Carrier: carriers[i%len(carriers)], TrackingID: fmt.Sprintf("TRK-%08d", 200000+i), ShippedAt: shippedText, Delivered: i%4 == 0},
			},
			Discount:         discount,
			IsGift:           i%2 == 0,
			Priority:         i%5 + 1,
			TotalUSD:         round2(total),
			LocalizedNotes:   localizedNotes,
			EmojiStatus:      emojiStatus[i%len(emojiStatus)],
			SpecialInstr:     special,
			AuditTrail:       audit,
			FulfillmentCodes: []string{"A-1", "B_2", "C.3"},
		}

		ordersByID[order.OrderID] = order
		orderTable = append(orderTable, OrderSummary{
			OrderID:     order.OrderID,
			CustomerID:  order.Customer.ID,
			Locale:      order.Customer.Locale,
			Priority:    order.Priority,
			TotalUSD:    order.TotalUSD,
			EmojiStatus: order.EmojiStatus,
		})
	}

	meta := []TextEntry{
		{Key: "project", Value: "TOON vs Multi-format showcase"},
		{Key: "languages", Value: "English / 日本語 / 中文"},
		{Key: "emoji", Value: "🚀✨🙂"},
		{Key: "explicitIndent", Value: "configured"},
	}

	return BenchmarkData{
		Version:       "1.0.0",
		GeneratedAt:   base.Format(time.RFC3339Nano),
		Owner:         owners[int(seed)%len(owners)],
		DefaultLocale: "en-US",
		Meta:          meta,
		OrdersByID:    ordersByID,
		OrderTable:    orderTable,
	}
}

func round2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func printBenchmarkReport(results []benchmarkResult, baseline benchmarkResult, records, iters, warmup int, seed int64, indentSize int) {
	fmt.Println()
	fmt.Println("=== Multi-Format Marshal/Unmarshal Showcase ===")
	fmt.Printf("Dataset records: %d\n", records)
	fmt.Printf("Iterations: %d (warmup: %d)\n", iters, warmup)
	fmt.Printf("Seed: %d\n", seed)
	fmt.Printf("Indent length used for JSON pretty + TOON + YAML + TOML + XML pretty: %d spaces\n", indentSize)
	fmt.Println()

	fmt.Printf("%-13s  %-13s  %-13s  %-13s  %-13s  %-8s  %-8s  %-10s\n",
		"Format", "Marshal(ms)", "Marshal(ns)", "Unmarshal(ms)", "Unmarshal(ns)", "Bytes", "Chars", "RoundTrip")
	for _, res := range results {
		chars := fmt.Sprintf("%d", res.EncodedChars)
		if res.EncodedChars < 0 {
			chars = "n/a"
		}
		fmt.Printf("%-13s  %-13.2f  %-13d  %-13.2f  %-13d  %-8d  %-8s  %-10s\n",
			res.Name,
			float64(res.MarshalTotal.Microseconds())/1000.0,
			res.MarshalAvg.Nanoseconds(),
			float64(res.UnmarshalTotal.Microseconds())/1000.0,
			res.UnmarshalAvg.Nanoseconds(),
			res.EncodedBytes,
			chars,
			res.RoundTripDetails,
		)
	}

	fmt.Println()
	fmt.Printf("Ratios vs selected baseline (lower is better):\n")
	fmt.Printf("Baseline format: %s\n", baseline.Name)
	fmt.Printf("%-13s  %-16s  %-16s  %-16s  %-16s\n", "Format", "MarshalRatio", "UnmarshalRatio", "ByteRatio", "CharRatio")
	for _, res := range results {
		fmt.Printf("%-13s  %-16s  %-16s  %-16s  %-16s\n",
			res.Name,
			ratio(res.MarshalAvg, baseline.MarshalAvg),
			ratio(res.UnmarshalAvg, baseline.UnmarshalAvg),
			ratioInt(res.EncodedBytes, baseline.EncodedBytes),
			ratioInt(res.EncodedChars, baseline.EncodedChars),
		)
	}

	fmt.Println()
	fmt.Printf("Human-friendly delta vs %s baseline:\n", baseline.Name)
	fmt.Println("(negative direction for time is faster; negative direction for size is smaller)")
	fmt.Printf("%-13s  %-26s  %-26s  %-24s  %-24s\n", "Format", "Marshal", "Unmarshal", "Bytes", "Chars")
	for _, res := range results {
		fmt.Printf("%-13s  %-26s  %-26s  %-24s  %-24s\n",
			res.Name,
			deltaDurationText(res.MarshalAvg, baseline.MarshalAvg),
			deltaDurationText(res.UnmarshalAvg, baseline.UnmarshalAvg),
			deltaIntText(res.EncodedBytes, baseline.EncodedBytes),
			deltaIntText(res.EncodedChars, baseline.EncodedChars),
		)
	}

	fmt.Println()
	fmt.Println("Encoded output sample (truncated):")
	for _, res := range results {
		fmt.Printf("- %s: %s\n", res.Name, res.Sample)
	}
	fmt.Println()
	fmt.Println("Notes:")
	fmt.Println("- MessagePack is binary; char metrics shown as n/a.")
	fmt.Println("- XML is included in two lanes: compact and pretty.")
	fmt.Println()
}

func ratio(v, base time.Duration) string {
	if base == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.3fx", float64(v)/float64(base))
}

func ratioInt(v, base int) string {
	if base <= 0 || v < 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.3fx", float64(v)/float64(base))
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

func runCharsetMatrix(formats []formatDef) charsetMatrix {
	tests := []charsetCase{
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

	resultsByFormat := make([][]charsetResult, len(formats))
	passCounts := make([]int, len(formats))
	formatNames := make([]string, len(formats))

	for fi, f := range formats {
		formatNames[fi] = f.Name
		rows := make([]charsetResult, 0, len(tests))
		for _, tc := range tests {
			res := evalCharsetCase(f, tc)
			if res.MarshalOK && res.RoundTripOK {
				passCounts[fi]++
			}
			rows = append(rows, res)
		}
		resultsByFormat[fi] = rows
	}

	return charsetMatrix{Tests: tests, FormatNames: formatNames, ResultsByFormat: resultsByFormat, PassCounts: passCounts}
}

func printCharsetMatrix(report charsetMatrix) {
	fmt.Println("=== Character Set Support Comparison ===")
	fmt.Println("Scope: practical behavior in selected Go libraries")
	fmt.Println()

	fmt.Printf("%-17s", "Case")
	for _, name := range report.FormatNames {
		fmt.Printf("  %-16s", name)
	}
	fmt.Println()

	for i, tc := range report.Tests {
		fmt.Printf("%-17s", tc.Name)
		for fi := range report.FormatNames {
			cell := "FAIL"
			if report.ResultsByFormat[fi][i].MarshalOK && report.ResultsByFormat[fi][i].RoundTripOK {
				cell = "OK"
			}
			fmt.Printf("  %-16s", cell)
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Support summary:")
	for fi, name := range report.FormatNames {
		fmt.Printf("- %s: %d/%d cases roundtrip\n", name, report.PassCounts[fi], len(report.Tests))
	}

	fmt.Println()
	fmt.Println("Notable details:")
	for fi, name := range report.FormatNames {
		for ti, tc := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			if res.MarshalOK && res.RoundTripOK {
				continue
			}
			fmt.Printf("- %s / %s: %s\n", name, tc.Name, res.Detail)
		}
	}
	fmt.Println()
}

func evalCharsetCase(def formatDef, tc charsetCase) charsetResult {
	probe := charsetProbe{Owner: tc.Value, Meta: []TextEntry{{Key: "probe", Value: tc.Value}}}

	raw, err := def.MarshalAny(probe)
	if err != nil {
		return charsetResult{MarshalOK: false, RoundTripOK: false, Detail: err.Error()}
	}
	var out charsetProbe
	if err := def.UnmarshalAny(raw, &out); err != nil {
		return charsetResult{MarshalOK: true, RoundTripOK: false, Detail: err.Error()}
	}
	if out.Owner != probe.Owner || len(out.Meta) == 0 || out.Meta[0].Value != probe.Meta[0].Value {
		return charsetResult{MarshalOK: true, RoundTripOK: false, Detail: "decoded value differs"}
	}
	return charsetResult{MarshalOK: true, RoundTripOK: true, Detail: "ok"}
}

func runShapeMatrix(formats []formatDef) shapeMatrix {
	tests := []shapeCase{
		{
			Name: "slice_struct",
			Prepare: func() (any, any) {
				in := probeSliceStruct{Value: []OrderSummary{{OrderID: "A", CustomerID: "C1", Locale: "en-US", Priority: 1, TotalUSD: 11.2, EmojiStatus: "🚀"}}}
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
				in := probeMapToSliceStruct{Value: map[string][]OrderSummary{"group": {{OrderID: "A", CustomerID: "X", Locale: "ja-JP", Priority: 2, TotalUSD: 33.4, EmojiStatus: "✨"}}}}
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
				v := &OrderSummary{OrderID: "X", CustomerID: "C", Locale: "zh-CN", Priority: 1, TotalUSD: 9.9, EmojiStatus: "🙂"}
				in := probeSlicePointer{Value: []*OrderSummary{v, nil}}
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

	resultsByFormat := make([][]shapeResult, len(formats))
	passCounts := make([]int, len(formats))
	formatNames := make([]string, len(formats))

	for fi, f := range formats {
		formatNames[fi] = f.Name
		rows := make([]shapeResult, 0, len(tests))
		for _, tc := range tests {
			res := evalShapeCase(f, tc)
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				passCounts[fi]++
			}
			rows = append(rows, res)
		}
		resultsByFormat[fi] = rows
	}

	return shapeMatrix{Tests: tests, FormatNames: formatNames, ResultsByFormat: resultsByFormat, PassCounts: passCounts}
}

func evalShapeCase(def formatDef, tc shapeCase) shapeResult {
	in, outPtr := tc.Prepare()
	raw, err := def.MarshalAny(in)
	if err != nil {
		return shapeResult{EncodeOK: false, DecodeOK: false, RoundTripOK: false, Detail: err.Error()}
	}
	if err := def.UnmarshalAny(raw, outPtr); err != nil {
		return shapeResult{EncodeOK: true, DecodeOK: false, RoundTripOK: false, Detail: err.Error()}
	}
	outVal := reflect.ValueOf(outPtr).Elem().Interface()
	if !normalizedEqual(in, outVal) {
		return shapeResult{EncodeOK: true, DecodeOK: true, RoundTripOK: false, Detail: "decoded value differs"}
	}
	return shapeResult{EncodeOK: true, DecodeOK: true, RoundTripOK: true, Detail: "ok"}
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

func printShapeMatrix(report shapeMatrix) {
	fmt.Println("=== Shape Compatibility Matrix ===")
	fmt.Println("Scope: structural compatibility of complex object shapes")
	fmt.Println()

	fmt.Printf("%-27s", "Shape case")
	for _, name := range report.FormatNames {
		fmt.Printf("  %-16s", name)
	}
	fmt.Println()

	for i, tc := range report.Tests {
		fmt.Printf("%-27s", tc.Name)
		for fi := range report.FormatNames {
			res := report.ResultsByFormat[fi][i]
			cell := "FAIL"
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				cell = "OK"
			}
			fmt.Printf("  %-16s", cell)
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Support summary:")
	for fi, name := range report.FormatNames {
		fmt.Printf("- %s: %d/%d shape cases roundtrip\n", name, report.PassCounts[fi], len(report.Tests))
	}

	fmt.Println()
	fmt.Println("Notable details:")
	for fi, name := range report.FormatNames {
		for ti, tc := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				continue
			}
			fmt.Printf("- %s / %s: %s\n", name, tc.Name, res.Detail)
		}
	}
	fmt.Println()
}

func buildBenchmarkCSV(results []benchmarkResult, baseline benchmarkResult, records, iters, warmup int, seed int64, indentSize int) (string, error) {
	if len(results) == 0 {
		return "", nil
	}

	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)
	err := w.Write([]string{
		"format", "baseline_format", "records", "iters", "warmup", "seed", "indent_size", "binary",
		"marshal_total_ms", "marshal_avg_ns", "unmarshal_total_ms", "unmarshal_avg_ns",
		"encoded_bytes", "encoded_chars", "roundtrip_ok", "roundtrip_details",
		"marshal_ratio_vs_baseline", "unmarshal_ratio_vs_baseline",
		"byte_ratio_vs_baseline", "char_ratio_vs_baseline",
	})
	if err != nil {
		return "", err
	}

	for _, res := range results {
		chars := ""
		if res.EncodedChars >= 0 {
			chars = fmt.Sprintf("%d", res.EncodedChars)
		}
		err = w.Write([]string{
			res.Name,
			baseline.Name,
			fmt.Sprintf("%d", records),
			fmt.Sprintf("%d", iters),
			fmt.Sprintf("%d", warmup),
			fmt.Sprintf("%d", seed),
			fmt.Sprintf("%d", indentSize),
			fmt.Sprintf("%t", res.Binary),
			fmt.Sprintf("%.6f", float64(res.MarshalTotal.Microseconds())/1000.0),
			fmt.Sprintf("%d", res.MarshalAvg.Nanoseconds()),
			fmt.Sprintf("%.6f", float64(res.UnmarshalTotal.Microseconds())/1000.0),
			fmt.Sprintf("%d", res.UnmarshalAvg.Nanoseconds()),
			fmt.Sprintf("%d", res.EncodedBytes),
			chars,
			fmt.Sprintf("%t", res.RoundTripOK),
			res.RoundTripDetails,
			ratioFloat(res.MarshalAvg, baseline.MarshalAvg),
			ratioFloat(res.UnmarshalAvg, baseline.UnmarshalAvg),
			ratioIntFloat(res.EncodedBytes, baseline.EncodedBytes),
			ratioIntFloat(res.EncodedChars, baseline.EncodedChars),
		})
		if err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func buildCharsetCSV(report charsetMatrix) (string, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	err := w.Write([]string{"format", "case", "marshal_ok", "roundtrip_ok", "status", "detail"})
	if err != nil {
		return "", err
	}

	for fi, name := range report.FormatNames {
		for ti, tc := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			status := "FAIL"
			if res.MarshalOK && res.RoundTripOK {
				status = "OK"
			}
			err = w.Write([]string{name, tc.Name, fmt.Sprintf("%t", res.MarshalOK), fmt.Sprintf("%t", res.RoundTripOK), status, res.Detail})
			if err != nil {
				return "", err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func buildShapeCSV(report shapeMatrix) (string, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	err := w.Write([]string{"format", "case", "encode_ok", "decode_ok", "roundtrip_ok", "status", "detail"})
	if err != nil {
		return "", err
	}

	for fi, name := range report.FormatNames {
		for ti, tc := range report.Tests {
			res := report.ResultsByFormat[fi][ti]
			status := "FAIL"
			if res.EncodeOK && res.DecodeOK && res.RoundTripOK {
				status = "OK"
			}
			err = w.Write([]string{
				name,
				tc.Name,
				fmt.Sprintf("%t", res.EncodeOK),
				fmt.Sprintf("%t", res.DecodeOK),
				fmt.Sprintf("%t", res.RoundTripOK),
				status,
				res.Detail,
			})
			if err != nil {
				return "", err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func emitCSVs(mode, benchmarkPath, benchmarkCSV, charsetCSV, shapeCSV string) {
	if mode == "stdout" {
		fmt.Println("=== CSV Benchmark Output ===")
		fmt.Print(benchmarkCSV)
		fmt.Println("=== CSV Charset Output ===")
		fmt.Print(charsetCSV)
		fmt.Println("=== CSV Shape Compatibility Output ===")
		fmt.Print(shapeCSV)
		return
	}

	if err := os.WriteFile(benchmarkPath, []byte(benchmarkCSV), 0644); err != nil {
		panic(fmt.Sprintf("failed to write benchmark csv: %v", err))
	}
	charsetPath := deriveCSVPath(benchmarkPath, ".charset.csv")
	if err := os.WriteFile(charsetPath, []byte(charsetCSV), 0644); err != nil {
		panic(fmt.Sprintf("failed to write charset csv: %v", err))
	}
	shapePath := deriveCSVPath(benchmarkPath, ".shape.csv")
	if err := os.WriteFile(shapePath, []byte(shapeCSV), 0644); err != nil {
		panic(fmt.Sprintf("failed to write shape csv: %v", err))
	}
	fmt.Printf("CSV written: %s\n", benchmarkPath)
	fmt.Printf("CSV written: %s\n", charsetPath)
	fmt.Printf("CSV written: %s\n", shapePath)
}

func deriveCSVPath(path, suffix string) string {
	if strings.HasSuffix(path, ".csv") {
		return strings.TrimSuffix(path, ".csv") + suffix
	}
	return path + suffix
}

func ratioFloat(v, base time.Duration) string {
	if base == 0 {
		return ""
	}
	return fmt.Sprintf("%.6f", float64(v)/float64(base))
}

func ratioIntFloat(v, base int) string {
	if base <= 0 || v < 0 {
		return ""
	}
	return fmt.Sprintf("%.6f", float64(v)/float64(base))
}

func deltaDurationText(v, base time.Duration) string {
	if base == 0 {
		return "n/a"
	}
	if v == base {
		return "same as baseline"
	}
	if v < base {
		pct := (float64(base-v) / float64(base)) * 100
		return fmt.Sprintf("faster %.1f%%", pct)
	}
	pct := (float64(v-base) / float64(base)) * 100
	return fmt.Sprintf("slower %.1f%%", pct)
}

func deltaIntText(v, base int) string {
	if base <= 0 || v < 0 {
		return "n/a"
	}
	if v == base {
		return "same as baseline"
	}
	if v < base {
		pct := (float64(base-v) / float64(base)) * 100
		return fmt.Sprintf("smaller %.1f%%", pct)
	}
	pct := (float64(v-base) / float64(base)) * 100
	return fmt.Sprintf("larger %.1f%%", pct)
}
