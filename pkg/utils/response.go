package utils

import (
	"errors"
	"net/http"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/cloudwego/hertz/pkg/app"

	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// WriteHertzErrorResponse for error response.
func WriteHertzErrorResponse(c *app.RequestContext, err error) {
	appError := new(apperrors.AppError)
	if !errors.As(err, &appError) {
		applog.Errorw("not apperror", "err", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	switch appError.Code {
	case apperrors.InvalidCode, apperrors.CannotExecCode, apperrors.BindErrorCode:
		c.JSON(http.StatusBadRequest, appError.Message)
	case apperrors.NotFoundCode, apperrors.RouteNotFoundCode:
		c.JSON(http.StatusNotFound, appError.Message)
	case apperrors.InternalCode:
		c.JSON(http.StatusInternalServerError, appError.Message)
	default:
		applog.Errorw("internal error", "err", appError.Inner)
		c.JSON(http.StatusInternalServerError, appError.Message)
	}
}

// WriteHertzOKResponse for all success request.
func WriteHertzOKResponse(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, data)
}
