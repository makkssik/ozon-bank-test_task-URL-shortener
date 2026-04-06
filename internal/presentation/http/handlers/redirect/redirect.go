package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/domain/entities"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/presentation/http/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLProvider --output=./mocks --outpkg=mocks --filename=url_provider.go
type URLProvider interface {
	GetOriginal(ctx context.Context, shortURL string) (*entities.URL, error)
}

func NewRedirectHandler(log *slog.Logger, provider URLProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.new"

		log := log.With(
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

		urlEntity, err := provider.GetOriginal(r.Context(), alias)

		if err != nil {
			if errors.Is(err, exceptions.ErrInvalidChars) {
				log.Info("invalid alias format", slog.String("alias", alias))
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, response.Error("invalid alias format"))
				return
			}

			log.Error("failed to get original url", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		if urlEntity == nil {
			log.Info("url not found", slog.String("alias", alias))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error("not found"))
			return
		}

		log.Info("got url", slog.String("url", urlEntity.Original))

		http.Redirect(w, r, urlEntity.Original, http.StatusFound)
	}
}
