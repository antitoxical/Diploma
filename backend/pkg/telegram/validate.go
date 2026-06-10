package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

var (
	ErrInvalidHash   = errors.New("invalid initData hash")
	ErrMissingHash   = errors.New("hash parameter missing in initData")
	ErrEmptyInitData = errors.New("empty initData")
)

// ValidateInitData verifies the Telegram Mini App initData string
// using HMAC-SHA256 as described in the official Telegram documentation.
//
// Algorithm:
//  1. Parse initData as URL-encoded query string.
//  2. Extract the "hash" parameter and remove it from the set.
//  3. Sort remaining key=value pairs alphabetically by key.
//  4. Join them with "\n" to form the data-check-string.
//  5. Compute HMAC-SHA256 of "WebAppData" using botToken as key → secretKey.
//  6. Compute HMAC-SHA256 of data-check-string using secretKey.
//  7. Compare with the extracted hash.
func ValidateInitData(initData, botToken string) error {
	if initData == "" {
		return ErrEmptyInitData
	}

	params, err := url.ParseQuery(initData)
	if err != nil {
		return fmt.Errorf("parse initData: %w", err)
	}

	hash := params.Get("hash")
	if hash == "" {
		return ErrMissingHash
	}
	params.Del("hash")

	var pairs []string
	for k := range params {
		pairs = append(pairs, k+"="+params.Get(k))
	}
	sort.Strings(pairs)
	dataCheckString := strings.Join(pairs, "\n")

	secretKey := hmacSHA256([]byte("WebAppData"), []byte(botToken))
	computedHash := hmacSHA256(secretKey, []byte(dataCheckString))

	if hex.EncodeToString(computedHash) != hash {
		return ErrInvalidHash
	}

	return nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
