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
	// this is a configuration of initial floating window.
	cfg := map[string]interface{}{
		"relative": "cursor",
		"col":      0,
		"row":      0,
		"width":    1,
		"height":   1,
	}
	if err := fw.vim.Call("nvim_open_win", nil, bufnr, true, cfg); err != nil {
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

	cfg := map[string]interface{}{
		"relative": "cursor",
		"col":      1,
		"row":      -height,
		"width":    width,
		"height":   height,
	}
	if err := fw.vim.Call("nvim_win_set_config", nil, fw.id, cfg); err != nil {
		return err
	}

	return fw.buffer.WriteStrings(ss)
}
