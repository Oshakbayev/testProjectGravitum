package v1

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"template/internal/delivery/http/middleware"
)

type Handler struct {
	logger    *zap.Logger
	responder Responder
	authSvc   authSvc
	userSvc   userSvc
}

func NewHandler(
	responder Responder,
	authSvc authSvc,
	userSvc userSvc,
	logger *zap.Logger) *Handler {
	return &Handler{
		responder: responder,
		authSvc:   authSvc,
		userSvc:   userSvc,
		logger:    logger,
	}
}

func SetHandler(
	router *chi.Mux,
	responder Responder,
	authSvc authSvc,
	userSvc userSvc,
	logger *zap.Logger,
) {
	handler := NewHandler(responder, authSvc, userSvc, logger)
	authMiddleware := middleware.NewAuthMiddleware(userSvc, authSvc, responder, logger).AuthorizedJWT
	authorizedAdminMdl := middleware.NewAuthMiddleware(userSvc, authSvc, responder, logger).AuthorizedAdmin
	recoveryMiddleware := middleware.NewRecoveryMiddleware(responder, logger).Recover
	setRoutes(handler, router, authMiddleware, authorizedAdminMdl, recoveryMiddleware)
}
