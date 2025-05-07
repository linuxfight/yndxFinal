package dto

type CalculateRequest struct {
	Expression string `json:"expression" validate:"expression,required" example:"2+2"`
}

type CalculateResponse struct {
	Id string `json:"id" example:"01JTE60CDWQ5R3QSWZBP8J6FG3"`
}
