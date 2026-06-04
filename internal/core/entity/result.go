package entity

import "time"

type BenchmarkResult struct {
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

type CharsetCase struct {
	Name  string
	Value string
}

type CharsetResult struct {
	MarshalOK   bool
	RoundTripOK bool
	Detail      string
}

type CharsetMatrix struct {
	Tests           []CharsetCase
	FormatNames     []string
	ResultsByFormat [][]CharsetResult
	PassCounts      []int
}

type ShapeResult struct {
	EncodeOK    bool
	DecodeOK    bool
	RoundTripOK bool
	Detail      string
}

type ShapeCaseResult struct {
	CaseName string
	Result   ShapeResult
}

type ShapeMatrix struct {
	Tests           []string
	FormatNames     []string
	ResultsByFormat [][]ShapeResult
	PassCounts      []int
}
