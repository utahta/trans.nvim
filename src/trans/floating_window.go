package trans

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type (
	floatingWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer *buffer
	}
)

func (fw *floatingWindow) Open() error {
	bufnr, err := fw.vim.CurrentBuffer()
	if err != nil {
		return err
	}
	if err := fw.vim.Call("nvim_open_win", nil, bufnr, true, fw.getConfig(1, 1)); err != nil {
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

func (fw *floatingWindow) Close() error {
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

func (fw *floatingWindow) SetLine(s string) error {
	var (
		ss    []string
		width int
	)
	if err := fw.vim.Call("strdisplaywidth", &width, s); err != nil {
		return err
	}

	const maxWidth = 80
	if width > maxWidth {
		for i := 0; i <= width/maxWidth; i++ {
			var (
				text   string
				start  = i * maxWidth
				length = i*maxWidth + maxWidth
			)
			if err := fw.vim.Call("strcharpart", &text, s, start, length); err != nil {
				return err
			}
			ss = append(ss, s)
		}
		width = maxWidth
	} else {
		ss = append(ss, s)
	}
	height := len(ss)

	if err := fw.vim.Call("nvim_win_set_config", nil, fw.id, fw.getConfig(width, height)); err != nil {
		return err
	}

	return fw.buffer.WriteStrings(ss)
}

func (fw *floatingWindow) getConfig(width, height int) map[string]interface{} {
	return map[string]interface{}{
		"relative":  "cursor",
		"col":       1,
		"row":       -height,
		"width":     width,
		"height":    height,
		"focusable": true,
	}
}
