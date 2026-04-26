package bootstrap

import (
	"fmt"
	"log"
	"time"

	services "github.com/cvudumbarainformatika/backend/app/Services"

	exceptions "github.com/cvudumbarainformatika/backend/app/Exceptions"
	middleware "github.com/cvudumbarainformatika/backend/app/Http/Middleware"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/cvudumbarainformatika/backend/database"
	"github.com/cvudumbarainformatika/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Application represents the main application
type Application struct {
	Router *gin.Engine
	DB     *database.Database
	Redis           *redis.Client
	Config          *config.Config
	BirthdayService *services.BirthdayService
	RSSService      *services.RSSService
}

// NewApplication creates and initializes a new application instance
func NewApplication() (*Application, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Set Gin mode based on environment
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Gin engine
	router := gin.New()

	// Register global middleware (in order)
	router.Use(gin.Recovery()) // Panic recovery
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS))

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)
	router.Use(rateLimiter.Middleware())

	// Error handler
	router.Use(exceptions.ErrorHandler())

	// Note: File serving is now handled via API endpoints:
	// GET /api/v1/files/:file_type/:filename
	// This provides security validation and access control
	// instead of direct static file access

	// Setup database connection
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	// Run migrations (Automatic execution enabled)
	migrator := database.NewMigrator(db.DB, "database/migrations")
	if err := migrator.RunMigrations(); err != nil {
		log.Printf("Warning: failed to run migrations: %v", err)
		// We don't return error here to allow app to start even if migrations fail
		// (e.g. table already exists or connection issue)
	}

	// Run seeders (Manual execution recommended)
	// if err := seeders.RunSeeders(db.DB); err != nil {
	// 	return nil, fmt.Errorf("failed to run seeders: %w", err)
	// }

	// Setup Redis connection
	rdb, err := database.InitRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	// Setup routes and get services that need scheduling
	birthdayService := routes.SetupRoutes(router, db.DB, rdb, cfg)

	app := &Application{
		Router:          router,
		DB:              db,
		Redis:           rdb,
		Config:          cfg,
		BirthdayService: birthdayService,
		RSSService:      services.NewRSSService(db.DB, rdb),
	}

	// Start background scheduler
	app.StartScheduler()

	return app, nil
}

// Run starts the application server
func (app *Application) Run() error {
	addr := fmt.Sprintf(":%s", app.Config.App.Port)
	log.Printf("Starting %s on %s (env: %s)", app.Config.App.Name, addr, app.Config.App.Env)
	return app.Router.Run(addr)
}

// StartScheduler initiates background tasks
func (app *Application) StartScheduler() {
	go func() {
		log.Println("Background scheduler started")
		
		// Run every 1 hour to check if it's birthday time
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		// Initial check on startup
		if err := app.BirthdayService.CheckAndSendGreetings(); err != nil {
			log.Printf("Error during startup birthday check: %v", err)
		}
		
		// Initial RSS sync on startup
		log.Println("Executing initial RSS sync on startup...")
		go func() {
			if err := app.RSSService.FetchAndStoreFeeds(); err != nil {
				log.Printf("Error during startup RSS sync: %v", err)
			}
		}()

		for range ticker.C {
			now := time.Now()
			
			// 1. Birthday Check (Only at 08:00 AM)
			if now.Hour() == 8 {
				log.Println("Executing scheduled birthday greetings...")
				if err := app.BirthdayService.CheckAndSendGreetings(); err != nil {
					log.Printf("Error executing scheduled birthday greetings: %v", err)
				}
			}

			// 2. RSS Sync (Every 6 hours)
			if now.Hour() % 6 == 0 {
				log.Println("Executing scheduled RSS sync...")
				if err := app.RSSService.FetchAndStoreFeeds(); err != nil {
					log.Printf("Error executing scheduled RSS sync: %v", err)
				}
			}
		}
	}()
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() error {
	log.Println("Shutting down application...")

	// Close database connections
	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Println("Database connection closed")
	}

	// Close Redis connection
	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			return fmt.Errorf("failed to close redis: %w", err)
		}
		log.Println("Redis connection closed")
	}

	log.Println("Application shutdown complete")
	return nil
}
