package internal

type Comic struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	ImgSrc string `json:"imgSrc"`
	Date   string `json:"date"`
}
