package trans

import (
	"context"
	"time"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
)

// Run runs a trans.nvim plugin.
func Run() {
	plugin.Main(func(p *plugin.Plugin) error {
		h := &handler{
			vim:           p.Nvim,
			windowHandler: windowHandler{vim: p.Nvim},
			translator:    translator{vim: p.Nvim},
			config:        config{vim: p.Nvim},
		}
		p.HandleCommand(&plugin.CommandOptions{Name: "Trans", NArgs: "?", Range: "%"}, h.Trans)
		p.HandleCommand(&plugin.CommandOptions{Name: "TransWord", NArgs: "?", Range: "%"}, h.TransWord)
		return nil
	})
}

type handler struct {
	vim           *nvim.Nvim
	windowHandler windowHandler
	translator    translator
	config        config
}

func (c *handler) Trans(args []string) error {
	target := c.config.Locale()
	if len(args) > 0 {
		target = args[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	text, err := c.translator.TranslatePos(ctx, translatorOption{
		Source:          "",
		Target:          target,
		Cutset:          c.config.Cutset(),
		CredentialsFile: c.config.CredentialsFile(),
	})
	if err != nil {
		return err
	}

	if err := c.windowHandler.CloseCurrentWindow(); err != nil {
		return err
	}
	w, err := c.windowHandler.OpenCurrentWindow(c.config.Output())
	if err != nil {
		return err
	}
	return w.SetLine(text)
}

func (c *handler) TransWord(args []string) error {
	target := c.config.Locale()
	if len(args) > 0 {
		target = args[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	text, err := c.translator.TranslateWord(ctx, translatorOption{
		Source:          "",
		Target:          target,
		Cutset:          c.config.Cutset(),
		CredentialsFile: c.config.CredentialsFile(),
	})
	if err != nil {
		return err
	}

	if err := c.windowHandler.CloseCurrentWindow(); err != nil {
		return err
	}
	w, err := c.windowHandler.OpenCurrentWindow(c.config.Output())
	if err != nil {
		return err
	}
	return w.SetLine(text)
}
