package routes

import (
	controllers "github.com/cvudumbarainformatika/backend/app/Http/Controllers"
	middleware "github.com/cvudumbarainformatika/backend/app/Http/Middleware"
	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// SetupRoutes configures all application routes and returns services that need scheduling
func SetupRoutes(router *gin.Engine, db *sqlx.DB, redis *redis.Client, cfg *config.Config) *services.BirthdayService {
	// Initialize Services
	mailService := services.NewMailService(cfg.Mail)
	waService := services.NewWhatsAppService(cfg.Zuwinda)
	birthdayService := services.NewBirthdayService(db, mailService, waService)

	// Initialize controllers
	authController := controllers.NewAuthController(db, cfg)
	avatarController := controllers.NewAvatarController()
	fileController := controllers.NewFileController()
	userController := controllers.NewUserController(db)
	beritaController := controllers.NewBeritaController(db, mailService)
	agendaController := controllers.NewAgendaController(db)
	uploadController := controllers.NewUploadController()
	homepageController := controllers.NewHomepageController(db)
	menuController := controllers.NewMenuController(db)
	contentController := controllers.NewContentController(db)
	contentController.InitTable()
	pdpiController := controllers.NewPDPIController(db, cfg)
	broadcastController := controllers.NewBroadcastController(mailService, waService, birthdayService, db, cfg.App)
	memberController := controllers.NewMemberController(db)
	dashboardController := controllers.NewDashboardController(db)
	documentController := controllers.NewDocumentController(db, cfg)

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

		// Homepage (Public)
		v1.GET("/homepage", homepageController.Get)

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
			berita.GET("/s/:slug", beritaController.GetBySlug)
		}

		// ==============================
		// Agenda Routes (Public GET)
		// ==============================
		agenda := v1.Group("/agenda")
		{
			agenda.GET("", agendaController.GetList)
			agenda.GET("/types", agendaController.GetTypes)
			agenda.GET("/s/:slug", agendaController.GetBySlug)
		}

		// ==============================
		// Menu Routes (Public GET)
		// ==============================
		menus := v1.Group("/menus")
		{
			menus.GET("", menuController.GetMenusByPosition)
			menus.GET("/:id", menuController.GetMenuByID)
		}

		// ==============================
		// Content Routes (Public GET)
		// ==============================
		v1.GET("/dynamic-content/*slug", contentController.GetContentBySlug)

		// ==============================
		// Members Directory (Public GET)
		// ==============================
		v1.GET("/members/search", pdpiController.SearchPublicMembers)

		// ==============================
		// Protected Routes (JWT Required)
		// ==============================
		protected := v1.Group("")
		protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
		{
			// Auth protected routes
			auth := protected.Group("/auth")
			{
				auth.GET("/me", authController.Me)
				auth.PUT("/profile", authController.UpdateProfile)
				auth.POST("/profile/change-password", authController.ChangePassword)
			}

			// Dashboard Stats
			protected.GET("/dashboard/stats", dashboardController.GetStats)

			// Homepage Management (Admin only)
			protected.POST("/homepage", homepageController.Update)

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
				beritaAdmin.GET("/:id", beritaController.GetByID)
				beritaAdmin.POST("", beritaController.Create)
				beritaAdmin.PUT("/:id", beritaController.Update)
				beritaAdmin.PATCH("/:id", beritaController.Patch)
				beritaAdmin.DELETE("/:id", beritaController.Delete)
			}

			// Agenda Management routes (Admin only)
			agendaAdmin := protected.Group("/agenda")
			{
				agendaAdmin.GET("/:id", agendaController.GetByID)
				agendaAdmin.POST("", agendaController.Create)
				agendaAdmin.PUT("/:id", agendaController.Update)
				agendaAdmin.PATCH("/:id", agendaController.Patch)
				agendaAdmin.DELETE("/:id", agendaController.Delete)
			}

			// Broadcast routes (Admin only)
			broadcastAdmin := protected.Group("/broadcast")
			{
				broadcastAdmin.POST("/berita/:id", broadcastController.BroadcastBerita)
				broadcastAdmin.POST("/berita-wa/:id", broadcastController.BroadcastBeritaWA)
				broadcastAdmin.POST("/agenda/:id", broadcastController.BroadcastAgenda)
				broadcastAdmin.POST("/agenda-wa/:id", broadcastController.BroadcastAgendaWA)
				broadcastAdmin.POST("/birthday-check", broadcastController.TriggerBirthdayGreetings)
			}

			// Menu Management routes (Admin only)
			menuAdmin := protected.Group("/menus")
			{
				menuAdmin.POST("", menuController.SaveMenus)
				menuAdmin.DELETE("/:id", menuController.DeleteMenu)
			}

			// Content Management routes (Admin only)
			contentAdmin := protected.Group("/dynamic-content")
			{
				contentAdmin.POST("", contentController.SaveContent)
			}

			// Members Management routes (Admin only)
			membersAdmin := protected.Group("/members")
			{
				membersAdmin.GET("", memberController.GetMembers)
				membersAdmin.GET("/filter-options", memberController.GetFilterOptions)
				membersAdmin.GET("/:id", memberController.GetMemberByID)
				membersAdmin.PUT("/:id", memberController.UpdateMember)
			}

			// Document Management routes (Protected)
			documents := protected.Group("/documents")
			{
				documents.GET("", documentController.GetList)
				documents.POST("", documentController.Upload)
				documents.DELETE("/:id", documentController.Delete)
			}

			// PDPI Integration routes (Protected)
			pdpi := protected.Group("/pdpi")
			{
				pdpi.POST("/sync-member", pdpiController.SyncMember)
				pdpi.POST("/sync-all-members", pdpiController.SyncAllMembers) // New: Sync all PDPI members
				pdpi.GET("/members", pdpiController.GetMembers)
				pdpi.GET("/member/:npa", pdpiController.GetMemberByNPA)
				pdpi.GET("/me", pdpiController.GetMyMemberData)
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

	return birthdayService
}
