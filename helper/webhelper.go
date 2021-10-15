package helper

type WebResponse struct {
	Draw            int         `json:"draw"`
	RecordsFiltered int         `json:"recordsFiltered"`
	Data            interface{} `json:"data"`
}

type Data struct {
	Number int    `json:"no"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

func WebApiResponse(draw int, recordsFiltered int, data interface{}) WebResponse {
	jsonResponse := WebResponse{
		Draw:            draw,
		RecordsFiltered: recordsFiltered,
		Data:            data,
	}
	return jsonResponse
}
