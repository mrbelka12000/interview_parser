package parser

import (
	"github.com/mrbelka12000/interview_parser/internal/config"
)

type (
	Parser struct {
		Cfg *config.Config
	}
)

func NewParser(cfg *config.Config) *Parser {
	return &Parser{
		Cfg: cfg,
	}
}
