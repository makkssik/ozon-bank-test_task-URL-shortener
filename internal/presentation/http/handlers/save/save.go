package save

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/presentation/http/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLShortener --output=./mocks --outpkg=mocks --filename=url_shortener.go
type URLShortener interface {
	ShortenURL(ctx context.Context, originalURL string) (*entities.URL, error)
}

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func NewSaveHandler(log *slog.Logger, shortener URLShortener) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.new"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Info("invalid request format", slog.String("error", err.Error()))

			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		urlEntity, err := shortener.ShortenURL(r.Context(), req.URL)
		if err != nil {
			log.Error("failed to shorten url", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to shorten url"))
			return
		}

		log.Info("url shortened", slog.String("alias", urlEntity.Short))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    urlEntity.Short,
		})
	}
}
