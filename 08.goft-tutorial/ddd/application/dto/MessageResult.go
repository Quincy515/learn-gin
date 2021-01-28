package dto

type MessageResult struct {
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}
