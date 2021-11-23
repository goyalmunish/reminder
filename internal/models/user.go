package models

import (
	"fmt"
)

type User struct {
	Name    string `json:"name"`
	EmailId string `json:"email_id"`
}

// provide basic string representation of a user
func (u User) String() string {
	return fmt.Sprintf("{Name: %v, EmailId: %v}", u.Name, u.EmailId)
}
