package entity

type Article struct {
	ID      int64  `json:"id" bson:"id"`
	Author  string `json:"author" bson:"author"`
	Title   string `json:"title" bson:"title"`
	Body    string `json:"body" bson:"body"`
	Created string `json:"created" bson:"created"`
}
