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
		SetLine([]string) error
	}
)

func NewHandler(vim *nvim.Nvim) Handler {
	return &handler{vim: vim}
}

func (h *handler) OpenCurrentWindow(winType string) (Window, error) {
	switch winType {
	case "preview":
		h.currentWin = &previewWindow{vim: h.vim}
	case "float":
		if h.canUseFloatingWindow() {
			h.currentWin = &floatingWindow{vim: h.vim}
		} else {
			h.currentWin = &previewWindow{vim: h.vim}
		}
	default:
		h.currentWin = &messageWindow{vim: h.vim}
	}

	if err := h.currentWin.Open(); err != nil {
		return nil, err
	}
	return h.currentWin, nil
}

func (h *handler) CloseCurrentWindow() error {
	if h.currentWin == nil {
		return nil
	}
	return h.currentWin.Close()
}

func (h *handler) canUseFloatingWindow() bool {
	for _, expr := range []string{`has('nvim')`, `exists('*nvim_win_set_config')`} {
		var valid bool
		if err := h.vim.Eval(expr, &valid); err != nil {
			h.vim.WriteErr(fmt.Sprintf("invalid: %v\n", err))
			return false
		}

		if !valid {
			return false
		}
	}
	return true
}
