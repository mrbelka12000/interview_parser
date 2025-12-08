package parser

import (
	"github.com/mrbelka12000/interview_parser/internal/client"
	"github.com/mrbelka12000/interview_parser/internal/config"
)

type (
	Parser struct {
		cfg      *config.Config
		aiClient *client.Client
	}
)

func NewParser(cfg *config.Config, aiClient *client.Client) *Parser {
	return &Parser{
		cfg:      cfg,
		aiClient: aiClient,
	}
}
