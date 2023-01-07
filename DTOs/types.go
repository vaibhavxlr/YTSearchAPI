package DTOs

type Video struct {
	VidId       string      `json:"vidId"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	PublishDate string      `json:"publishDate"`
	Thumbnails  []Thumbnail `json:"thumbnails"`
}

type Thumbnail struct {
	URL string `json:"url"`
}
