package qichacha

type Result struct {
	Status int `json:"Status"`
	Result []struct {
		HitReason struct {
			Field string `json:"Field"`
			Value string `json:"Value"`
		} `json:"HitReason"`
		HitReasons []struct {
			Field string `json:"Field"`
			Value string `json:"Value"`
		} `json:"HitReasons"`
		KeyNo         string        `json:"KeyNo"`
		Name          string        `json:"Name"`
		Type          int           `json:"Type"`
		No            string        `json:"No"`
		CreditCode    string        `json:"CreditCode"`
		OperName      string        `json:"OperName"`
		Status        string        `json:"Status"`
		StartDate     string        `json:"StartDate"`
		Address       string        `json:"Address"`
		RegistCapi    string        `json:"RegistCapi"`
		ContactNumber string        `json:"ContactNumber"`
		Email         string        `json:"Email"`
		ImageUrl      string        `json:"ImageUrl"`
		TelList       string        `json:"TelList"`
		IsAuth        bool          `json:"IsAuth"`
		OperInfo      string        `json:"OperInfo"`
		TagsInfo      interface{}   `json:"TagsInfo"`
		Flag          string        `json:"Flag"`
		Tag           string        `json:"Tag"`
		EmailList     []interface{} `json:"EmailList"`
		AuthLevel     int           `json:"AuthLevel"`
		IsHide        bool          `json:"IsHide"`
		AreaCode      string        `json:"AreaCode"`
		QccCode       string        `json:"QccCode"`
		GW            string        `json:"GW"`
		CountInfo     interface{}   `json:"CountInfo"`
		Industry      struct {
			IndustryCode    string `json:"IndustryCode"`
			Industry        string `json:"Industry"`
			SubIndustryCode string `json:"SubIndustryCode"`
			SubIndustry     string `json:"SubIndustry"`
		} `json:"Industry"`
		EconKind          string   `json:"EconKind"`
		ShortStatus       string   `json:"ShortStatus"`
		UpdatedDate       string   `json:"UpdatedDate"`
		OperId            string   `json:"OperId"`
		OperOrg           int      `json:"OperOrg"`
		RegistCapitalType string   `json:"RegistCapitalType"`
		LegalManType      string   `json:"legalManType"`
		StatusColor       string   `json:"StatusColor"`
		StockInfo         []string `json:"StockInfo"`
		StockCode         string   `json:"StockCode"`
		Tags              []struct {
			Tag   string `json:"tag"`
			Color string `json:"color"`
		} `json:"tags"`
		ShowFields []struct {
			K   string `json:"k"`
			V   string `json:"v"`
			Url string `json:"url,omitempty"`
		} `json:"showFields"`
		TelCount    int    `json:"TelCount"`
		City        string `json:"City"`
		PhoneStatus struct {
			StatusIcon string `json:"statusIcon"`
		} `json:"PhoneStatus"`
		GWStatus struct {
			StatusIcon string `json:"statusIcon"`
		} `json:"GWStatus"`
		HitReasonRouter string `json:"HitReasonRouter"`
		HasFollow       bool   `json:"hasFollow"`
		ShowFollow      bool   `json:"showFollow"`
		CardLabel       string `json:"CardLabel"`
	} `json:"Result"`
	Paging struct {
		PageSize     int `json:"PageSize"`
		PageIndex    int `json:"PageIndex"`
		TotalRecords int `json:"TotalRecords"`
	} `json:"Paging"`
	ShowTotal      string `json:"ShowTotal"`
	ActivityRouter string `json:"ActivityRouter"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  Result `json:"result"`
}
