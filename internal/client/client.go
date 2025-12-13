package client

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/mrbelka12000/interview_parser/internal/config"
)

type Client struct {
	cl  openai.Client
	cfg config.Config
}

func New(cfg *config.Config, apiKey string) *Client {
	return &Client{
		cl: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
		cfg: *cfg,
	}
}
