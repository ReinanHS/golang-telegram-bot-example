package dto

type Data struct {
	Data interface{} `json:"data"`
}

type Response struct {
	Version string `json:"version"`
	Data
}

type ResponseMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
