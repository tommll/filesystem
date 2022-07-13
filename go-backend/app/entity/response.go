package entity

type BaseDataResponse struct {
	ReturnCode int         `json:"return_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}
