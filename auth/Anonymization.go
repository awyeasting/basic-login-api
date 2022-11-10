package auth

import "strings"

// Removes most indentifying information from an email
// in order to hint at what the email is to real users
// and to hide it from malicious users.
func AnonymizeEmail(email string) string {
	if len(email) < 5 {
		return email
	}

	firstPart := ""
	secondPart := ""
	thirdPart := ""

	parts := strings.Split(email, "@")
	firstPart = string(parts[0][0]) + strings.Repeat("*", len(parts[0])-1)

	parts = strings.Split(parts[1], ".")
	secondPart = string(parts[0][0]) + strings.Repeat("*", len(parts[0])-1)
	thirdPart = parts[1]

	return firstPart + "@" + secondPart + "." + thirdPart
}
