package window

import (
	"fmt"

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
	return mw.vim.WriteOut(fmt.Sprintf("%s\n", s))
}
