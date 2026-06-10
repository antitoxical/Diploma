package handler

import (
	"encoding/json"
	"log"
	"strings"

	"telegramshop_backend/pkg/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const auditLocalKey = "audit"

// AuditMeta allows handlers to override inferred audit data.
type AuditMeta struct {
	Action    string
	TableName string
	RecordID  *string
	Details   map[string]any
	Skip      bool
}

func (h *Handler) writeAudit(c *fiber.Ctx, adminID, action, tableName string, recordID *string, details map[string]any) {
	if adminID == "" || action == "" || tableName == "" {
		return
	}
	if err := h.adminService.WriteAuditLog(c.Context(), adminID, action, tableName, recordID, details); err != nil {
		log.Printf("audit log write failed: %v", err)
	}
}

func (h *Handler) adminIDFromContext(c *fiber.Ctx) string {
	claims, ok := c.Locals("admin").(*auth.Claims)
	if !ok || claims == nil {
		return ""
	}
	return claims.AdminID
}

// AdminAuditMiddleware records successful mutating admin requests.
func (h *Handler) AdminAuditMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}

		status := c.Response().StatusCode()
		if status < 200 || status >= 300 {
			return nil
		}

		method := c.Method()
		if method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions {
			return nil
		}

		routePath := c.Route().Path
		if strings.Contains(routePath, "/audit-logs") {
			return nil
		}

		adminID := h.adminIDFromContext(c)
		if adminID == "" {
			return nil
		}

		meta := inferAuditMeta(c, method, routePath)
		if custom, ok := c.Locals(auditLocalKey).(*AuditMeta); ok && custom != nil {
			if custom.Skip {
				return nil
			}
			meta = custom
		}
		if meta == nil || meta.Action == "" || meta.TableName == "" {
			return nil
		}

		recordID := meta.RecordID
		if recordID == nil {
			recordID = extractRecordIDFromResponse(c)
		}

		details := meta.Details
		if details == nil {
			details = paramsToDetails(c)
		}

		h.writeAudit(c, adminID, meta.Action, meta.TableName, recordID, details)
		return nil
	}
}

func inferAuditMeta(c *fiber.Ctx, method, routePath string) *AuditMeta {
	action := methodToAction(method)
	if action == "" {
		return nil
	}

	tableName := resourceFromRoute(routePath)
	if tableName == "" {
		return nil
	}

	meta := &AuditMeta{
		Action:    action,
		TableName: tableName,
		RecordID:  firstUUIDParam(c),
		Details:   paramsToDetails(c),
	}

	switch {
	case strings.Contains(routePath, "/moderate"):
		meta.Action = "update"
		meta.TableName = "reviews"
		meta.Details = mergeDetails(meta.Details, map[string]any{"operation": "moderate"})
	case strings.Contains(routePath, "/status"):
		meta.Details = mergeDetails(meta.Details, map[string]any{"operation": "status_change"})
	case strings.Contains(routePath, "/stock"):
		meta.Details = mergeDetails(meta.Details, map[string]any{"operation": "stock_change"})
	case strings.Contains(routePath, "/uploads/image"):
		meta.TableName = "uploads"
		meta.Action = "create"
		meta.RecordID = nil
	case strings.Contains(routePath, "/tags/"):
		meta.Details = mergeDetails(meta.Details, map[string]any{"operation": "product_tag"})
	case strings.Contains(routePath, "/images"):
		meta.Details = mergeDetails(meta.Details, map[string]any{"operation": "images"})
	}

	return meta
}

func methodToAction(method string) string {
	switch method {
	case fiber.MethodPost:
		return "create"
	case fiber.MethodPut, fiber.MethodPatch:
		return "update"
	case fiber.MethodDelete:
		return "delete"
	default:
		return ""
	}
}

func resourceFromRoute(routePath string) string {
	const marker = "/admin/"
	if idx := strings.Index(routePath, marker); idx >= 0 {
		rest := routePath[idx+len(marker):]
		if rest == "" {
			return ""
		}
		return strings.Split(rest, "/")[0]
	}

	trimmed := strings.TrimPrefix(routePath, "/")
	if trimmed == "" {
		return ""
	}
	return strings.Split(trimmed, "/")[0]
}

func firstUUIDParam(c *fiber.Ctx) *string {
	for _, key := range []string{"id", "image_id", "product_id", "tag_id"} {
		v := c.Params(key)
		if v == "" {
			continue
		}
		if _, err := uuid.Parse(v); err != nil {
			continue
		}
		id := v
		return &id
	}
	return nil
}

func paramsToDetails(c *fiber.Ctx) map[string]any {
	params := c.AllParams()
	if len(params) == 0 {
		return nil
	}
	out := make(map[string]any, len(params))
	for k, v := range params {
		out[k] = v
	}
	return out
}

func mergeDetails(base map[string]any, extra map[string]any) map[string]any {
	if len(extra) == 0 {
		return base
	}
	if base == nil {
		base = make(map[string]any, len(extra))
	}
	for k, v := range extra {
		base[k] = v
	}
	return base
}

func extractRecordIDFromResponse(c *fiber.Ctx) *string {
	body := c.Response().Body()
	if len(body) == 0 {
		return nil
	}

	var resp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil
	}
	if resp.Data.ID == "" {
		return nil
	}
	if _, err := uuid.Parse(resp.Data.ID); err != nil {
		return nil
	}
	id := resp.Data.ID
	return &id
}
