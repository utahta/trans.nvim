package trans

import (
	"github.com/neovim/go-client/nvim"
)

type (
	buffer struct {
		vim    *nvim.Nvim
		number nvim.Buffer
	}

	bufferOption func(b *buffer)
)

func newBuffer(vim *nvim.Nvim, opts ...bufferOption) (*buffer, error) {
	b := &buffer{
		vim:    vim,
		number: -1,
	}
	for _, opt := range opts {
		opt(b)
	}

	if b.number < 0 {
		if err := b.New(); err != nil {
			return nil, err
		}
	}
	if err := b.applyOptions(); err != nil {
		return nil, err
	}
	return b, nil
}

func withBufnr(nr nvim.Buffer) bufferOption {
	return func(b *buffer) {
		b.number = nr
	}
}

func (b *buffer) New() error {
	if err := b.vim.Command("silent enew"); err != nil {
		return err
	}
	var err error
	b.number, err = b.vim.CurrentBuffer()
	if err != nil {
		return err
	}
	return nil
}

func (b *buffer) applyOptions() error {
	options := []struct {
		name  string
		value interface{}
	}{
		{"buftype", "nofile"},
		{"bufhidden", "wipe"},
		{"buflisted", false},
		{"swapfile", false},
		{"modeline", false},
	}
	for _, o := range options {
		if err := b.vim.SetBufferOption(b.number, o.name, o.value); err != nil {
			return err
		}
	}
	return nil
}

func (b *buffer) WriteString(s string) error {
	return b.WriteStrings([]string{s})
}

func (b *buffer) WriteStrings(ss []string) error {
	var bs [][]byte
	for _, s := range ss {
		bs = append(bs, []byte(s))
	}
	if err := b.vim.SetBufferLines(b.number, 0, 0, false, bs); err != nil {
		return err
	}
	options := []struct {
		name  string
		value interface{}
	}{
		{"modified", false},
		{"modifiable", false},
	}
	for _, o := range options {
		if err := b.vim.SetBufferOption(b.number, o.name, o.value); err != nil {
			return err
		}
	}
	return nil
}