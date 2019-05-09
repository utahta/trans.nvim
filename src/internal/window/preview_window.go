package window

import (
	"time"

	"github.com/neovim/go-client/nvim"
	"trans.nvim/src/internal/buffer"
	"trans.nvim/src/internal/event"
)

type (
	previewWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer buffer.Buffer
	}
)

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

	bufnr, err := pw.vim.CurrentBuffer()
	if err != nil {
		return err
	}
	pw.buffer, err = buffer.New(pw.vim, buffer.WithBufnr(bufnr))
	if err != nil {
		return err
	}

	if err := pw.vim.Command("wincmd p"); err != nil {
		return err
	}

	event.On(event.TypeMoveEvent, func() error {
		timer := time.NewTimer(1500 * time.Millisecond)
		select {
		case <-timer.C:
			return pw.Close()
		}
	})
	return nil
}

func (pw *previewWindow) Close() error {
	if err := pw.vim.Command("silent pclose"); err != nil {
		return err
	}
	return nil
}

func (pw *previewWindow) SetLine(s string) error {
	return pw.buffer.WriteString(s)
}
