package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type update struct {
	UpdateID int64  `json:"update_id"`
	Message  *tgMsg `json:"message"`
}

type tgMsg struct {
	Chat chatInfo `json:"chat"`
	From fromInfo `json:"from"`
	Text string   `json:"text"`
}

type chatInfo struct {
	ID int64 `json:"id"`
}

type fromInfo struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
}

func (b *BotClient) StartPolling(ctx context.Context) {
	log.Println("[Bot] starting long-polling…")
	offset := int64(0)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Bot] polling stopped")
			return
		default:
		}

		updates, err := b.getUpdates(offset, 30)
		if err != nil {
			log.Printf("[Bot] getUpdates error: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, u := range updates {
			offset = u.UpdateID + 1
			if u.Message == nil || u.Message.Text == "" {
				continue
			}
			cmd := strings.Fields(strings.TrimSpace(u.Message.Text))[0]
			switch cmd {
			case "/start":
				b.handleStart(u.Message)
			case "/help":
				b.handleHelp(u.Message)
			}
		}
	}
}

func (b *BotClient) getUpdates(offset int64, timeout int) ([]update, error) {
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=%d&allowed_updates=[\"message\"]",
		b.token, offset, timeout,
	)

	longPollClient := &http.Client{Timeout: time.Duration(timeout+5) * time.Second}
	resp, err := longPollClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		OK     bool     `json:"ok"`
		Result []update `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if !result.OK {
		return nil, fmt.Errorf("telegram getUpdates: not ok")
	}
	return result.Result, nil
}

func (b *BotClient) handleStart(msg *tgMsg) {
	text := FormatWelcomeMessage(msg.From.FirstName)
	if err := b.SendMessageWithShopButton(msg.Chat.ID, text); err != nil {
		log.Printf("[Bot] failed to send welcome: %v", err)
	}
}

func (b *BotClient) handleHelp(msg *tgMsg) {
	text := FormatHelpMessage(b.SupportContact)
	if err := b.SendMessageWithHelpButtons(msg.Chat.ID, text); err != nil {
		log.Printf("[Bot] failed to send help: %v", err)
	}
}
