package v1

import (
	"errors"
	"github.com/go-chi/chi/v5"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"net/http"
	_ "template/internal/delivery/http"
	"template/internal/dto/user"
	"template/internal/utils"
	"time"
)

// CreateUser
//
//	@Summary		creates new user
//	@Tags			users
//	@Description	creates new user.
//	@Produce		json
//	@Accept			json
//	@Param			user	body		user.CreateUserRequest	true	"User details"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	http.myResponse
//	@Failure		401		{object}	http.myResponse
//	@Failure		403		{object}	http.myResponse
//	@Failure		500		{object}	http.myResponse
//	@Security		Bearer
//	@Router			/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	defer r.Body.Close()

	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.WithBadRequest(w, "INVALID_REQUEST_BODY")
		return
	}

	if err := req.Validate(); err != nil {
		h.responder.WithBadRequest(w, err.Error())
		return
	}

	logger := h.logger.With(
		zap.String("method", "CreateUser"),
	)

	userId, err := h.userSvc.CreatUser(r.Context(), req)
	var expError utils.ExplorerError

	if errors.As(err, &expError) {
		switch expError.ErrorCategory {
		case utils.ErrorCategoryInternal:
			uuid := h.responder.WithInternalError(w, "handled internal error")
			logger.Error("Failed to create user",
				zap.String("type", "ERROR"),
				zap.String("error_type", "CREATE_USER_FAILED"),
				zap.String("request_status", "fail"),
				zap.String("error_uuid", uuid),
				zap.Duration("processing_time", time.Since(startTime)),
				zap.String("reason", err.Error()),
			)
		case utils.ErrorCategoryUserError:
			h.responder.With(expError.StatusCode, w, expError.Message)
		}

		return
	} else if err != nil {
		uuid := h.responder.WithInternalError(w, "unhandled internal error")
		logger.Error("Failed create user",
			zap.String("type", "ERROR"),
			zap.String("error_type", "CREATE_USER_FAILED"),
			zap.String("request_status", "fail"),
			zap.String("error_uuid", uuid),
			zap.Duration("processing_time", time.Since(startTime)),
			zap.String("reason", err.Error()),
		)
		return
	}
	h.responder.WithOK(w, userId)

	logger.Info("Created user",
		zap.String("type", "AUTH_LOGIN"),
		zap.String("request_status", "success"),
		zap.String("created_userId", userId),
		zap.Int("level", 2),
		zap.Duration("processing_time", time.Since(startTime)),
	)
}

// UpdateUser
//
//	@Summary		updates user
//	@Tags			users
//	@Description	updates user
//	@Produce		json
//	@Accept			json
//	@Param			user	body		user.UpdateUserRequest	true	"User details"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	http.myResponse
//	@Failure		401		{object}	http.myResponse
//	@Failure		403		{object}	http.myResponse
//	@Failure		500		{object}	http.myResponse
//	@Security		Bearer
//	@Router			/users [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	defer r.Body.Close()

	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.WithBadRequest(w, "INVALID_REQUEST_BODY")
		return
	}
	userID := chi.URLParam(r, "id")
	if err := req.Validate(); err != nil {
		h.responder.WithBadRequest(w, err.Error())
		return
	}

	logger := h.logger.With(
		zap.String("method", "UpdateUser"),
	)

	err := h.userSvc.UpdateUser(r.Context(), userID, req)
	var expError utils.ExplorerError

	if errors.As(err, &expError) {
		switch expError.ErrorCategory {
		case utils.ErrorCategoryInternal:
			uuid := h.responder.WithInternalError(w, "handled internal error")
			logger.Error("Failed update for ",
				zap.String("type", "ERROR"),
				zap.String("error_type", "UPDATE_USER_FAILED"),
				zap.String("request_status", "fail"),
				zap.String("error_uuid", uuid),
				zap.Duration("processing_time", time.Since(startTime)),
				zap.String("reason", err.Error()),
			)
		case utils.ErrorCategoryUserError:
			h.responder.With(expError.StatusCode, w, expError.Message)
		}

		return
	} else if err != nil {
		uuid := h.responder.WithInternalError(w, "unhandled internal error")
		logger.Error("Failed update for user ",
			zap.String("type", "ERROR"),
			zap.String("error_type", "UPDATE_USER_FAILED"),
			zap.String("request_status", "fail"),
			zap.String("error_uuid", uuid),
			zap.Duration("processing_time", time.Since(startTime)),
			zap.String("reason", err.Error()),
		)
		return
	}
	h.responder.WithOK(w, nil)

	logger.Info("Edit User for ",
		zap.String("type", "UPDATE_USER_SUCCESS"),
		zap.String("request_status", "success"),
		zap.Duration("processing_time", time.Since(startTime)),
	)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := h.logger.With(
		zap.String("method", "GetUser"),
	)

	userID := chi.URLParam(r, "id")
	me, err := h.userSvc.GetUserById(r.Context(), userID)
	var expError utils.ExplorerError

	if errors.As(err, &expError) {
		switch expError.ErrorCategory {
		case utils.ErrorCategoryInternal:
			uuid := h.responder.WithInternalError(w, "handled internal error")
			logger.Error("Failed get user",
				zap.String("type", "ERROR"),
				zap.String("error_type", "GET_USER_FAILED"),
				zap.String("request_status", "fail"),
				zap.String("error_uuid", uuid),
				zap.Duration("processing_time", time.Since(startTime)),
				zap.String("reason", err.Error()),
			)
		case utils.ErrorCategoryUserError:
			h.responder.With(expError.StatusCode, w, expError.Message)
		}

		return
	} else if err != nil {
		uuid := h.responder.WithInternalError(w, "unhandled internal error")
		logger.Error("Failed get user",
			zap.String("type", "ERROR"),
			zap.String("error_type", "GET_USER_FAILED"),
			zap.String("request_status", "fail"),
			zap.String("error_uuid", uuid),
			zap.Duration("processing_time", time.Since(startTime)),
			zap.String("reason", err.Error()),
		)
		return
	}
	h.responder.WithOK(w, me)

	logger.Info(" Get user",
		zap.String("type", "GET_USER_SUCCESS"),
		zap.String("request_status", "success"),
		zap.Duration("processing_time", time.Since(startTime)),
	)
}
