package v1

import (
	"errors"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
	_ "template/internal/delivery/http"
	"template/internal/dto/auth"
	"template/internal/utils"
	"time"
)

// Login
//
//	@Summary		login
//	@Tags			auth
//	@Description	login
//	@Produce		json
//	@Accept			json
//	@Param			user	body		auth.LoginRequest	true	"user creds"
//	@Success		200		{object}	auth.LoginResponse
//	@Failure		401		{object}	http.myResponse
//	@Failure		403		{object}	http.myResponse
//	@Failure		500		{object}	http.myResponse
//	@Router			/auth/local [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.WithBadRequest(w, "INVALID_REQUEST_BODY")
		return
	}

	if err := req.Validate(); err != nil {
		h.responder.WithBadRequest(w, err.Error())
		return
	}

	logger := h.logger.With(
		zap.String("method", "login"),
	)

	res, err := h.authSvc.Login(r.Context(), req)
	var expError utils.ExplorerError

	if errors.As(err, &expError) {
		switch expError.ErrorCategory {
		case utils.ErrorCategoryInternal:
			uuid := h.responder.WithInternalError(w, "internal error")
			logger.Error("Failed authorization",
				zap.String("type", "ERROR"),
				zap.String("error_type", "AUTH_ERROR"),
				zap.String("request_status", "fail"),
				zap.String("error_uuid", uuid),
				zap.Duration("processing_time", time.Since(startTime)),
				zap.String("reason", err.Error()),
			)
		case utils.ErrorCategoryUserError:
			h.responder.With(expError.StatusCode, w, expError.Message)
			logger.Error("Failed authorization",
				zap.String("type", "ERROR"),
				zap.String("error_type", "AUTH_ERROR"),
				zap.Duration("processing_time", time.Since(startTime)),
				zap.String("reason", err.Error()))
		}

		return
	} else if err != nil {
		uuid := h.responder.WithInternalError(w, "internal error")
		logger.Error("Failed authorization for",
			zap.String("type", "ERROR"),
			zap.String("error_type", "AUTH_ERROR"),
			zap.String("request_status", "fail"),
			zap.String("error_uuid", uuid),
			zap.Duration("processing_time", time.Since(startTime)),
			zap.String("reason", err.Error()),
		)
		return
	}

	logger.Info("Successful authorization for "+res.User.Name,
		zap.String("user_id", res.User.ID),
		zap.Int("level", 2),
		zap.String("type", "AUTH_SUCCESS"),
		zap.String("request_status", "success"),
		zap.Duration("processing_time", time.Since(startTime)),
	)
	h.responder.WithOK(w, res)
}
