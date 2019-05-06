package trans

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type (
	floatWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer *buffer
	}
)

func (fw *floatWindow) Open() error {
	bufnr, err := fw.vim.CurrentBuffer()
	if err != nil {
		return err
	}
	opts := map[string]interface{}{
		"relative": "cursor",
		"col":      1,
		"row":      1,
		"width":    100,
		"height":   2,
	}
	if err := fw.vim.Call("nvim_open_win", nil, bufnr, true, opts); err != nil {
		return err
	}
	fw.id, err = fw.vim.CurrentWindow()
	if err != nil {
		return err
	}

	fw.buffer, err = newBuffer(fw.vim)
	if err != nil {
		return err
	}

	if err := fw.vim.Command("wincmd p"); err != nil {
		return err
	}
	return nil
}

func (fw *floatWindow) Close() error {
	var winnr int
	if err := fw.vim.Call("win_id2win", &winnr, fw.id); err != nil {
		return err
	}

	if winnr > 0 {
		if err := fw.vim.Command(fmt.Sprintf("%dwincmd c", winnr)); err != nil {
			return err
		}
	}
	return nil
}

func (fw *floatWindow) SetLine(s string) error {
	return fw.buffer.Write(s)
}
