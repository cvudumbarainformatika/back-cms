package bootstrap

import (
	"fmt"
	"log"

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
	Redis  *redis.Client
	Config *config.Config
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

	// Setup Redis connection
	rdb, err := database.InitRedis(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	// Setup routes
	routes.SetupRoutes(router, db.DB, rdb, cfg)

	return &Application{
		Router: router,
		DB:     db,
		Redis:  rdb,
		Config: cfg,
	}, nil
}

// Run starts the application server
func (app *Application) Run() error {
	addr := fmt.Sprintf(":%s", app.Config.App.Port)
	log.Printf("Starting %s on %s (env: %s)", app.Config.App.Name, addr, app.Config.App.Env)
	return app.Router.Run(addr)
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
