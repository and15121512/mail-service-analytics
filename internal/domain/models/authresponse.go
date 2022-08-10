package models

const (
	OkAuthRespStatus = iota
	RefreshedAuthRespStatus
	RefusedAuthRespStatus
)

type AuthResult struct {
	Status       int
	AccessToken  string
	RefreshToken string
	Login        string
}
