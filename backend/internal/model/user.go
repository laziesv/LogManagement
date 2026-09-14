package model

type User struct {
	ID     string `json:"id"`
	Tenant string `json:"tenant"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Hash   string `json:"-"`
}
