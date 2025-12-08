package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
)

func (c *Client) IsValidAPIKeysProvided() error {

	_, err := c.cl.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		Model: c.cfg.GPTTranscribeModel,
	})
	if err != nil && !strings.Contains(err.Error(), "Field required") {
		if strings.Contains(err.Error(), "does not have access to model") {
			return errors.New(fmt.Sprintf(`Your project does not have access to model: %s
You can add it here: https://platform.openai.com/settings/proj_qJmmMi9x0FgRhq9KgNLsq2xw/limits`, c.cfg.GPTTranscribeModel))
		}
		return errors.New("Invalid API Key")
	}

	_, err = c.cl.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello, it is test request!"),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		if strings.Contains(err.Error(), "does not have access to model") {
			return errors.New(fmt.Sprintf(`Your project does not have access to model: %s
You can add it here: https://platform.openai.com/settings/proj_qJmmMi9x0FgRhq9KgNLsq2xw/limits`, c.cfg.GPTClassifyQuestionsModel))
		}
		return errors.New("Invalid API Key")
	}

	return nil
}
