package codec

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	toon "github.com/toon-format/toon-go"
	"github.com/vmihailenco/msgpack/v5"
	"go.yaml.in/yaml/v3"

	"toon-showcase/internal/core/entity"
	"toon-showcase/internal/core/usecase"
)

type xmlBenchmarkData struct {
	XMLName       xml.Name              `xml:"showcase"`
	Version       string                `xml:"version"`
	GeneratedAt   string                `xml:"generatedAt"`
	Owner         string                `xml:"owner"`
	DefaultLocale string                `xml:"defaultLocale"`
	Meta          []entity.TextEntry    `xml:"meta>entry"`
	Orders        []entity.Order        `xml:"orders>order"`
	OrderTable    []entity.OrderSummary `xml:"orderTable>summary"`
}

type formatCodec struct {
	key          string
	name         string
	binary       bool
	encodeData   func(entity.BenchmarkData) ([]byte, error)
	decodeData   func([]byte) (entity.BenchmarkData, error)
	marshalAny   func(any) ([]byte, error)
	unmarshalAny func([]byte, any) error
}

func (f formatCodec) Key() string { return f.key }

func (f formatCodec) Name() string { return f.name }

func (f formatCodec) Binary() bool { return f.binary }

func (f formatCodec) EncodeData(data entity.BenchmarkData) ([]byte, error) { return f.encodeData(data) }

func (f formatCodec) DecodeData(raw []byte) (entity.BenchmarkData, error) { return f.decodeData(raw) }

func (f formatCodec) MarshalAny(v any) ([]byte, error) { return f.marshalAny(v) }

func (f formatCodec) UnmarshalAny(raw []byte, out any) error { return f.unmarshalAny(raw, out) }

func Build(indentSize int) []usecase.Codec {
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

	marshalXMLPretty := func(data entity.BenchmarkData) ([]byte, error) {
		x := xmlBenchmarkData{
			Version:       data.Version,
			GeneratedAt:   data.GeneratedAt,
			Owner:         data.Owner,
			DefaultLocale: data.DefaultLocale,
			Meta:          data.Meta,
			OrderTable:    data.OrderTable,
			Orders:        orderedOrders(data.OrdersByID),
		}
		return xml.MarshalIndent(x, "", indent)
	}

	return []usecase.Codec{
		formatCodec{
			key:    "json-compact",
			name:   "JSON compact",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return json.Marshal(data)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				var out entity.BenchmarkData
				err := json.Unmarshal(raw, &out)
				return out, err
			},
			marshalAny: json.Marshal,
			unmarshalAny: func(raw []byte, out any) error {
				return json.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "json-pretty",
			name:   "JSON pretty",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return json.MarshalIndent(data, "", indent)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				var out entity.BenchmarkData
				err := json.Unmarshal(raw, &out)
				return out, err
			},
			marshalAny: func(v any) ([]byte, error) { return json.MarshalIndent(v, "", indent) },
			unmarshalAny: func(raw []byte, out any) error {
				return json.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "toon",
			name:   "TOON",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return toon.Marshal(data, toon.WithIndent(indentSize))
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				var out entity.BenchmarkData
				err := toon.Unmarshal(raw, &out)
				return out, err
			},
			marshalAny: func(v any) ([]byte, error) { return toon.Marshal(v, toon.WithIndent(indentSize)) },
			unmarshalAny: func(raw []byte, out any) error {
				return toon.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "yaml",
			name:   "YAML",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return marshalYAML(data)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				var out entity.BenchmarkData
				err := yaml.Unmarshal(raw, &out)
				return out, err
			},
			marshalAny: marshalYAML,
			unmarshalAny: func(raw []byte, out any) error {
				return yaml.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "toml",
			name:   "TOML",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return marshalTOML(data)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				var out entity.BenchmarkData
				err := toml.Unmarshal(raw, &out)
				return out, err
			},
			marshalAny: marshalTOML,
			unmarshalAny: func(raw []byte, out any) error {
				return toml.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "xml-compact",
			name:   "XML compact",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return marshalXMLBenchmark(data)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				return unmarshalXMLBenchmark(raw)
			},
			marshalAny: xml.Marshal,
			unmarshalAny: func(raw []byte, out any) error {
				return xml.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "xml-pretty",
			name:   "XML pretty",
			binary: false,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return marshalXMLPretty(data)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				return unmarshalXMLBenchmark(raw)
			},
			marshalAny: func(v any) ([]byte, error) {
				x, ok := v.(entity.BenchmarkData)
				if ok {
					return marshalXMLPretty(x)
				}
				return xml.MarshalIndent(v, "", indent)
			},
			unmarshalAny: func(raw []byte, out any) error {
				return xml.Unmarshal(raw, out)
			},
		},
		formatCodec{
			key:    "messagepack",
			name:   "MessagePack",
			binary: true,
			encodeData: func(data entity.BenchmarkData) ([]byte, error) {
				return msgpack.Marshal(data)
			},
			decodeData: func(raw []byte) (entity.BenchmarkData, error) {
				var out entity.BenchmarkData
				err := msgpack.Unmarshal(raw, &out)
				return out, err
			},
			marshalAny: msgpack.Marshal,
			unmarshalAny: func(raw []byte, out any) error {
				return msgpack.Unmarshal(raw, out)
			},
		},
	}
}

func marshalXMLBenchmark(data entity.BenchmarkData) ([]byte, error) {
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

func unmarshalXMLBenchmark(raw []byte) (entity.BenchmarkData, error) {
	var x xmlBenchmarkData
	if err := xml.Unmarshal(raw, &x); err != nil {
		return entity.BenchmarkData{}, err
	}
	ordersByID := make(map[string]entity.Order, len(x.Orders))
	for _, order := range x.Orders {
		ordersByID[order.OrderID] = order
	}
	return entity.BenchmarkData{
		Version:       x.Version,
		GeneratedAt:   x.GeneratedAt,
		Owner:         x.Owner,
		DefaultLocale: x.DefaultLocale,
		Meta:          x.Meta,
		OrdersByID:    ordersByID,
		OrderTable:    x.OrderTable,
	}, nil
}

func orderedOrders(m map[string]entity.Order) []entity.Order {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	orders := make([]entity.Order, 0, len(keys))
	for _, k := range keys {
		orders = append(orders, m[k])
	}
	return orders
}
