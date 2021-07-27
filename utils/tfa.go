package utils

type TFARequest struct {
	Data Data `json:"data"`
}

type Data struct {
	Id      string `json:"id"`
	Project string `json:"project"`
}

type TFAResponse struct {
	Prediction string `json:"prediction"`
}
