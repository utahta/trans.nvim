package window

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type (
	Handler interface {
		OpenCurrentWindow(winType string) (Window, error)
		CloseCurrentWindow() error
	}

	handler struct {
		vim        *nvim.Nvim
		currentWin Window
	}

	Window interface {
		Open() error
		Close() error
		SetLine(s string) error
	}

	messageWindow struct {
		vim *nvim.Nvim
	}
)

func NewHandler(vim *nvim.Nvim) Handler {
	return &handler{vim: vim}
}

func (w *handler) OpenCurrentWindow(winType string) (Window, error) {
	switch winType {
	case "preview":
		w.currentWin = &previewWindow{vim: w.vim}
	case "float", "floating":
		w.currentWin = &floatingWindow{vim: w.vim}
	default:
		w.currentWin = &messageWindow{vim: w.vim}
	}

	if err := w.currentWin.Open(); err != nil {
		return nil, err
	}
	return w.currentWin, nil
}

func (w *handler) CloseCurrentWindow() error {
	if w.currentWin == nil {
		return nil
	}
	return w.currentWin.Close()
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
