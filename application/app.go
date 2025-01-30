package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"githum.com/leebrouse/urlshortener/config"
	"githum.com/leebrouse/urlshortener/database"
	"githum.com/leebrouse/urlshortener/internal/api"
	"githum.com/leebrouse/urlshortener/internal/cache"
	"githum.com/leebrouse/urlshortener/internal/service"
	"githum.com/leebrouse/urlshortener/pkg/shortcode"
	"githum.com/leebrouse/urlshortener/pkg/validator"
)

type Application struct {
	e                   *echo.Echo
	db                  *sql.DB
	redisClient         *cache.RedisCache
	URLService          *service.URLService
	urlHandler          *api.URLHandler
	cfg                 *config.Config
	shortCondeGenerator *shortcode.ShortCode
}

func (a *Application) Init(filePath string) error {
	cfg, err := config.LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.cfg = cfg

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		return err
	}
	a.db = db

	redisClient, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		return err
	}
	a.redisClient = redisClient
	a.shortCondeGenerator = shortcode.NewShortCode(cfg.ShortCode)

	urlService := service.NewURLService(db, a.shortCondeGenerator, cfg.App.DefaultDuration, redisClient, cfg.App.BaseURL)

	a.urlHandler = api.NewURLHandler(urlService)

	e := echo.New()
	e.Server.WriteTimeout = cfg.Server.ReadTimeout
	e.Server.ReadTimeout = cfg.Server.WriteTimeout
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/api/url", a.urlHandler.CreateURL)
	e.GET("/:code", a.urlHandler.RedirectURL)
	e.Validator=validator.NewCustomValidator()
	a.e = e
	return nil
}

func (a *Application) Run() {
	go a.startServer()
	go a.cleanUp()
	a.shutdown()
}

func (a *Application) startServer() {
	if err := a.e.Start(a.cfg.Server.Addr); err != nil {
		log.Println(err)
	}
}

func (a *Application) cleanUp() {
	ticker := time.NewTicker(a.cfg.App.CleanupInterval)
	defer ticker.Stop()

	for _ = range ticker.C {
		if err := a.URLService.DeleteURL(context.Background()); err != nil {
			log.Println(err)
		}
	}
}

func (a *Application) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer func() {
		if err := a.db.Close(); err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		if err := a.redisClient.Close(); err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
