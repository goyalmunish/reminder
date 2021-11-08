package models

import (
	"fmt"
)

type User struct {
	Name    string `json:"name"`
	EmailId string `json:"email_id"`
}

func (u User) String() string {
	return fmt.Sprintf("{Name: %v, EmailId: %v}", u.Name, u.EmailId)
}
