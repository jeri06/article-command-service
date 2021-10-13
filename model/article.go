package model

type Article struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
