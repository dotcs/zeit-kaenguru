package internal

type Comic struct {
	Id    int      `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Img   ComicImg `json:"img"`
}

type ComicImg struct {
	Src    string  `json:"src"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Ratio  float32 `json:"ratio"`
}
