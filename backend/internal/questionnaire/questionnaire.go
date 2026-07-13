package questionnaire

import (
	_ "embed"
	"strings"
)

//go:embed client_onboarding.md
var clientOnboarding string

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

func (Provider) Load() string {
	return strings.TrimSpace(clientOnboarding)
}
