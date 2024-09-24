package responses

type URL struct {
	URL string `json:"url"`
}

func NewURL(ok string) URL {
	return URL{ok}
}
