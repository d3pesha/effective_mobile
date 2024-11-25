package router

import (
	"em/internal/adapters/logger"
	"em/internal/config"
	"gorm.io/gorm"
	"time"
)

type Server interface {
	Listen()
}

func NewWebServer(
	cfg config.Config,
	logger logger.Logger,
	db *gorm.DB,
	ctxTimeout time.Duration,
) Server {
	return newGinServer(cfg, logger, db, ctxTimeout)
}
