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

	previewWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer *buffer
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
	default:
		w.currentWin = &messageWindow{vim: w.vim}
	}

	if err := w.currentWin.Open(); err != nil {
		return nil, err
	}
	return w.currentWin, nil
}

func (pw *previewWindow) Open() error {
	if err := pw.Close(); err != nil {
		return err
	}

	if err := pw.vim.Command("silent pedit translated"); err != nil {
		return err
	}
	if err := pw.vim.Command("wincmd P"); err != nil {
		return err
	}

	var err error
	pw.id, err = pw.vim.CurrentWindow()
	if err != nil {
		return err
	}
	if err := pw.vim.SetWindowHeight(pw.id, 5); err != nil {
		return err
	}

	bufID, err := pw.vim.CurrentBuffer()
	if err != nil {
		return err
	}
	pw.buffer, err = newBuffer(pw.vim, withBufferID(bufID))
	if err != nil {
		return err
	}

	if err := pw.vim.Command("wincmd p"); err != nil {
		return err
	}
	return nil
}

func (pw *previewWindow) Close() error {
	if err := pw.vim.Command("silent pclose"); err != nil {
		return err
	}
	return nil
}

func (pw *previewWindow) SetLine(s string) error {
	return pw.buffer.Write(s)
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
