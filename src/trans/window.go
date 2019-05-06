package trans

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type (
	windowHandler struct {
		vim        *nvim.Nvim
		currentWin window
	}

	window interface {
		Open() error
		Close() error
		SetLine(s string) error
	}

	messageWindow struct {
		vim *nvim.Nvim
	}
)

func (w *windowHandler) Open(winType string) (window, error) {
	if w.currentWin != nil {
		if err := w.currentWin.Close(); err != nil {
			return nil, err
		}
	}

	switch winType {
	case "preview":
		w.currentWin = &previewWindow{vim: w.vim}
	case "float", "floating":
		w.currentWin = &floatWindow{vim: w.vim}
	default:
		w.currentWin = &messageWindow{vim: w.vim}
	}

	if err := w.currentWin.Open(); err != nil {
		return nil, err
	}
	return w.currentWin, nil
}

func (mw *messageWindow) Open() error {
	return nil
}

func (mw *messageWindow) Close() error {
	return nil
}

func (mw *messageWindow) SetLine(s string) error {
	return mw.vim.WriteOut(fmt.Sprintf("%s\n", s))
}
