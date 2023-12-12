package models

type HttpResponse struct {
	Error Error       `json:"error"`
	Data  interface{} `json:"data"`
}

type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (ht HttpResponse) NewResponse(error Error, data interface{}) HttpResponse {
	ht.Error = error
	ht.Data = data
	return ht
}
