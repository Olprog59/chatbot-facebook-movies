package models

import (
	"fmt"
	"os"
	"strings"
)

type User struct {
	FirstName string
	LastName  string
}

func GetAuthorizedUsers() []User {
	var users []User
	for i := 1; ; i++ {
		userEnv := os.Getenv(fmt.Sprintf("USER%d", i))
		if userEnv == "" {
			break
		}
		parts := strings.Split(userEnv, ";")
		if len(parts) == 2 {
			users = append(users, User{FirstName: parts[0], LastName: parts[1]})
		}
	}
	return users
}
