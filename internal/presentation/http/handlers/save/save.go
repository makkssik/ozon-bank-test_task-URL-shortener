package save

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/application/contracts/url"
	"url-shortener/internal/application/contracts/url/operations"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/presentation/http/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func NewSaveHandler(log *slog.Logger, service url.Service) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.new"

		log = log.With(
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

		res, err := service.ShortenURL(r.Context(), operations.ShortenRequest{
			OriginalURL: req.URL,
		})

		if err != nil {
			if errors.Is(err, exceptions.ErrAlreadyExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.Status(r, http.StatusConflict)
				render.JSON(w, r, response.Error("url already exists"))
				return
			}

			log.Error("failed to shorten url", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to shorten url"))
			return
		}

		log.Info("url shortened", slog.String("alias", res.ShortURL))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    res.ShortURL,
		})
	}
}
