package infrastructure

import (
	"em/internal/adapters/logger"
	"em/internal/config"
	"em/internal/infrastructure/database"
	"em/internal/infrastructure/log"
	"em/internal/infrastructure/router"
	"gorm.io/gorm"
	"time"
)

type app struct {
	cfg        config.Config
	logger     logger.Logger
	db         *gorm.DB
	webServer  router.Server
	ctxTimeout time.Duration
}

func NewConfig(config config.Config) *app {
	return &app{
		cfg: config,
	}
}

func (a *app) ContextTimeout(t time.Duration) *app {
	a.ctxTimeout = t
	return a
}

func (a *app) Logger() *app {
	log := log.NewLogrusLogger()
	a.logger = log
	a.logger.Infof("Success log configured")

	return a
}
func (a *app) Database() *app {
	db, err := database.NewGormHandler(&a.cfg)
	if err != nil {
		a.logger.Fatalln(err, "Could not make a connection database")
	}
	a.db = db
	a.logger.Infof("Success connected to database")

	return a
}

func (a *app) WebServer() *app {
	server := router.NewWebServer(
		a.cfg,
		a.logger,
		a.db,
		a.ctxTimeout,
	)

	a.logger.Infof("Success connected to web server")
	a.webServer = server
	return a
}

func (a *app) Start() {
	a.webServer.Listen()
}
