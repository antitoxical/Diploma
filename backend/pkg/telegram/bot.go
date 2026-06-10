package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type BotClient struct {
	token          string
	client         *http.Client
	WebAppURL      string
	AdminPanelURL  string
	AdminChatID    int64
	SupportContact string
}

func NewBotClient(token, webAppURL, adminPanelURL, supportContact string, adminChatID int64) *BotClient {
	return &BotClient{
		token:          token,
		client:         &http.Client{Timeout: 10 * time.Second},
		WebAppURL:      webAppURL,
		AdminPanelURL:  adminPanelURL,
		AdminChatID:    adminChatID,
		SupportContact: supportContact,
	}
}

func (b *BotClient) HasAdmin() bool {
	return b != nil && b.AdminChatID != 0
}

func (b *BotClient) webAppPath(path string) string {
	if b.WebAppURL == "" {
		return ""
	}
	if path == "" {
		return b.WebAppURL
	}
	if strings.HasPrefix(path, "/") {
		return b.WebAppURL + path
	}
	return b.WebAppURL + "/" + path
}

type inlineKeyboardButton struct {
	Text   string      `json:"text"`
	URL    string      `json:"url,omitempty"`
	WebApp *webAppInfo `json:"web_app,omitempty"`
}

type webAppInfo struct {
	URL string `json:"url"`
}

type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type sendMessageRequest struct {
	ChatID      int64                 `json:"chat_id"`
	Text        string                `json:"text"`
	ParseMode   string                `json:"parse_mode,omitempty"`
	ReplyMarkup *inlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (b *BotClient) SendMessage(chatID int64, text string) error {
	return b.sendMsg(chatID, text, nil)
}

func (b *BotClient) SendMessageWithShopButton(chatID int64, text string) error {
	var markup *inlineKeyboardMarkup
	if b.WebAppURL != "" {
		markup = &inlineKeyboardMarkup{
			InlineKeyboard: [][]inlineKeyboardButton{
				{{Text: "🛒 Открыть магазин", WebApp: &webAppInfo{URL: b.WebAppURL}}},
			},
		}
	}
	return b.sendMsg(chatID, text, markup)
}

func (b *BotClient) SendMessageWithOrderButton(chatID int64, text string) error {
	var markup *inlineKeyboardMarkup
	if b.WebAppURL != "" {
		markup = &inlineKeyboardMarkup{
			InlineKeyboard: [][]inlineKeyboardButton{
				{{Text: "📦 Мои заказы", WebApp: &webAppInfo{URL: b.webAppPath("/orders")}}},
				{{Text: "🛒 Продолжить покупки", WebApp: &webAppInfo{URL: b.webAppPath("")}}},
			},
		}
	}
	return b.sendMsg(chatID, text, markup)
}

func (b *BotClient) SendMessageWithHelpButtons(chatID int64, text string) error {
	var rows [][]inlineKeyboardButton
	if url := b.webAppPath(""); url != "" {
		rows = append(rows, []inlineKeyboardButton{
			{Text: "🛒 Открыть магазин", WebApp: &webAppInfo{URL: url}},
		})
	}
	if url := b.webAppPath("/orders"); url != "" {
		rows = append(rows, []inlineKeyboardButton{
			{Text: "📦 Мои заказы", WebApp: &webAppInfo{URL: url}},
		})
	}
	if supportURL := supportTelegramURL(b.SupportContact); supportURL != "" {
		rows = append(rows, []inlineKeyboardButton{
			{Text: "💬 Поддержка", URL: supportURL},
		})
	}
	var markup *inlineKeyboardMarkup
	if len(rows) > 0 {
		markup = &inlineKeyboardMarkup{InlineKeyboard: rows}
	}
	return b.sendMsg(chatID, text, markup)
}

func (b *BotClient) SendMessageWithReviewButton(chatID int64, orderID, productID, text string) error {
	reviewPath := "/orders"
	if productID != "" {
		reviewPath = "/product/" + productID + "/reviews"
	} else if orderID != "" {
		reviewPath = "/order/" + orderID
	}
	var markup *inlineKeyboardMarkup
	if url := b.webAppPath(reviewPath); url != "" {
		markup = &inlineKeyboardMarkup{
			InlineKeyboard: [][]inlineKeyboardButton{
				{{Text: "⭐ Оценить покупку", WebApp: &webAppInfo{URL: url}}},
				{{Text: "📦 Мои заказы", WebApp: &webAppInfo{URL: b.webAppPath("/orders")}}},
			},
		}
	}
	return b.sendMsg(chatID, text, markup)
}

func (b *BotClient) SendMessageWithProductButton(chatID int64, productID, text string) error {
	var markup *inlineKeyboardMarkup
	if url := b.webAppPath("/product/" + productID); url != "" {
		markup = &inlineKeyboardMarkup{
			InlineKeyboard: [][]inlineKeyboardButton{
				{{Text: "🛍 Открыть товар", WebApp: &webAppInfo{URL: url}}},
			},
		}
	}
	return b.sendMsg(chatID, text, markup)
}

func (b *BotClient) SendAdminNewOrderNotification(orderID, text string) error {
	if !b.HasAdmin() {
		return nil
	}
	var markup *inlineKeyboardMarkup
	if b.AdminPanelURL != "" {
		adminURL := strings.TrimRight(b.AdminPanelURL, "/") + "/orders/show/" + orderID
		markup = &inlineKeyboardMarkup{
			InlineKeyboard: [][]inlineKeyboardButton{
				{{Text: "⚙️ Открыть в админ-панели", URL: adminURL}},
			},
		}
	}
	return b.sendMsg(b.AdminChatID, text, markup)
}

func (b *BotClient) BroadcastMessage(chatIDs []int64, text string, withShopButton bool) {
	for _, chatID := range chatIDs {
		if chatID == 0 {
			continue
		}
		var err error
		if withShopButton {
			err = b.SendMessageWithShopButton(chatID, text)
		} else {
			err = b.SendMessage(chatID, text)
		}
		if err != nil {
			continue
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func supportTelegramURL(contact string) string {
	contact = strings.TrimSpace(contact)
	if contact == "" {
		return ""
	}
	contact = strings.TrimPrefix(contact, "https://t.me/")
	contact = strings.TrimPrefix(contact, "t.me/")
	contact = strings.TrimPrefix(contact, "@")
	if contact == "" {
		return ""
	}
	return "https://t.me/" + contact
}

func (b *BotClient) sendMsg(chatID int64, text string, markup *inlineKeyboardMarkup) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	msg := sendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: markup,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal sendMessage: %w", err)
	}

	resp, err := b.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram sendMessage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram sendMessage: status %d", resp.StatusCode)
	}
	return nil
}

var statusLabels = map[string]string{
	"created":    "✅ Создан",
	"paid":       "💳 Оплачен",
	"delivering": "🚚 Доставляется",
	"completed":  "🎉 Завершён",
	"cancelled":  "❌ Отменён",
}

func FormatOrderStatusMessage(orderID, status string) string {
	label, ok := statusLabels[status]
	if !ok {
		label = status
	}
	return fmt.Sprintf(
		"📋 <b>Статус заказа обновлён</b>\n\n"+
			"Заказ: <code>#%s</code>\n"+
			"Статус: %s",
		shortOrderID(orderID), label,
	)
}

func FormatOrderCancelledMessage(orderID, comment string) string {
	msg := fmt.Sprintf(
		"❌ <b>Заказ отменён</b>\n\n"+
			"Заказ <code>#%s</code> был отменён.\n\n"+
			"Если оплата уже прошла, средства вернутся в течение 3–5 рабочих дней. "+
			"По вопросам возврата обратитесь в поддержку.",
		shortOrderID(orderID),
	)
	if comment != "" {
		msg += "\n\n💬 Комментарий: " + comment
	}
	return msg
}

func FormatOrderCompletedReviewMessage(orderID string) string {
	return fmt.Sprintf(
		"🎉 <b>Заказ доставлен!</b>\n\n"+
			"Заказ <code>#%s</code> выполнен. Поделитесь впечатлением — "+
			"ваш отзыв поможет другим покупателям.",
		shortOrderID(orderID),
	)
}

func FormatPromocodeMessage(code string, discountPercent int) string {
	return fmt.Sprintf(
		"🎁 <b>Для вас скидка %d%%</b>\n\n"+
			"Используйте промокод <code>%s</code> при оформлении заказа в магазине.",
		discountPercent, code,
	)
}

func FormatBackInStockMessage(productName string) string {
	return fmt.Sprintf(
		"🔔 <b>Товар снова в наличии</b>\n\n"+
			"«%s» снова доступен для заказа. Успейте оформить покупку!",
		productName,
	)
}

func FormatAdminNewOrderMessage(orderID string, totalKopeks int64) string {
	rubles := float64(totalKopeks) / 100
	return fmt.Sprintf(
		"🆕 <b>Новый заказ</b>\n\n"+
			"Заказ: <code>#%s</code>\n"+
			"Сумма: <b>%.2f Br</b>",
		shortOrderID(orderID), rubles,
	)
}

func FormatHelpMessage(supportContact string) string {
	support := strings.TrimSpace(supportContact)
	if support == "" {
		support = "через кнопку «Поддержка» ниже"
	} else if !strings.HasPrefix(support, "http") {
		if !strings.HasPrefix(support, "@") {
			support = "@" + support
		}
		support = "в Telegram: " + support
	}
	return "ℹ️ <b>Справка AlphaSales</b>\n\n" +
		"🛒 <b>Магазин</b> — нажмите «Открыть магазин», чтобы перейти в каталог Mini App.\n\n" +
		"📦 <b>Заказы</b> — раздел «Мои заказы»: статусы, история и детали доставки.\n\n" +
		"💬 <b>Поддержка</b> — " + support + ".\n\n" +
		"Команды бота:\n" +
		"/start — приветствие и кнопка магазина\n" +
		"/help — эта справка"
}

func shortOrderID(orderID string) string {
	if len(orderID) > 8 {
		return orderID[:8]
	}
	return orderID
}

func OrderStatusLabel(status string) string {
	if label, ok := statusLabels[status]; ok {
		return label
	}
	return status
}

func FormatOrderCreatedInApp(orderID string, totalKopeks int64, itemCount int) (title, message string) {
	rubles := float64(totalKopeks) / 100
	title = "Заказ принят"
	message = fmt.Sprintf(
		"Заказ #%s оформлен. Товаров: %d. Сумма: %.2f Br.",
		shortOrderID(orderID), itemCount, rubles,
	)
	return title, message
}

func FormatOrderStatusInApp(orderID, status, comment string) (title, message string) {
	switch status {
	case "cancelled":
		return FormatOrderCancelledInApp(orderID, comment)
	case "completed":
		return FormatOrderCompletedInApp(orderID)
	default:
		title = "Статус заказа изменён"
		message = fmt.Sprintf("Заказ #%s: %s", shortOrderID(orderID), OrderStatusLabel(status))
		if comment != "" {
			message += "\nКомментарий: " + comment
		}
		return title, message
	}
}

func FormatOrderCancelledInApp(orderID, comment string) (title, message string) {
	title = "Заказ отменён"
	message = fmt.Sprintf(
		"Заказ #%s отменён. При оплате средства вернутся в течение 3–5 рабочих дней.",
		shortOrderID(orderID),
	)
	if comment != "" {
		message += " Комментарий: " + comment
	}
	return title, message
}

func FormatOrderCompletedInApp(orderID string) (title, message string) {
	title = "Заказ выполнен"
	message = fmt.Sprintf(
		"Заказ #%s доставлен. Оцените покупку в разделе отзывов.",
		shortOrderID(orderID),
	)
	return title, message
}

func FormatOrderCreatedMessage(orderID string, totalKopeks int64, itemCount int) string {
	return FormatOrderCreatedMessageDetailed(OrderCreatedNotify{
		OrderID:     orderID,
		TotalKopeks: totalKopeks,
	})
}

type OrderLineNotify struct {
	Name       string
	Quantity   int
	UnitKopeks int64
}

type OrderCreatedNotify struct {
	OrderID      string
	TotalKopeks  int64
	Items        []OrderLineNotify
	DeliveryName string
	PaymentName  string
	Address      string
}

func FormatOrderCreatedMessageDetailed(n OrderCreatedNotify) string {
	short := shortOrderID(n.OrderID)
	rubles := float64(n.TotalKopeks) / 100
	totalQty := 0
	for _, it := range n.Items {
		totalQty += it.Quantity
	}

	var b strings.Builder
	b.WriteString("🛍 <b>Заказ принят!</b>\n\n")
	b.WriteString(fmt.Sprintf("Заказ: <code>#%s</code>\n", short))

	if len(n.Items) > 0 {
		b.WriteString("\n<b>Состав:</b>\n")
		for _, it := range n.Items {
			lineRub := float64(it.UnitKopeks*int64(it.Quantity)) / 100
			name := it.Name
			if name == "" {
				name = "Товар"
			}
			b.WriteString(fmt.Sprintf("• %s × %d — %.2f Br\n", escapeHTML(name), it.Quantity, lineRub))
		}
	}

	if len(n.Items) > 0 {
		b.WriteString(fmt.Sprintf("\nПозиций: %d • Единиц: %d\n", len(n.Items), totalQty))
	}
	b.WriteString(fmt.Sprintf("Сумма: <b>%.2f Br</b>\n", rubles))

	if n.DeliveryName != "" {
		b.WriteString(fmt.Sprintf("Доставка: %s\n", escapeHTML(n.DeliveryName)))
	}
	if n.PaymentName != "" {
		b.WriteString(fmt.Sprintf("Оплата: %s\n", escapeHTML(n.PaymentName)))
	}
	if n.Address != "" {
		addr := strings.TrimSpace(n.Address)
		if len(addr) > 220 {
			addr = addr[:217] + "..."
		}
		b.WriteString("\n📍 ")
		b.WriteString(escapeHTML(strings.ReplaceAll(addr, "\n", "\n📍 ")))
		b.WriteString("\n")
	}

	b.WriteString("\nМы уведомим вас об изменении статуса.")
	return b.String()
}

func FormatOrderCreatedInAppDetailed(n OrderCreatedNotify) (title, message string) {
	title = "Заказ принят"
	short := shortOrderID(n.OrderID)
	rubles := float64(n.TotalKopeks) / 100
	totalQty := 0
	for _, it := range n.Items {
		totalQty += it.Quantity
	}
	message = fmt.Sprintf(
		"Заказ #%s оформлен. Позиций: %d, единиц: %d. Сумма: %.2f Br.",
		short, len(n.Items), totalQty, rubles,
	)
	if n.DeliveryName != "" {
		message += " Доставка: " + n.DeliveryName + "."
	}
	return title, message
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func FormatWelcomeMessage(firstName string) string {
	name := firstName
	if name == "" {
		name = "друг"
	}
	return fmt.Sprintf(
		"👋 <b>Привет, %s!</b>\n\n"+
			"Добро пожаловать в <b>AlphaSales</b> — магазин электроники.\n\n"+
			"🔹 Широкий выбор товаров\n"+
			"🔹 Быстрая доставка\n"+
			"🔹 Уведомления о статусе заказа\n\n"+
			"Нажмите кнопку ниже, чтобы открыть каталог 👇",
		name,
	)
}
