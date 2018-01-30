package trans

import (
	"context"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// Client represents google translate api client
type Client struct {
	*translate.Client
}

// New returns Client
func New(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	c, err := translate.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{Client: c}, nil
}

// Translate translates input text
func (c *Client) Translate(ctx context.Context, input string, s string, t string) (string, error) {
	var (
		source language.Tag
		target language.Tag
		err    error
	)
	if s != "" {
		source, err = language.Parse(s)
		if err != nil {
			return "", err
		}
	}
	target, err = language.Parse(t)
	if err != nil {
		return "", err
	}

	res, err := c.Client.Translate(ctx, []string{input}, target, &translate.Options{Source: source, Format: translate.Text})
	if err != nil {
		return "", err
	}

	return res[0].Text, nil
}
