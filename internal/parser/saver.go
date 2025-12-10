package parser

import (
	"fmt"
	"os"

	"github.com/mrbelka12000/interview_parser/internal/client"
)

func (p *Parser) SaveTranscript(outputFile, transcript string) error {
	// save text into file
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(transcript))
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func (p *Parser) SaveAnalyzeResponse(outputFile string, response client.AnalyzeResponse) error {
	if err := os.Remove(outputFile); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("Error removing file %s\n", outputFile)
		}
	}

	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, q := range response.Questions {
		switch {
		case q.Accuracy > 0.7:
			_, err = f.WriteString(fmt.Sprintf(`
%v
%v
%v
`, q.Question, q.FullAnswer, q.Accuracy))
		default:
			_, err = f.WriteString(fmt.Sprintf(`
%v
%v
%v
%v
`, q.Question, q.FullAnswer, q.ReasonUnanswered, q.Accuracy))
		}
	}

	return nil
}

func (p *Parser) SaveCallAnalysis(outputFile, analysisText string) error {
	// save analysis text into file
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open output file: %v", err)
	}
	defer f.Close()

	_, err = f.Write([]byte(analysisText))
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}
