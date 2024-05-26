package models

type Response struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

type Meta struct {
	Message string `json:"message"`
}
