package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
)

type (
	Config interface {
		Locale() string
		Cutset() []string
		CredentialsFile() string
		Output() string
	}

	config struct {
		vim *nvim.Nvim
	}
)

func New(vim *nvim.Nvim) Config {
	return &config{vim: vim}
}

func (c *config) Locale() string {
	var lang string
	if err := c.vim.Var("trans_lang_locale", &lang); err != nil {
		return "ja"
	}
	return lang
}

func (c *config) Cutset() []string {
	var cutset string
	if err := c.vim.Var("trans_lang_cutset", &cutset); err != nil {
		return []string{"//", "#"}
	}
	return strings.Split(cutset, " ")
}

func (c *config) CredentialsFile() string {
	var creds string
	if err := c.vim.Var("trans_lang_credentials_file", &creds); err != nil {
		return ""
	}
	if strings.HasPrefix(creds, "~") {
		creds = strings.TrimLeft(creds, "~")
		creds = fmt.Sprintf("%s/%s", os.Getenv("HOME"), creds)
	}
	return creds
}

// Output returns a `preview` or `` or `float`.
//
// `preview` displays text with preview window.
// `` (empty string) displays text with message.
// `float` displays text with floating window.
func (c *config) Output() string {
	var o string
	if err := c.vim.Var("trans_lang_output", &o); err != nil {
		return ""
	}
	return o
}
