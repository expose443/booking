package model

type Data struct {
	Posts        []Post
	Lol          string
	User         User
	Notification Notification
}

type Notification struct {
	Message  string
	CheckIn  string
	CheckOut string
	Room     string
}
