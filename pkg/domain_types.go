package pkg

import "time"

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

type Routine struct {
	Id        int32     `json:"routine_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RoutineGist struct {
	RoutineId int32  `json:"routine_id"`
	GistId    string `json:"gist_id"`
	UserId    int32  `json:"user_id"`
}

type Session struct {
	SessionId int32     `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UserId    int32     `json:"user_id"`
}
