package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AleksandrVishniakov/jwt-auth/internal/e"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) (err error)

func Error(next ErrorHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err == nil {
			return
		}

		var httpError *e.HTTPError
		if !errors.As(err, &httpError) {
			slog.Warn("got non-wrapped service error", e.SlogErr(err))
			httpError = e.NewError(e.WithError(err))
		}

		slog.Debug("http error occured", e.SlogErr(httpError.Unwrap()))

		err = EncodeResponse(w, httpError, httpError.Code)
		if err != nil {
			slog.Error("encoding response error", e.SlogErr(err))
		}
	})
}
