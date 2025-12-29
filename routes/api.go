package routes

import (
	controllers "github.com/cvudumbarainformatika/backend/app/Http/Controllers"
	middleware "github.com/cvudumbarainformatika/backend/app/Http/Middleware"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *sqlx.DB, redis *redis.Client, cfg *config.Config) {
	// Initialize controllers
	authController := controllers.NewAuthController(db, cfg)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// ==============================
		// Public Routes (No Auth Required)
		// ==============================
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.Refresh)
			auth.POST("/logout", authController.Logout)
		}

		// ==============================
		// Protected Routes (JWT Required)
		// ==============================
		protected := v1.Group("/")
		protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
		{
			// Auth protected routes
			auth := protected.Group("/auth")
			{
				auth.GET("/me", authController.Me)
			}

			// TODO: Add your protected routes here
			// Example: User CRUD routes
			// users := protected.Group("/users")
			// {
			//     users.GET("/get-list", userController.GetAllUsers)
			//     users.GET("/:id", userController.GetUserByID)
			//     users.POST("/create", userController.CreateUser)
			//     users.PUT("/update/:id", userController.UpdateUser)
			//     users.DELETE("/delete/:id", userController.DeleteUser)
			// }
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})
}
