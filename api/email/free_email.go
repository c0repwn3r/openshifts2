package email

import (
	_ "embed"
	"slices"
	"strings"
)

//go:embed domains.txt
var free_email_providers string

func IsFreeEmailProvider(domain string) bool {
	split := strings.Split(free_email_providers, "\n")
	return slices.Contains(split, domain)
}
