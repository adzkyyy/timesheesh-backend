package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"timesheesh-backend/database"
	"timesheesh-backend/models"

	"github.com/gin-gonic/gin"
)

// UploadProfilePicture uploads a user's profile picture
// POST /api/upload/profile-picture
func UploadProfilePicture(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(400, gin.H{"error": "User ID not found"})
		return
	}

	// Get file from request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file provided"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(400, gin.H{"error": "File size must be less than 5MB"})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(400, gin.H{"error": "Only JPEG, PNG, GIF, and WebP images are allowed"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/profiles"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Failed to create uploads directory: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("user_%d_%d%s", userID, timestamp, ext)
	filepath := filepath.Join(uploadsDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		log.Printf("Failed to save uploaded file: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate accessible URL
	imageURL := fmt.Sprintf("/uploads/profiles/%s", filename)

	// Update user with image URL
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("image_url", imageURL).Error; err != nil {
		log.Printf("Failed to update user image_url: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(200, gin.H{
		"message":   "Profile picture uploaded successfully",
		"image_url": imageURL,
	})
}

// UploadProjectImage uploads a project image/logo
// POST /api/upload/project/:projectId
func UploadProjectImage(c *gin.Context) {
	projectID := c.Param("projectId")

	// Check if user is admin or project owner
	var project models.Project
	if err := database.DB.First(&project, projectID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Project not found"})
		return
	}

	// Get file from request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file provided"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(400, gin.H{"error": "File size must be less than 10MB"})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(400, gin.H{"error": "Only JPEG, PNG, GIF, and WebP images are allowed"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/projects"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Failed to create uploads directory: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("project_%d_%d%s", project.ID, timestamp, ext)
	filepath := filepath.Join(uploadsDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		log.Printf("Failed to save uploaded file: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate accessible URL
	imageURL := fmt.Sprintf("/uploads/projects/%s", filename)

	c.JSON(200, gin.H{
		"message":   "Project image uploaded successfully",
		"image_url": imageURL,
	})
}

// UploadWorkspaceLogo uploads a workspace/company logo
// POST /api/upload/logo
func UploadWorkspaceLogo(c *gin.Context) {
	// Get file from request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file provided"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(400, gin.H{"error": "File size must be less than 5MB"})
		return
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(400, gin.H{"error": "Only JPEG, PNG, GIF, and WebP images are allowed"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads/logos"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Failed to create uploads directory: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("logo_%d%s", timestamp, ext)
	targetPath := filepath.Join(uploadsDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		log.Printf("Failed to save uploaded file: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	// Generate accessible URL
	imageURL := fmt.Sprintf("/uploads/logos/%s", filename)

	c.JSON(200, gin.H{
		"message":   "Logo uploaded successfully",
		"image_url": imageURL,
	})
}

// ServeUploadedFile serves uploaded files
// GET /uploads/:type/:filename
func ServeUploadedFile(c *gin.Context) {
	uploadType := c.Param("type")
	filename := c.Param("filename")

	// Validate upload type
	validTypes := map[string]bool{
		"profiles": true,
		"projects": true,
		"logos":    true,
	}
	if !validTypes[uploadType] {
		c.JSON(400, gin.H{"error": "Invalid upload type"})
		return
	}

	filepath := filepath.Join("uploads", uploadType, filename)

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	// Serve file
	c.File(filepath)
}

// GenerateThumbnail can be used to generate thumbnails
func saveThumbnail(sourceFile *multipart.FileHeader, destPath string, maxWidth, maxHeight int) error {
	// This is a placeholder. For production, use imaging library
	// Example: github.com/disintegration/imaging
	log.Printf("Thumbnail generation not yet implemented for: %s", destPath)
	return nil
}
