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
	avatarController := controllers.NewAvatarController()
	fileController := controllers.NewFileController()
	userController := controllers.NewUserController(db)
	beritaController := controllers.NewBeritaController(db)
	uploadController := controllers.NewUploadController()

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

		// Avatar routes (public or protected based on config)
		avatars := v1.Group("/avatars")
		{
			// Public access to avatars by user ID or filename
			avatars.GET("/:user_id", avatarController.GetAvatar)
			avatars.GET("/file/:filename", avatarController.GetAvatarByName)
		}

		// Generic file routes (for all file types)
		files := v1.Group("/files")
		{
			// Serve files: GET /api/v1/files/:file_type/:filename
			files.GET("/:file_type/:filename", fileController.ServeFile)
			// List files (admin): GET /api/v1/files/:file_type/list
			files.GET("/:file_type/list", fileController.ListFiles)
		}

		// ==============================
		// Berita Routes (Public GET)
		// ==============================
		berita := v1.Group("/berita")
		{
			berita.GET("", beritaController.GetList)
			berita.GET("/categories", beritaController.GetCategories)
			berita.GET("/:slug", beritaController.GetBySlug)
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
				auth.PUT("/profile", authController.UpdateProfile)
				auth.POST("/profile/change-password", authController.ChangePassword)
			}

			// Upload endpoint (protected)
			protected.POST("/upload", uploadController.UploadFile)

			// User Management routes (Admin only)
			users := protected.Group("/users")
			{
				users.GET("/get-lists", userController.GetList)
				users.GET("/:id", userController.GetByID)
				users.POST("/create", userController.Create)
				users.PUT("update/:id", userController.Update)
				users.PATCH("patch/:id", userController.Patch)
				users.DELETE("delete/:id", userController.Delete)
			}

			// Berita Management routes (Admin only)
			beritaAdmin := protected.Group("/berita")
			{
				beritaAdmin.POST("", beritaController.Create)
				beritaAdmin.PUT("/:id", beritaController.Update)
				beritaAdmin.PATCH("/:id", beritaController.Patch)
				beritaAdmin.DELETE("/:id", beritaController.Delete)
			}
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
