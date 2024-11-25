// @title Music info
// @version 1.0
// @BasePath /
package router

import (
	"context"
	_ "em/docs"
	"em/internal/adapters/api/action"
	"em/internal/adapters/logger"
	"em/internal/adapters/presenter"
	"em/internal/config"
	"em/internal/repo"
	"em/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ginEngine struct {
	cfg        config.Config
	router     *gin.Engine
	log        logger.Logger
	db         *gorm.DB
	ctxTimeout time.Duration
}

func newGinServer(
	cfg config.Config,
	log logger.Logger,
	db *gorm.DB,
	ctxTimeout time.Duration,
) *ginEngine {
	return &ginEngine{
		cfg:        cfg,
		router:     gin.New(),
		log:        log,
		db:         db,
		ctxTimeout: ctxTimeout,
	}
}

func (g *ginEngine) Listen() {
	gin.Recovery()

	g.setupRoutes(g.router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", g.cfg.AppPort),
		Handler: g.router,
	}

	go func() {
		g.log.Infof(fmt.Sprintf("Starting server on port %d", g.cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			g.log.Fatalln(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	g.log.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		g.log.Fatalln(fmt.Sprintf("Server forced to shutdown: %s", err))
	}

	g.log.Infof("Server exiting")
}

func (g ginEngine) setupRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/", g.buildSongCreateAction())
	router.DELETE("/:id", g.buildSongDeleteAction())
	router.PATCH("/:id", g.buildSongUpdateAction())
	router.GET("/info", g.buildSongFindAllAction())
	router.GET("/info/:id", g.buildSongFindTextVersesAction())

}

func buildParams(c *gin.Context, params ...string) {
	q := c.Request.URL.Query()

	for _, value := range params {
		switch value {
		case "page":
			if _, exists := q["page"]; !exists {
				q.Set("page", "1")
			}
		case "limit":
			if _, exists := q["limit"]; !exists {
				q.Set("limit", "10")
			}
		default:
			q.Add(value, c.Param(value))
		}
	}
	c.Request.URL.RawQuery = q.Encode()
}

// buildSongCreateAction создаёт новый музыкальный трек.
// @Summary Создание нового трека
// @Description Создает новый музыкальный трек в базе данных.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body usecase.SongCreateInput true "Данные о песне"
// @Success 201 {object} usecase.SongCreateOutput
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router / [post]
func (g *ginEngine) buildSongCreateAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewSongCreateInteractor(
				repo.NewLibraryRepository(g.db),
				presenter.NewSongCreatePresenter(),
				g.ctxTimeout,
			)
			act = action.NewSongCreateAction(uc, g.log)
		)
		act.Execute(c.Writer, c.Request)
	}
}

// buildSongDeleteAction удаляет музыкальный трек по ID.
// @Summary Удаление трека
// @Description Удаляет музыкальный трек по ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID трека"
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /{id} [delete]
func (g *ginEngine) buildSongDeleteAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewSongDeleteInteractor(
				repo.NewLibraryRepository(g.db),
				g.ctxTimeout,
			)
			act = action.NewSongDeleteAction(uc, g.log)
		)
		buildParams(c, "id")
		act.Execute(c.Writer, c.Request)
	}
}

// buildSongUpdateAction обновляет данные музыкального трека по ID.
// @Summary Обновление данных трека
// @Description Обновляет информацию о музыкальном треке по ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID трека"
// @Param song body usecase.SongUpdateInput true "Данные для обновления"
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /{id} [patch]
func (g *ginEngine) buildSongUpdateAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewSongUpdateInteractor(
				repo.NewLibraryRepository(g.db),
				g.ctxTimeout,
			)
			act = action.NewSongUpdateAction(uc, g.log)
		)
		buildParams(c, "id")

		act.Execute(c.Writer, c.Request)
	}
}

// buildSongFindAllAction возвращает список всех песен с параметрами пагинации и фильтрации.
// @Summary Получение списка песен
// @Description Возвращает список всех песен с возможностью пагинации и фильтрации по группе и названию.
// @Tags songs
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию"
// @Param orderBy query string false "Сортировка (формат: поле:направление, например: release_date:asc,song:desc). Поля: release_date, song, text"
// @Param text query string false "Фильтр по тексту"
// @Success 200 {array} usecase.SongFindAllOutput "Список песен"
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /info [get]
func (g *ginEngine) buildSongFindAllAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewSongFindAllInteractor(
				repo.NewLibraryRepository(g.db),
				presenter.NewSongFindAllPresenter(),
				g.ctxTimeout,
			)
			act = action.NewSongFindAllAction(uc, g.log)
		)
		buildParams(c, "page", "limit", "group", "song", "orderBy", "text")

		act.Execute(c.Writer, c.Request)
	}
}

// buildSongFindTextVersesAction возвращает текст песни по ID.
// @Summary Получение текста песни по ID
// @Description Возвращает текст песни по ID.
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {object} usecase.SongFindTextVersesOutput "Текст песни"
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /info/{id} [get]
func (g *ginEngine) buildSongFindTextVersesAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewSongFindTextVersesInteractor(
				repo.NewLibraryRepository(g.db),
				presenter.NewSongFindTextVersesPresenter(),
				g.ctxTimeout,
			)
			act = action.NewSongFindTextVersesAction(uc, g.log)
		)
		buildParams(c, "page", "limit", "id")

		act.Execute(c.Writer, c.Request)
	}
}
