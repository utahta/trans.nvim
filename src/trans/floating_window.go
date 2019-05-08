package trans

import (
	"fmt"
	"math"
	"strings"

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
	if err := fw.vim.Call("nvim_open_win", nil, bufnr, true, fw.getWindowConfig(0, 0, 1, 80)); err != nil {
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
	// we considered to use nvim_win_get_number api, but it occurs an error when window id is invalid
	// so we still use win_id2win api.
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
		width  int
		height int
	)
	s = strings.TrimSpace(s)
	if err := fw.vim.Call("strdisplaywidth", &width, s); err != nil {
		return err
	}

	const maxWidth = 80
	if width > maxWidth {
		height = int(math.Ceil(float64(width) / float64(maxWidth)))
		width = maxWidth
	} else {
		height = 1
	}

	if err := fw.vim.Call("nvim_win_set_config", nil, fw.id, fw.getWindowConfig(-height, 0, height, width)); err != nil {
		return err
	}
	return fw.buffer.WriteString(s)
}

func (fw *floatingWindow) getWindowConfig(row, col, height, width int) map[string]interface{} {
	return map[string]interface{}{
		"relative":  "cursor",
		"row":       row,
		"col":       col,
		"height":    height,
		"width":     width,
		"focusable": true,
	}
}
