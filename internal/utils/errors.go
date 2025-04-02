package utils

import "net/http"

type ErrorCategory uint8

const (
	ErrorCategoryInternal  ErrorCategory = 1
	ErrorCategoryUserError ErrorCategory = 2
)

type ExplorerError struct {
	ErrorCategory ErrorCategory
	StatusCode    int
	Message       string
}

func (e ExplorerError) Error() string {
	return e.Message
}

func NewInternalError(message string) ExplorerError {
	return ExplorerError{
		ErrorCategory: ErrorCategoryInternal,
		StatusCode:    http.StatusInternalServerError,
		Message:       message,
	}
}

func NewUserError(status int, message string) ExplorerError {
	return ExplorerError{
		ErrorCategory: ErrorCategoryUserError,
		StatusCode:    status,
		Message:       message,
	}
}

var (
	ErrUnauthorized            = NewUserError(http.StatusUnauthorized, "unauthorized")
	ErrForbidden               = NewUserError(http.StatusForbidden, "ACTION_FORBIDDEN")
	ErrNotFound                = NewUserError(http.StatusNotFound, "not found")
	ErrWrongCredentials        = NewUserError(http.StatusBadRequest, "WRONG_CREDENTIALS")
	ErrUserBlocked             = NewUserError(http.StatusForbidden, "USER_BLOCKED")
	ErrUserNotFound            = NewUserError(http.StatusBadRequest, "USER_NOT_FOUND")
	ErrInvalidToken            = NewUserError(http.StatusUnauthorized, "INVALID_TOKEN")
	ErrWrongUserId             = NewUserError(http.StatusBadRequest, "WRONG_USER_ID")
	ErrWrongDashboardId        = NewUserError(http.StatusBadRequest, "WRONG_DASHBOARD_ID")
	ErrDashboardNotFound       = NewUserError(http.StatusNotFound, "DASHBOARD_NOT_FOUND")
	ErrDashboardNotARRESTED    = NewUserError(http.StatusNotFound, "DASHBOARD_NOT_ARRESTED")
	ErrActivityNotFound        = NewUserError(http.StatusNotFound, "ACTIVITY_NOT_FOUND")
	ErrNotUserCreator          = NewUserError(http.StatusForbidden, "NOT_USER_CREATOR")
	ErrAlreadyExists           = NewUserError(http.StatusConflict, "USERALREADY_EXISTS")
	ErrInvalidEmail            = NewUserError(http.StatusBadRequest, "INVALID_EMAIL")
	ErrPasswordsNotEqual       = NewUserError(http.StatusBadRequest, "PASSWORDS_NOT_EQUAL")
	ErrEmailAlreadyExists      = NewUserError(http.StatusBadRequest, "EMAIL_ALREADY_EXISTS")
	ErrUsernameAlreadyExists   = NewUserError(http.StatusBadRequest, "USERNAME_ALREADY_EXISTS")
	ErrPasswordNotMatch        = NewUserError(http.StatusBadRequest, "PASSWORDS_NOT_MATCH")
	ErrInternalRole            = NewInternalError("ROLE_NOT_FOUNDED")
	ErrVideoLimit              = NewUserError(http.StatusBadRequest, "Вы превысили лимит загружаемых видео в день. ")
	ErrImageLimit              = NewUserError(http.StatusBadRequest, "Вы превысили лимит загружаемых фотографий в день.")
	ErrParallelFileLimit       = NewUserError(http.StatusBadRequest, "Вы превысили лимит параллельно обрабатываемых медиа. ")
	ErrTooManyRequests         = NewUserError(http.StatusTooManyRequests, "too many requests")
	ErrStatusBadRequest        = NewUserError(http.StatusBadRequest, "bad request")
	ErrStatusExternalServer    = NewInternalError("external server error")
	ErrUnknownStatus           = NewInternalError("unknown status")
	ErrPersonNotFound          = NewUserError(http.StatusNotFound, "PERSON_NOT_FOUND")
	ErrReportNotFound          = NewUserError(http.StatusNotFound, "REPORT_NOT_FOUND")
	ErrDataLookUpFieldNotFound = NewInternalError("DATA_LOOKUP_FIELD_NOT_FOUND")
	ErrWrongPersonId           = NewUserError(http.StatusBadRequest, "WRONG_PERSON_ID")
	ErrWrongFileId             = NewUserError(http.StatusBadRequest, "WRONG_FILE_ID")
	ErrNoTransferFound         = NewUserError(http.StatusBadRequest, "No transfer found")
	ErrCameraNotFound          = NewUserError(http.StatusBadRequest, "CAMERA_NOT_FOUND")
	ErrCameraVectorNotFound    = NewUserError(http.StatusBadRequest, "CAMERA_VECTOR_NOT_FOUND")
	ErrCameraStatus            = NewUserError(http.StatusBadRequest, "CAMERA_STATUS_IS_DELETED_OR_CAMERA_ALREADY_STREAMS")
	ErrWrongCameraId           = NewUserError(http.StatusBadRequest, "WRONG_CAMERA_ID")
	ErrWrongPid                = NewInternalError("WRONG_PID")
	ErrCameraAlreadyStreams    = NewInternalError("CAMERA_ALREADY_STREAMS")
	ErrWrongCommentID          = NewUserError(http.StatusBadRequest, "WRONG_COMMENT_ID")
	ErrAnalyzerWebsocket       = NewInternalError("ANALYZER_WEBSOCKET")
	ErrInvalidUrl              = NewUserError(http.StatusBadRequest, "INVALID_URL")
	ErrActiveZone              = NewUserError(http.StatusBadRequest, "INVALID_ACTIVE_ZONE")
	ErrBlindZone               = NewUserError(http.StatusBadRequest, "INVALID_ACTIVE_ZONE")
	ErrFindMetaDataDocs        = NewUserError(http.StatusNotFound, "PERSON_BY_METADATA_NOT_FOUND")
	ErrInvalidMetaType         = NewUserError(http.StatusBadRequest, "INVALID_META_TYPE")
	ErrInvalidMetaValue        = NewUserError(http.StatusBadRequest, "INVALID_META_VALUE")
)
