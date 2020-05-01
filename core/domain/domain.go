package domain

type BrazilCases struct{
	Country      string `json:"country"`
	NumberCases  int `json:"cases"`
	Deaths	 	 int `json:"deaths"`
	TodayCases 	 int `json:"todayCases"`
}
