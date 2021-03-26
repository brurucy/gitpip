package pkg

type GistOwner struct {
	Login string `json:"login"`
	Id    int32  `json:"id"`
}

type GistFile struct {
	Filename string `json:"filename"`
	RawUrl   string `json:"raw_url"`
}

type Gist struct {
	GistUUID  string              `json:"id"`
	GistOwner GistOwner           `json:"owner"`
	GistFiles map[string]GistFile `json:"files,omitempty"`
}

type PipeAddActivityRequest struct {
	Note    string `json:"note"`
	Subject string `json:"subject"`
	Done    bool   `json:"done"`
}
