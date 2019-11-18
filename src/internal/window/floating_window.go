package window

import (
	"fmt"

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
		return fw.Close()
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
		width    int
		height   int
		winwidth int
	)

	if err := fw.vim.Call("winwidth", &winwidth, 0); err != nil {
		return err
	}
	winwidth -= 4 // for padding

	for i := 0; i < len(ss); i++ {
		var w int
		if err := fw.vim.Call("strdisplaywidth", &w, ss[i]); err != nil {
			return err
		}

		if w > winwidth {
			width = winwidth

			// Use binary search to find the maximum width of floating window.
			// But there may be a better way
			var (
				rs = []rune(ss[i])
				lb = -1
				ub = len(rs)
			)
			for ub-lb > 1 {
				mid := (ub + lb) / 2
				var midwidth int
				if err := fw.vim.Call("strdisplaywidth", &midwidth, string(rs[0:mid])); err != nil {
					return err
				}
				if midwidth >= winwidth {
					ub = mid
				} else {
					lb = mid
				}
			}

			ss[i] = string(rs[0:ub])
			if i == len(ss)-1 {
				ss = append(ss, string(rs[ub:]))
			} else {
				ss[i+1] = string(rs[ub:]) + ss[i+1]
			}
		} else if w > width {
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
