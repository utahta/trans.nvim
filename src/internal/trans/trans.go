package trans

import (
	"context"
	"time"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
	"trans.nvim/src/internal/config"
	"trans.nvim/src/internal/event"
	"trans.nvim/src/internal/translator"
	"trans.nvim/src/internal/window"
)

// Run runs a trans.nvim plugin.
func Run() {
	plugin.Main(func(p *plugin.Plugin) error {
		h := &handler{
			vim:           p.Nvim,
			windowHandler: window.NewHandler(p.Nvim),
			translator:    translator.New(p.Nvim),
			config:        config.New(p.Nvim),
		}
		p.HandleCommand(&plugin.CommandOptions{Name: "Trans", NArgs: "?", Range: "%"}, h.Trans)
		p.HandleCommand(&plugin.CommandOptions{Name: "TransWord", NArgs: "?", Range: "%"}, h.TransWord)

		event.Init(p.Nvim)
		p.HandleAutocmd(&plugin.AutocmdOptions{
			Group:   event.Group,
			Event:   "CursorMoved,CursorMovedI",
			Pattern: "*",
		}, event.HandleFunc(event.TypeCursorMoved))
		return nil
	})
}

type handler struct {
	vim           *nvim.Nvim
	windowHandler window.Handler
	translator    translator.Translator
	config        config.Config
}

func (c *handler) Trans(args []string) error {
	target := c.config.Locale()
	if len(args) > 0 {
		target = args[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ss, err := c.translator.TranslatePos(ctx, translator.Option{
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
	return w.SetLine(ss)
}

func (c *handler) TransWord(args []string) error {
	target := c.config.Locale()
	if len(args) > 0 {
		target = args[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ss, err := c.translator.TranslateWord(ctx, translator.Option{
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
	return w.SetLine(ss)
}
