package handlers

import (
	"net/http"
	"strconv"

	"github.com/Hotelmurah2025/romansyah-11mei2025/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomHandler struct {
	DB *gorm.DB
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{DB: db}
}

type RoomRequest struct {
	HotelID      uint    `json:"hotel_id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	PriceWeekday float64 `json:"price_weekday" binding:"required"`
	PriceWeekend float64 `json:"price_weekend" binding:"required"`
}

func (h *RoomHandler) GetAllRooms(c *gin.Context) {
	var rooms []models.Room
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin {
		if err := h.DB.Joins("JOIN hotels ON rooms.hotel_id = hotels.id").
			Where("hotels.user_id = ?", userID.(uint)).
			Find(&rooms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
			return
		}
	} else {
		if err := h.DB.Find(&rooms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
	})
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	
	if err := h.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	
	var hotel models.Hotel
	if err := h.DB.First(&hotel, room.HotelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotel information"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var hotel models.Hotel
	if err := h.DB.First(&hotel, req.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	room := models.Room{
		HotelID:      req.HotelID,
		Name:         req.Name,
		Description:  req.Description,
		PriceWeekday: req.PriceWeekday,
		PriceWeekend: req.PriceWeekend,
	}
	
	if err := h.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Room created successfully",
		"room":    room,
	})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	
	if err := h.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	
	var hotel models.Hotel
	if err := h.DB.First(&hotel, room.HotelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotel information"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	var req RoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if room.HotelID != req.HotelID && role.(models.Role) != models.RoleAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change hotel ID"})
		return
	}
	
	room.HotelID = req.HotelID
	room.Name = req.Name
	room.Description = req.Description
	room.PriceWeekday = req.PriceWeekday
	room.PriceWeekend = req.PriceWeekend
	
	if err := h.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Room updated successfully",
		"room":    room,
	})
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	
	if err := h.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	
	var hotel models.Hotel
	if err := h.DB.First(&hotel, room.HotelID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotel information"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists && role.(models.Role) != models.RoleAdmin && hotel.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		return
	}
	
	var count int64
	h.DB.Model(&models.Booking{}).Where("room_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete room with bookings. Cancel bookings first."})
		return
	}
	
	if err := h.DB.Where("room_id = ?", id).Delete(&models.RatePlan{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated rate plans"})
		return
	}
	
	if err := h.DB.Where("room_id = ?", id).Delete(&models.RoomAvailability{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated availabilities"})
		return
	}
	
	if err := h.DB.Delete(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Room deleted successfully",
	})
}
