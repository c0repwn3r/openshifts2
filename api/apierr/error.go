package apierr

const (
	EUnauthorized    = "unauthorized"
	EEmailInUse      = "email_in_use"
	EFreeEmail       = "free_email"
	EDomainInUse     = "domain_in_use"
	EPasswordTooLong = "password_too_long"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}
