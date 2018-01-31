package trans

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/neovim/go-client/nvim/plugin"
	"github.com/utahta/trans"
	"google.golang.org/api/option"
)

func Run() {
	plugin.Main(func(p *plugin.Plugin) error {
		c := &Command{Nvim: p.Nvim}

		p.HandleCommand(&plugin.CommandOptions{Name: "Trans", NArgs: "?", Range: "%"}, c.Trans)
		p.HandleCommand(&plugin.CommandOptions{Name: "TransWord", NArgs: "?", Range: "%"}, c.TransWord)
		return nil
	})
}

type Command struct {
	Nvim *nvim.Nvim
}

func (c *Command) Trans(args []string) error {
	text, err := c.getText()
	if err != nil {
		return err
	}

	to := c.langLocale()
	if len(args) > 0 {
		to = args[0]
	}
	return c.translateOutput(text, "", to)
}

func (c *Command) TransWord(args []string) error {
	var text string
	if err := c.Nvim.Eval("expand('<cword>')", &text); err != nil {
		return err
	}

	to := c.langLocale()
	if len(args) > 0 {
		to = args[0]
	}
	return c.translateOutput(text, "", to)
}

func (c *Command) getText() (string, error) {
	var startpos []int
	if err := c.Nvim.Eval("getpos(\"'<\")", &startpos); err != nil {
		return "", err
	}
	var endpos []int
	if err := c.Nvim.Eval("getpos(\"'>\")", &endpos); err != nil {
		return "", err
	}

	if startpos[1] == 0 && startpos[2] == 0 &&
		endpos[1] == 0 && endpos[2] == 0 {
		return "", nil
	}

	var text string
	if startpos[1] == endpos[1] {
		b, err := c.Nvim.CurrentLine()
		if err != nil {
			return "", err
		}
		text = string(b)
		if endpos[2] > len(text) {
			endpos[2] = len(text)
		}
		text = text[startpos[2]-1 : endpos[2]]
	} else {
		b, err := c.Nvim.CurrentBuffer()
		if err != nil {
			return "", err
		}
		bytes, err := c.Nvim.BufferLines(b, startpos[1]-1, endpos[1], true)
		if err != nil {
			return "", err
		}

		tmp := make([]string, len(bytes))
		for i, b := range bytes {
			tmp[i] = string(b)
			if i == 0 {
				tmp[i] = tmp[i][startpos[2]-1:]
			} else if i == len(tmp)-1 {
				if endpos[2] > len(tmp[i]) {
					endpos[2] = len(tmp[i])
				}
				tmp[i] = tmp[i][:endpos[2]]
			}

			tmp[i] = strings.TrimSpace(tmp[i])
			for _, cutset := range c.langCutset() {
				tmp[i] = strings.TrimLeft(tmp[i], cutset)
			}
			tmp[i] = strings.TrimSpace(tmp[i])
		}
		text = strings.Join(tmp, " ")
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

func (c *Command) langCutset() []string {
	var cutset string
	if err := c.Nvim.Var("trans_lang_cutset", &cutset); err != nil {
		return []string{"//", "#"}
	}
	return strings.Split(cutset, " ")
}

func (c *Command) langCredentialsFile() string {
	var creds string
	if err := c.Nvim.Var("trans_lang_credentials_file", &creds); err != nil {
		return ""
	}
	if strings.HasPrefix(creds, "~") {
		creds = strings.TrimLeft(creds, "~")
		creds = fmt.Sprintf("%s/%s", os.Getenv("HOME"), creds)
	}
	return creds
}

func (c *Command) translateOutput(text string, source string, target string) error {
	if text == "" {
		return nil
	}

	var opts []option.ClientOption
	creds := c.langCredentialsFile()
	if creds != "" {
		opts = append(opts, option.WithCredentialsFile(creds))
	}

	ctx := context.Background()
	cli, err := trans.New(ctx, opts...)
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
