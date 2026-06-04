package entity

import "encoding/xml"

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
