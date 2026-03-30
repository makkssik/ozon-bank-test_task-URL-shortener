package operations

type ShortenRequest struct {
	OriginalURL string
}

type ShortenResponse struct {
	ShortURL string
}

type GetOriginalRequest struct {
	ShortURL string
}

type GetOriginalResponse struct {
	OriginalURL string
}
