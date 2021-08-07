package main

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"scope"`
}

type CB_TimeStampData struct {
	Data CB_TimeStamp
}

type CB_TimeStamp struct {
	ISO   string
	Epoch int64
}

type CB_Balance struct {
	Amount   string
	Currency string
}

type CB_CurrencyInfo struct {
	Code string `json:"code"`
}

type CB_Data struct {
	Name     string          `json:"name"`
	Currency CB_CurrencyInfo `json:"currency"`
	Balance  CB_Balance      `json:"balance"`
}

type CB_Listed_Crypto struct {
	Currencies []CB_Data `json:"data"`
}

type CB_Single_Crypto struct {
	Data CB_Data `json:"data"`
}

type CB_Spot_Price struct {
	Data CB_Balance `json:"Data"`
}
