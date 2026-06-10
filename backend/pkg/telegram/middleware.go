package telegram

import (
	"encoding/json"
	"net/url"

	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
)

type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

// InitDataMiddleware returns Fiber middleware that validates
// the Telegram initData sent in the X-Telegram-Init-Data header.
// On success it stores the parsed TelegramUser in c.Locals("tg_user").
func InitDataMiddleware(botToken string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		initData := c.Get("X-Telegram-Init-Data")
		if initData == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				web.ErrorResp("error_missing_init_data", "X-Telegram-Init-Data header is required"),
			)
		}

		if botToken != "" {
			if err := ValidateInitData(initData, botToken); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(
					web.ErrorResp("error_invalid_init_data", err.Error()),
				)
			}
		}

		params, _ := url.ParseQuery(initData)
		userJSON := params.Get("user")
		if userJSON != "" {
			var u TelegramUser
			if err := json.Unmarshal([]byte(userJSON), &u); err == nil {
				c.Locals("tg_user", &u)
			}
		}

		return c.Next()
	}
}
