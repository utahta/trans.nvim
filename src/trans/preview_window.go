package trans

import "github.com/neovim/go-client/nvim"

type (
	previewWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer *buffer
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
	pw.buffer, err = newBuffer(pw.vim, withBufnr(bufnr))
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
