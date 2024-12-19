package delivery

import (
	"context"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
	_ "todo/docs"
	"todo/internal/application"
)

type Config struct {
	Port         string        `env:"PORT"`
	ReadTimeOut  time.Duration `env:"READ_TIMEOUT"`
	WriteTimeOut time.Duration `env:"WRITE_TIMEOUT"`
}

type Handler struct {
	cfg        *Config
	services   *application.Service
	httpServer *http.Server
	router     *gin.Engine
}

func NewHandler(services *application.Service, cfg *Config) *Handler {
	return &Handler{
		services: services,
		cfg:      cfg,
	}
}

func (h *Handler) Run(_ context.Context) error {

	h.httpServer = &http.Server{
		Addr:         ":" + h.cfg.Port,
		Handler:      h.router,
		ReadTimeout:  h.cfg.ReadTimeOut,
		WriteTimeout: h.cfg.WriteTimeOut,
	}
	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			log.Println("listen: %s\n", err.Error())
			return
		}
	}()
	return nil
}

func (h *Handler) Stop(ctx context.Context) error {
	return h.httpServer.Shutdown(ctx)
}

func (h *Handler) Init() error {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}

		items := api.Group("items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	h.router = router
	return nil
}
