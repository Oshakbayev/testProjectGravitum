package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
	_ "template/docs"
	"template/internal/config"
	http2 "template/internal/delivery/http"
	v1 "template/internal/delivery/http/v1"
	"template/internal/repository/mongo"
	"template/internal/service/auth"
	"template/internal/service/user"
	"time"
)

type App struct {
	cfg    *config.Config
	router *chi.Mux
	logger *zap.Logger
	db     db
}

func NewApp(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Initialize() error {

	a.logger, _ = zap.NewProduction()
	a.logger = a.logger.With(zap.String("facility", a.cfg.Facility))
	a.router = chi.NewRouter()
	if err := a.setHandler(); err != nil {
		a.logger.Error(err.Error())
		return err
	}
	return nil
}

func (a *App) setHandler() error {
	var err error

	mongoCli, err := mongo.NewMongoClient(a.cfg.Repository.Mongo)
	if err != nil {
		return err
	}
	a.db = mongoCli
	authRepo := mongo.NewAuthRepo(a.cfg.Repository.Mongo, mongoCli)
	userRepo := mongo.NewUserRepo(a.cfg.Repository.Mongo, mongoCli)

	authSvc := auth.NewService(authRepo, a.cfg.AuthConf, a.logger)
	userSvc := user.NewService(userRepo, a.logger)

	err = userSvc.InitAdmin(a.cfg.App.AdmConf)
	if err != nil {
		return err
	}

	responder := http2.NewResponder(a.logger)

	a.router.Get("/swagger/*", httpSwagger.WrapHandler)
	v1.SetHandler(a.router, responder, authSvc, userSvc, a.logger)

	return nil
}

func (a *App) Run(ctx context.Context) {
	var err error
	defer a.closeConnections()

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:        a.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    a.cfg.App.RTO,
		WriteTimeout:   a.cfg.App.WTO,
	}
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			a.logger.Info(err.Error())
			cancel()
			return
		}
	}()
	a.logger.Info("started", zap.Int("port", a.cfg.App.Port))
	<-ctx.Done()
	a.logger.Info("shutting down server ...\n")

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.Info("server forced shutdown ...")
	}

	a.logger.Info("server exiting ...")
}

func (a *App) closeConnections() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer a.logger.Sync()
	if err := a.db.Disconnect(ctx); err != nil {
		a.logger.Error(err.Error())
	}
	// если будет бд и т.д, закрывать коннект здесь
}
