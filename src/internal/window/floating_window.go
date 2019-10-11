package window

import (
	"fmt"
	"time"

	"github.com/neovim/go-client/nvim"
	"trans.nvim/src/internal/buffer"
	"trans.nvim/src/internal/event"
)

type (
	floatingWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer buffer.Buffer
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

	fw.buffer, err = buffer.New(fw.vim)
	if err != nil {
		return err
	}

	if err := fw.vim.Command("wincmd p"); err != nil {
		return err
	}

	event.Once("CursorMoved,CursorMovedI", "<buffer>", func() error {
		select {
		case <-time.After(500 * time.Millisecond):
			return fw.Close()
		}
	})
	return nil
}

func (fw *floatingWindow) Close() error {
	// we have decided to use win_id2_win api because nvim_win_get_number api occurs an error when window id is invalid.
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

func (fw *floatingWindow) SetLine(ss []string) error {
	var (
		width  int
		height int
	)

	for i := range ss {
		var w int
		if err := fw.vim.Call("strdisplaywidth", &w, ss[i]); err != nil {
			return err
		}
		if w > width {
			width = w
		}
	}
	height = len(ss)

	// padding
	width += 4
	height += 2
	for i := range ss {
		ss[i] = "  " + ss[i]
	}
	ss = append([]string{""}, ss...)

	var (
		winline int
		row     = 1
		col     = 0
	)
	if err := fw.vim.Call("winline", &winline); err != nil {
		return err
	}
	if (winline - height) > 0 {
		row = -height
	}

	if err := fw.vim.Call("nvim_win_set_config", nil, fw.id, fw.getWindowConfig(row, col, height, width)); err != nil {
		return err
	}
	return fw.buffer.WriteStrings(ss)
}

func (fw *floatingWindow) getWindowConfig(row, col, height, width int) map[string]interface{} {
	return map[string]interface{}{
		"relative":  "cursor",
		"row":       row,
		"col":       col,
		"height":    height,
		"width":     width,
		"focusable": true,
		"style":     "minimal",
	}
}
