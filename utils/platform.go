package utils

type PlatformRequest struct {
	Project    string
	Auth_Token string
	URL        string
	TFAURL     string
	Header     []string
}

type PlatformResponse struct {
	Id       string `json:"id"`
	Project  string `json:"project"`
	Messages string `json:"messages"`
}

type Platform struct {
	PlatformURL string `mapstructure:"PLATFORMURL"`
}

//fill the header field
func (p *Platform) buildHeader() {
}

// post
func (p *Platform) post() {
}

//get with the full p.URL + full query
func (p *Platform) get() (resp *PlatformResponse) {
	return &PlatformResponse{}
}

func (p *Platform) updateOriginal() (err error) {
	return nil
}
