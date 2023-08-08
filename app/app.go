package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/middleware"
	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/router"
	"lim/tools/log"
)

type App struct {
	s *http.Server
	r *gin.Engine
}

func New() (*App, error) {
	err := config.Init()
	if err != nil {
		return nil, err
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), middleware.Logger(), middleware.Cors(), middleware.AuthToken())
	router.Set(r.Group("/api"))

	return &App{
		r: r,
	}, nil
}

func (a *App) Run() (err error) {
	sig := make(chan os.Signal, 0)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = a.run(); err != nil {
			log.Errorf("服务运行失败: %s", err.Error())
			sig <- syscall.SIGTERM
		}
	}()

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = a.s.Shutdown(ctx); err != nil {
		log.Errorf("服务关闭失败: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Errorf("Mongo关闭失败: %s", err.Error())
	}
	if err = cache.Close(); err != nil {
		log.Errorf("Redis关闭失败: %s", err.Error())
	}

	return
}

func (a *App) run() error {
	addr, tlsCfg := config.GetApp().Addr, config.GetApp().TLS
	a.s = &http.Server{
		Addr:    addr,
		Handler: a.r,
	}

	if tlsCfg == nil || tlsCfg.CertFile == "" || tlsCfg.KeyFile == "" {
		return a.s.ListenAndServe()
	}

	return a.s.ListenAndServeTLS(tlsCfg.CertFile, tlsCfg.KeyFile)
}

func Register(rg *gin.RouterGroup, c *config.Config) {
	config.SetFromManual(c)

	rg.Use(middleware.AuthToken())
	router.Set(rg)
}
