package parser

import (
	"fmt"
	"os"
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
