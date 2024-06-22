package src

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Repo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Details  string `json:"details"`
	Username string `json:"user_name"`
}

type VerifyCallback struct {
	State    string
	Verifier string
	Expires  time.Time
}
