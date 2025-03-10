package domain

type User struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
	Age      int    `json:"Age"`
}

type Kafka struct {
	ID   int  `json:"id"`
	Data User `json:"data"`
}
