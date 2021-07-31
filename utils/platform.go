package utils

type PlatformRequest struct {
	Project    string   `mapstructure:"PROJECT" json:"project"`
	Auth_Token string   `mapstructure:"Auth_Token" json:"project"`
	URL        string   `mapstructure:"PROJECT" json:"project"`
	TFAURL     string   `mapstructure:"PROJECT" json:"project"`
	Header     []string `mapstructure:"PROJECT" json:"project"`
}

type PlatformResponse struct {
	Id       string `json:"id"`
	Project  string `json:"project"`
	Messages string `json:"messages"`
}

type Platform struct {
	PlatformURL string `mapstructure:"PLATFORMURL" json:"platformURL"`
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
