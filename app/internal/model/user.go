package model

type User struct {
	Id             int32             `json:"id"`
	Email          string            `json:"email"`
	Name           string            `json:"name"`
	Surname        string            `json:"surname"`
	EmergencyGroup map[string]string `json:"emergencyGroup"`
	Password       string
}
