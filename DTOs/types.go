package DTOs

type Video struct {
	VidId       string `json:"vidId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishDate string `json:"publishDate"`
	Thumbnails  string `json:"thumbnails"`
}

type ErrorResp struct {
	ErrorMsg string `json:"errMsg"`
}

type SearchReqBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ExactMatch  bool   `json:"exactMatch"`
}
