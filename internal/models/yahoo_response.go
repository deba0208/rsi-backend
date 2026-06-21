package models

type YahooResponse struct {
	Chart struct {
		Result []struct {
			Timestamp []int64 `json:"timestamp"`

			Indicators struct {
				Quote []struct {
					Close []*float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}
