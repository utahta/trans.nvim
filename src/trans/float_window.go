package trans

import "github.com/neovim/go-client/nvim"

type (
	floatWindow struct {
		vim    *nvim.Nvim
		id     nvim.Window
		buffer *buffer
	}
)

func (fw *floatWindow) Open() error {
	return nil
}

func (fw *floatWindow) Close() error {
	return nil
}

func (fw *floatWindow) SetLine(s string) error {
	return fw.buffer.Write(s)
}
