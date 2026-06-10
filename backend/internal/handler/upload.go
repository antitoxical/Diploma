package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const maxImageUploadBytes = 5 * 1024 * 1024

var allowedUploadExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

// имя файла = UUID + допустимое расширение
var uploadedImageFileRx = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\.(jpg|jpeg|png|webp|gif)$`)

func uploadRootDir() string {
	d := strings.TrimSpace(os.Getenv("UPLOAD_DIR"))
	if d == "" {
		return "uploads"
	}
	return d
}

func uploadImagesDir() string {
	return filepath.Join(uploadRootDir(), "images")
}

func safeJoinUploadImages(filename string) (full string, ok bool) {
	if !uploadedImageFileRx.MatchString(filename) {
		return "", false
	}
	base := filepath.Clean(uploadImagesDir())
	full = filepath.Join(base, filename)
	rel, err := filepath.Rel(base, full)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return full, true
}

// ServeUploadedImage отдаёт сохранённое изображение (публичный GET).
func (h *Handler) ServeUploadedImage(c *fiber.Ctx) error {
	name := c.Params("filename")
	full, ok := safeJoinUploadImages(name)
	if !ok {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if _, err := os.Stat(full); os.IsNotExist(err) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.SendFile(full)
}

// UploadAdminImage сохраняет файл из multipart field "file" (JWT admin).
func (h *Handler) UploadAdminImage(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_no_file", "Нужно поле multipart «file»"))
	}
	if file.Size > maxImageUploadBytes {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_file_too_large", "Не более 5 МБ"))
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" || !allowedUploadExt[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_bad_extension", "Допустимы jpg, jpeg, png, webp, gif"))
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_open_upload", err.Error()))
	}
	defer src.Close()

	dir := uploadImagesDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_mkdir", err.Error()))
	}

	id := uuid.New().String()
	safeName := id + ext
	dstPath := filepath.Join(dir, safeName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_file", err.Error()))
	}
	defer dst.Close()

	lr := &io.LimitedReader{R: src, N: maxImageUploadBytes + 1}
	written, err := io.Copy(dst, lr)
	if err != nil {
		_ = os.Remove(dstPath)
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_save_file", err.Error()))
	}
	if written > maxImageUploadBytes {
		_ = os.Remove(dstPath)
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_file_too_large", "Не более 5 МБ"))
	}

	publicPath := fmt.Sprintf("%s/uploads/images/%s", basePath, safeName)
	return c.JSON(web.OkResp("success_upload", fiber.Map{"url": publicPath}))
}
