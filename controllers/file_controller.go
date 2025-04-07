package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"turtle-stash/config"
	"turtle-stash/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadRequest struct {
	Filename       string `json:"filename" validate:"required"`
	CreatedOnInUTC int64  `json:"createdOnInUTC"`
	CreatedBy      string `json:"createdBy"`
	UpdatedOnInUTC int64  `json:"updatedOnInUTC"`
	UpdatedBy      string `json:"updatedBy"`
}

func UploadFile(c *fiber.Ctx) error {
	var req UploadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	fileID := uuid.New().String()
	file := models.File{
		ID:             fileID,
		Filename:       req.Filename,
		CreatedOnInUTC: time.Unix(req.CreatedOnInUTC, 0),
		CreatedBy:      req.CreatedBy,
		UpdatedOnInUTC: time.Unix(req.UpdatedOnInUTC, 0),
		UpdatedBy:      req.UpdatedBy,
	}

	if err := config.DB.Create(&file).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB insert failed"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"fileId":      fileID,
		"downloadUrl": fmt.Sprintf("http://localhost:8080/file/%s", fileID),
	})
}

func GetDownloadURL(c *fiber.Ctx) error {
	fileId := c.Params("fileId")
	var file models.File
	if err := config.DB.First(&file, "id = ?", fileId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}

	return c.JSON(fiber.Map{
		"fileId":      file.ID,
		"downloadUrl": fmt.Sprintf("http://localhost:8080/file/%s/download", file.ID),
	})
}

func DeleteFile(c *fiber.Ctx) error {
	fileId := c.Params("fileId")
	var file models.File
	if err := config.DB.First(&file, "id = ?", fileId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}

	os.Remove(file.FilePath) // optional

	if err := config.DB.Delete(&file).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Delete failed"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetFolderSnapshot(c *fiber.Ctx) error {
	folderId := c.Params("folderId")
	startIndex, _ := strconv.Atoi(c.Query("startIndex", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var files []models.File
	if err := config.DB.
		Where("folder_id = ?", folderId).
		Offset(startIndex).
		Limit(limit).
		Order("updated_on_in_utc desc").
		Find(&files).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DB query failed"})
	}

	var result []map[string]interface{}
	for _, f := range files {
		result = append(result, map[string]interface{}{
			"fileId":                f.ID,
			"filename":              f.Filename,
			"thumbnail_img":         f.ThumbnailImg,
			"lastModifiedDateInUTC": f.UpdatedOnInUTC,
			"creationDateInUTC":     f.CreatedOnInUTC,
		})
	}

	return c.JSON(fiber.Map{
		"folderId": folderId,
		"fileList": result,
	})
}
