package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/application/contracts/url"
	"url-shortener/internal/application/contracts/url/operations"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/presentation/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRedirectHandler(log *slog.Logger, service url.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.new"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		res, err := service.GetOriginal(r.Context(), operations.GetOriginalRequest{
			ShortURL: alias,
		})

		if err != nil {
			if errors.Is(err, exceptions.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, response.Error("not found"))
				return
			}

			log.Error("failed to get original url", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", res.OriginalURL))

		http.Redirect(w, r, res.OriginalURL, http.StatusFound)
	}
}
