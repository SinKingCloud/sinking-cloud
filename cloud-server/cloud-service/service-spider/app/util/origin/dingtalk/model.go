package dingtalk

type Data struct {
	Total int `json:"total"`
	Data  []struct {
		CompanyName            string      `json:"companyName"`
		CompanyNameWithSummary string      `json:"companyNameWithSummary"`
		LegalPerson            string      `json:"legalPerson"`
		RegisterDate           string      `json:"registerDate"`
		Address                string      `json:"address"`
		Ocid                   string      `json:"ocid"`
		Lon                    float64     `json:"lon"`
		Lat                    float64     `json:"lat"`
		RegisterCapital        string      `json:"registerCapital"`
		Logo                   interface{} `json:"logo"`
		StaffSize              interface{} `json:"staffSize"`
		Tags                   []string    `json:"tags"`
		BizStatus              string      `json:"bizStatus"`
		BAddMonitor            bool        `json:"bAddMonitor"`
		BFavorite              bool        `json:"bFavorite"`
		BTraced                interface{} `json:"bTraced"`
		Industry               string      `json:"industry"`
		IndustryLevel1         interface{} `json:"industryLevel1"`
		IndustryLevel2         interface{} `json:"industryLevel2"`
		Phone                  string      `json:"phone"`
		Email                  string      `json:"email"`
		Garden                 interface{} `json:"garden"`
		Brands                 interface{} `json:"brands"`
		OfficialSite           string      `json:"officialSite"`
		FinDate                interface{} `json:"finDate"`
		FinType                interface{} `json:"finType"`
		FinAmount              interface{} `json:"finAmount"`
		Investor               interface{} `json:"investor"`
		FlushTime              int         `json:"flushTime"`
		TaxNum                 string      `json:"taxNum"`
		SocialCreditCode       string      `json:"socialCreditCode"`
		HitSummary             bool        `json:"hitSummary"`
		SummaryFieldName       string      `json:"summaryFieldName"`
		SummaryFieldValue      string      `json:"summaryFieldValue"`
		IdentityId             string      `json:"identityId"`
		LatestBiddingItem      interface{} `json:"latestBiddingItem"`
		BiddingCnt             interface{} `json:"biddingCnt"`
		ParkId                 interface{} `json:"parkId"`
		MaxConnectionTelephone interface{} `json:"maxConnectionTelephone"`
		TotalTelephoneNum      interface{} `json:"totalTelephoneNum"`
		TelephoneDetailLists   []struct {
			Telephone              string      `json:"telephone"`
			ConnectRate            string      `json:"connectRate"`
			SameTelephoneCrops     string      `json:"sameTelephoneCrops"`
			SameTelephoneCropsOcid []string    `json:"sameTelephoneCropsOcid"`
			SpecialTags            interface{} `json:"specialTags"`
			TelephoneSource        string      `json:"telephoneSource"`
		} `json:"telephoneDetailLists"`
		EntBrief         interface{} `json:"entBrief"`
		EmailList        interface{} `json:"emailList"`
		OfficialSiteList interface{} `json:"officialSiteList"`
		IsMobilephone    interface{} `json:"isMobilephone"`
		RegCity          string      `json:"regCity"`
		RegProvince      string      `json:"regProvince"`
		RegArea          string      `json:"regArea"`
		OpScope          string      `json:"opScope"`
		IsStrictDingOrg  bool        `json:"isStrictDingOrg"`
		CompanyPhoneInfo interface{} `json:"companyPhoneInfo"`
	} `json:"data"`
}

type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    Data        `json:"data"`
}
