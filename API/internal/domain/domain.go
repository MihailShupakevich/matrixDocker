package domain

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}
