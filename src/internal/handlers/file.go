package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/usecases/file"
	"github.com/gofiber/fiber/v2"
)

type FileHandler interface {
	Post(ctx *fiber.Ctx) error
	Get(c *fiber.Ctx) error
}

type fileHandler struct {
	fileUseCase file.UseCase
}

func NewFileHandler(fileUseCase file.UseCase) FileHandler {
	return &fileHandler{
		fileUseCase: fileUseCase,
	}
}

func (r *fileHandler) Post(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer src.Close()

	if file.Size > (100 * 1024) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file exceeds the maximum limit of 100KiB",
		})
	}

	fileName := file.Filename
	fileType := file.Header.Get("Content-Type")

	if !isAllowedFileType(fileName, fileType) {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "file type is not allowed",
		})
	}

	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileName)

	fullURL := ctx.Protocol() + "://" + ctx.Hostname() + "/v1/file?name="

	ceteex := context.WithValue(ctx.UserContext(), "request_url", fullURL)
	ctx.SetUserContext(ceteex)

	publicUrl, err := r.fileUseCase.UploadFile(ctx.UserContext(), file, src, filename)
	if err != nil {
		return ctx.JSON(fiber.Map{"message": err.Error()})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"uri": publicUrl,
	})
}

func isAllowedFileType(fileName, fileType string) bool {
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		// "application/octet-stream": true,
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedMimeTypes[fileType] {
		return false
	}

	if fileType == "application/octet-stream" {
		ext := strings.ToLower(filepath.Ext(fileName))
		if !allowedExtensions[ext] {
			return false
		}
	}

	return true
}

func (f *fileHandler) Get(c *fiber.Ctx) error {
	ctx := c.Context()

	fileName := c.Query("name")
	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file name is required",
		})
	}

	file, err := f.fileUseCase.GetFile(ctx, fileName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	// Detect content type
	buffer := make([]byte, 512)
	n, _ := file.Read(buffer)
	contentType := http.DetectContentType(buffer[:n])

	// Reset pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read file",
		})
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "inline; filename="+fileName)

	// Stream file
	return c.SendStream(file)
}
