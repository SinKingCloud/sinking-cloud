package qiyoucha

type RespBody struct {
	Total int `json:"total"`
	Items []struct {
		Name            string `json:"name"`
		LegalPersonName string `json:"legal_person_name"`
		Code            string `json:"code"`
		EstiblishTime   string `json:"estiblish_time"`
		IsCertified     int    `json:"is_certified"`
	} `json:"items"`
}

type Response struct {
	RespCode int      `json:"respCode"`
	RespBody RespBody `json:"respBody"`
}
