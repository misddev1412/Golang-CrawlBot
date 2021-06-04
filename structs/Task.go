package structs

type Task struct {
	CurrentPage int64 `json:"current_page"`
	Data        []struct {
		FeedURL string `json:"feedURL"`
	} `json:"data"`
}

type DataTask struct {
	FeedURL string
	Feed    FeedData
}

type FeedData struct {
	Path string
}
