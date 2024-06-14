package models

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type RequestBatch struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type ResponseBatch struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
