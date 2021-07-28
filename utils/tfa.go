package utils

type TFARequest struct {
	Data Data `json:"data"`
}

type Data struct {
	Id      string `json:"id"`
	Project string `json:"project"`
}

type TFAResponse struct {
	Prediction string `json:"prediction"`
}

type TFA struct {
	TFAURL string `mapstructure:"TFAURL"`
}

func (t *TFA) buildHeader() {
	//fill the header field
}

func (t *TFA) post() (err error) {
	// post
	return nil
}
func (t *TFA) get() (resp *TFAResponse) {
	//get with the full p.URL + full query
	return &TFAResponse{}
}
