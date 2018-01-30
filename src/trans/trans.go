package trans

import (
	"context"
	"fmt"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
	"github.com/utahta/trans"
)

func Run() {
	plugin.Main(func(p *plugin.Plugin) error {
		c := &Command{Nvim: p.Nvim}

		p.HandleCommand(&plugin.CommandOptions{Name: "Trans", NArgs: "?", Range: "%"}, c.trans)
		p.HandleCommand(&plugin.CommandOptions{Name: "TransWord", NArgs: "?", Range: "%"}, c.transWord)
		p.HandleCommand(&plugin.CommandOptions{Name: "TransYank", NArgs: "?", Range: "%"}, c.transYank)
		return nil
	})
}

type Command struct {
	Nvim *nvim.Nvim
}

func (c *Command) trans(args []string) error {
	var sline []int
	if err := c.Nvim.Eval("getpos(\"'<\")", &sline); err != nil {
		return err
	}
	var eline []int
	if err := c.Nvim.Eval("getpos(\"'>\")", &eline); err != nil {
		return err
	}

	var text string
	if sline[1] == eline[1] {
		b, err := c.Nvim.CurrentLine()
		if err != nil {
			return err
		}
		text = string(b)
		text = text[sline[2]-1 : eline[2]-1]
	} else {
		b, err := c.Nvim.CurrentBuffer()
		if err != nil {
			return err
		}
		bytes, err := c.Nvim.BufferLines(b, sline[1]-1, eline[1], true)
		if err != nil {
			return err
		}

		tmp := make([]string, len(bytes))
		for i, b := range bytes {
			tmp[i] = string(b)
			if i == 0 {
				tmp[i] = tmp[i][sline[2]-1:]
			} else if i == len(tmp)-1 {
				tmp[i] = tmp[i][:eline[2]-1]
			}

			tmp[i] = strings.TrimSpace(tmp[i])
			tmp[i] = strings.TrimLeft(tmp[i], "//")
			tmp[i] = strings.TrimLeft(tmp[i], "#")
			tmp[i] = strings.TrimSpace(tmp[i])
		}
		text = strings.Join(tmp, " ")
	}

	to := c.langLocale()
	if len(args) > 0 {
		to = args[0]
	}
	return c.transOutput(text, "", to)
}

func (c *Command) transWord(args []string) error {
	text, err := c.getText("expand('<cword>')")
	if err != nil {
		return err
	}

	to := c.langLocale()
	if len(args) > 0 {
		to = args[0]
	}
	return c.transOutput(text, "", to)
}

func (c *Command) transYank(args []string) error {
	text, err := c.getText("@")
	if err != nil {
		return err
	}

	to := c.langLocale()
	if len(args) > 0 {
		to = args[0]
	}
	return c.transOutput(text, "", to)
}

func (c *Command) getText(expr string) (string, error) {
	var text string
	if err := c.Nvim.Eval(expr, &text); err != nil {
		return "", err
	}
	return strings.Replace(text, "\n", " ", -1), nil
}

func (c *Command) langLocale() string {
	var lang string
	if err := c.Nvim.Var("trans_lang_locale", &lang); err != nil {
		return "ja"
	}
	return lang
}

func (c *Command) transOutput(text string, source string, target string) error {
	ctx := context.Background()
	cli, err := trans.New(ctx)
	if err != nil {
		return err
	}

	output, err := cli.Translate(ctx, text, source, target)
	if err != nil {
		return err
	}

	c.Nvim.WriteOut(fmt.Sprintf("%s\n", output))
	return nil
}
