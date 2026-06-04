package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"toon-showcase/internal/core/entity"
)

func GenerateDataset(records int, seed int64) entity.BenchmarkData {
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

	ordersByID := make(map[string]entity.Order, records)
	orderTable := make([]entity.OrderSummary, 0, records)

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

		items := []entity.LineItem{
			{SKU: fmt.Sprintf("ENG-%03d", i), Title: englishTitles[i%len(englishTitles)], Quantity: item1Qty, UnitUSD: item1Price, Labels: []string{"english", "office", "stable"}},
			{SKU: fmt.Sprintf("JPN-%03d", i), Title: japaneseTitles[i%len(japaneseTitles)], Quantity: item2Qty, UnitUSD: item2Price, Labels: []string{"日本語", "craft", "gift"}},
			{SKU: fmt.Sprintf("CHN-%03d", i), Title: chineseTitles[i%len(chineseTitles)], Quantity: item3Qty, UnitUSD: item3Price, Labels: []string{"中文", "home", "limited"}},
		}

		total := float64(item1Qty)*item1Price + float64(item2Qty)*item2Price + float64(item3Qty)*item3Price

		created := base.Add(time.Duration(i) * time.Hour)
		createdText := created.Format(time.RFC3339Nano)
		shippedText := created.Add(3 * time.Hour).Format(time.RFC3339Nano)

		localizedNotes := []entity.LocalizedNote{
			{Lang: "en", Text: "Ship fast; include invoice copy."},
			{Lang: "ja", Text: "丁寧に梱包してください。"},
			{Lang: "zh", Text: "请小心包装并附上发票。"},
			{Lang: "emoji", Text: "🚀✨🙂"},
		}

		audit := []entity.AuditEvent{
			{At: createdText, Actor: "system", Action: "created", Message: "Order created via API"},
			{At: created.Add(20 * time.Minute).Format(time.RFC3339Nano), Actor: "qa.bot", Action: "validated", Message: "UTF-8 fields validated: 你好 / こんにちは / hello"},
			{At: shippedText, Actor: "ops", Action: "shipped", Message: "Tracking assigned 📦"},
		}

		var discount *entity.Discount
		if i%3 == 0 {
			discount = &entity.Discount{Code: fmt.Sprintf("SAVE-%02d", i%20), Percentage: 5 + float64(i%10)}
		}

		var special *string
		if i%2 == 0 {
			msg := fmt.Sprintf("Gift wrap: yes; note=%s / 注文 / 订单", emojiStatus[i%len(emojiStatus)])
			special = &msg
		}

		order := entity.Order{
			OrderID:   fmt.Sprintf("ORD-%06d", i+1),
			CreatedAt: createdText,
			Customer: entity.Customer{
				ID:          fmt.Sprintf("CUS-%05d", i+1000),
				DisplayName: fmt.Sprintf("User %d 你好 こんにちは", i+1),
				Email:       fmt.Sprintf("user%03d@example.com", i+1),
				Locale:      loc,
				Address: entity.Address{
					Line1:      fmt.Sprintf("%d Market Street", 100+i),
					City:       city,
					Region:     region,
					PostalCode: fmt.Sprintf("%05d", 10000+i),
					Country:    "Global",
				},
				Tags: []string{"vip", "newsletter", "emoji-🙂"},
			},
			Items: []entity.LineItem(items),
			Shipments: []entity.Shipment{
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
		orderTable = append(orderTable, entity.OrderSummary{
			OrderID:     order.OrderID,
			CustomerID:  order.Customer.ID,
			Locale:      order.Customer.Locale,
			Priority:    order.Priority,
			TotalUSD:    order.TotalUSD,
			EmojiStatus: order.EmojiStatus,
		})
	}

	meta := []entity.TextEntry{
		{Key: "project", Value: "TOON vs Multi-format showcase"},
		{Key: "languages", Value: "English / 日本語 / 中文"},
		{Key: "emoji", Value: "🚀✨🙂"},
		{Key: "explicitIndent", Value: "configured"},
	}

	return entity.BenchmarkData{
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
