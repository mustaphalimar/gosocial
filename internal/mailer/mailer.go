package mailer

import "embed"

const (
	FromName            = "GoSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

// the line below ensures the template files will be embedded with the go binary at build time!

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
