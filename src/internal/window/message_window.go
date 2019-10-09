package window

import (
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
)

type (
	messageWindow struct {
		vim *nvim.Nvim
	}
)

func (mw *messageWindow) Open() error {
	return nil
}

func (mw *messageWindow) Close() error {
	return nil
}

func (mw *messageWindow) SetLine(s string) error {
	s = strings.Replace(s, "\n", " ", -1)
	return mw.vim.WriteOut(fmt.Sprintf("%s\n", s))
}
