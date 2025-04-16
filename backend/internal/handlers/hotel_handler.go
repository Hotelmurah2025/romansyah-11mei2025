package handlers

import (
	"net/http"
	"strconv"

	"github.com/Hotelmurah2025/romansyah-11mei2025/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HotelHandler struct {
	DB *gorm.DB
}

func NewHotelHandler(db *gorm.DB) *HotelHandler {
	return &HotelHandler{DB: db}
}

type HotelRequest struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Description string `json:"description"`
}

func (h *HotelHandler) GetAllHotels(c *gin.Context) {
	var hotels []models.Hotel
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin {
		if err := h.DB.Where("user_id = ?", userID.(uint)).Find(&hotels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels"})
			return
		}
	} else {
		if err := h.DB.Find(&hotels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels"})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"hotels": hotels,
	})
}

func (h *HotelHandler) GetHotel(c *gin.Context) {
	id := c.Param("id")
	var hotel models.Hotel
	
	if err := h.DB.First(&hotel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"hotel": hotel,
	})
}

func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var req HotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	
	hotel := models.Hotel{
		UserID:      userID.(uint),
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
	}
	
	if err := h.DB.Create(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Hotel created successfully",
		"hotel":   hotel,
	})
}

func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	id := c.Param("id")
	var hotel models.Hotel
	
	if err := h.DB.First(&hotel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	var req HotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	hotel.Name = req.Name
	hotel.Location = req.Location
	hotel.Description = req.Description
	
	if err := h.DB.Save(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel updated successfully",
		"hotel":   hotel,
	})
}

func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	id := c.Param("id")
	var hotel models.Hotel
	
	if err := h.DB.First(&hotel, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	var count int64
	h.DB.Model(&models.Room{}).Where("hotel_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete hotel with rooms. Delete rooms first."})
		return
	}
	
	if err := h.DB.Delete(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hotel"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Hotel deleted successfully",
	})
}

func (h *HotelHandler) GetHotelRooms(c *gin.Context) {
	hotelID := c.Param("id")
	var hotel models.Hotel
	
	if err := h.DB.First(&hotel, hotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	var rooms []models.Room
	if err := h.DB.Where("hotel_id = ?", hotelID).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

func (h *HotelHandler) SearchHotels(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}
	
	var hotels []models.Hotel
	if err := h.DB.Where("name LIKE ? OR location LIKE ?", "%"+query+"%", "%"+query+"%").Find(&hotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search hotels"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"hotels": hotels,
	})
}
