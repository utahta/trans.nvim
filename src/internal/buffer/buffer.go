package buffer

import (
	"github.com/neovim/go-client/nvim"
)

type (
	Buffer interface {
		WriteString(s string) error
		WriteStrings(ss []string) error
	}

	buffer struct {
		vim    *nvim.Nvim
		number nvim.Buffer
	}

	Option func(b *buffer)
)

func New(vim *nvim.Nvim, opts ...Option) (Buffer, error) {
	b := &buffer{
		vim:    vim,
		number: -1,
	}
	for _, opt := range opts {
		opt(b)
	}

	if b.number < 0 {
		if err := b.new(); err != nil {
			return nil, err
		}
	}
	if err := b.applyOptions(); err != nil {
		return nil, err
	}
	return b, nil
}

func WithBufnr(nr nvim.Buffer) Option {
	return func(b *buffer) {
		b.number = nr
	}
}

func (b *buffer) new() error {
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
