package entity

type Article struct {
	ID      int64  `bson:"id"`
	Author  string `bson:"author"`
	Title   string `bson:"title"`
	Body    string `bson:"body"`
	Created string `bson:"created"`
}
