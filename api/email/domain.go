package email

import "strings"

func GetDomain(email string) string {
	v := strings.Split(email, "@")
	return v[1]
}
