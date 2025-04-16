package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Hotelmurah2025/romansyah-11mei2025/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingHandler struct {
	DB *gorm.DB
}

func NewBookingHandler(db *gorm.DB) *BookingHandler {
	return &BookingHandler{DB: db}
}

type BookingRequest struct {
	UserID       uint      `json:"user_id" binding:"required"`
	RoomID       uint      `json:"room_id" binding:"required"`
	RatePlanID   uint      `json:"rate_plan_id" binding:"required"`
	CheckinDate  string    `json:"checkin_date" binding:"required"`
	CheckoutDate string    `json:"checkout_date" binding:"required"`
	TotalPrice   float64   `json:"total_price" binding:"required"`
	Status       string    `json:"status"`
}

func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	var bookings []models.Booking
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists {
		if role.(models.Role) == models.RoleAdmin {
			if err := h.DB.Find(&bookings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
				return
			}
		} else if role.(models.Role) == models.RoleHotelOwner {
			if err := h.DB.Joins("JOIN rooms ON bookings.room_id = rooms.id").
				Joins("JOIN hotels ON rooms.hotel_id = hotels.id").
				Where("hotels.user_id = ?", userID.(uint)).
				Find(&bookings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
				return
			}
		} else {
			if err := h.DB.Where("user_id = ?", userID.(uint)).Find(&bookings).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
				return
			}
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"bookings": bookings,
	})
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	
	if err := h.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists {
		if role.(models.Role) == models.RoleAdmin {
		} else if role.(models.Role) == models.RoleHotelOwner {
			var count int64
			if err := h.DB.Table("rooms").
				Joins("JOIN hotels ON rooms.hotel_id = hotels.id").
				Where("rooms.id = ? AND hotels.user_id = ?", booking.RoomID, userID.(uint)).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
				return
			}
			
			if count == 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
				return
			}
		} else if booking.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"booking": booking,
	})
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	checkinDate, err := time.Parse("2006-01-02", req.CheckinDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkin date format. Use YYYY-MM-DD"})
		return
	}
	
	checkoutDate, err := time.Parse("2006-01-02", req.CheckoutDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkout date format. Use YYYY-MM-DD"})
		return
	}
	
	if checkinDate.After(checkoutDate) || checkinDate.Equal(checkoutDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Checkout date must be after checkin date"})
		return
	}
	
	var room models.Room
	if err := h.DB.First(&room, req.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	
	var ratePlan models.RatePlan
	if err := h.DB.First(&ratePlan, req.RatePlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rate plan not found"})
		return
	}
	
	if ratePlan.RoomID != req.RoomID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rate plan does not belong to the selected room"})
		return
	}
	
	var availabilityCount int64
	if err := h.DB.Model(&models.RoomAvailability{}).
		Where("room_id = ? AND rate_plan_id = ? AND date >= ? AND date < ? AND available_rooms > 0", 
			req.RoomID, req.RatePlanID, checkinDate, checkoutDate).
		Count(&availabilityCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check room availability"})
		return
	}
	
	nights := int(checkoutDate.Sub(checkinDate).Hours() / 24)
	if availabilityCount < int64(nights) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room not available for the selected dates"})
		return
	}
	
	status := req.Status
	if status == "" {
		status = string(models.BookingStatusPending)
	}
	
	booking := models.Booking{
		UserID:       req.UserID,
		RoomID:       req.RoomID,
		RatePlanID:   req.RatePlanID,
		CheckinDate:  checkinDate,
		CheckoutDate: checkoutDate,
		TotalPrice:   req.TotalPrice,
		Status:       models.BookingStatus(status),
	}
	
	if err := h.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}
	
	if err := h.DB.Model(&models.RoomAvailability{}).
		Where("room_id = ? AND rate_plan_id = ? AND date >= ? AND date < ?", 
			req.RoomID, req.RatePlanID, checkinDate, checkoutDate).
		UpdateColumn("available_rooms", gorm.Expr("available_rooms - 1")).Error; err != nil {
		h.DB.Delete(&booking)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room availability"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
		"booking": booking,
	})
}

func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	
	if err := h.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists {
		if role.(models.Role) == models.RoleAdmin {
		} else if role.(models.Role) == models.RoleHotelOwner {
			var count int64
			if err := h.DB.Table("rooms").
				Joins("JOIN hotels ON rooms.hotel_id = hotels.id").
				Where("rooms.id = ? AND hotels.user_id = ?", booking.RoomID, userID.(uint)).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
				return
			}
			
			if count == 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
				return
			}
		} else if booking.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	booking.Status = models.BookingStatus(req.Status)
	
	if err := h.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking updated successfully",
		"booking": booking,
	})
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	
	if err := h.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}
	
	userID, exists := c.Get("userID")
	role, _ := c.Get("role")
	
	if exists {
		if role.(models.Role) == models.RoleAdmin {
		} else if role.(models.Role) == models.RoleHotelOwner {
			var count int64
			if err := h.DB.Table("rooms").
				Joins("JOIN hotels ON rooms.hotel_id = hotels.id").
				Where("rooms.id = ? AND hotels.user_id = ?", booking.RoomID, userID.(uint)).
				Count(&count).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
				return
			}
			
			if count == 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
				return
			}
		} else if booking.UserID != userID.(uint) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}
	
	if err := h.DB.Model(&models.RoomAvailability{}).
		Where("room_id = ? AND rate_plan_id = ? AND date >= ? AND date < ?", 
			booking.RoomID, booking.RatePlanID, booking.CheckinDate, booking.CheckoutDate).
		UpdateColumn("available_rooms", gorm.Expr("available_rooms + 1")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room availability"})
		return
	}
	
	if err := h.DB.Delete(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking deleted successfully",
	})
}
