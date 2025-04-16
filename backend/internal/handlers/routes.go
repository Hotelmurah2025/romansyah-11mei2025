package handlers

import (
	"github.com/Hotelmurah2025/romansyah-11mei2025/internal/middleware"
	"github.com/Hotelmurah2025/romansyah-11mei2025/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	cfg := config.LoadConfig()

	api := router.Group("/api/v1")

	authHandler := NewAuthHandler(db)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))

	hotelHandler := NewHotelHandler(db)
	hotels := protected.Group("/hotels")
	{
		hotels.GET("", hotelHandler.GetAllHotels)
		hotels.GET("/:id", hotelHandler.GetHotel)
		hotels.POST("", hotelHandler.CreateHotel)
		hotels.PUT("/:id", hotelHandler.UpdateHotel)
		hotels.DELETE("/:id", hotelHandler.DeleteHotel)
		hotels.GET("/:id/rooms", hotelHandler.GetHotelRooms)
		hotels.GET("/search", hotelHandler.SearchHotels)
	}

	roomHandler := NewRoomHandler(db)
	rooms := protected.Group("/rooms")
	{
		rooms.GET("", roomHandler.GetAllRooms)
		rooms.GET("/:id", roomHandler.GetRoom)
		rooms.POST("", roomHandler.CreateRoom)
		rooms.PUT("/:id", roomHandler.UpdateRoom)
		rooms.DELETE("/:id", roomHandler.DeleteRoom)
	}
}
