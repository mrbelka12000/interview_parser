package internal

import (
	"fmt"
	"os"

	"github.com/mrbelka12000/interview_parser/internal/client"
)

func SaveTranscript(outputFile, transcript string) error {
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

func SaveAnalyzeResponse(outputFile string, response client.AnalyzeResponse) error {
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

	for _, v := range response.QuestionsAnswered {
		_, err = f.WriteString(fmt.Sprintf(`
###  %v
%v
%v
`, v.Question, v.AnswerSummary, v.Accuracy))
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	for _, v := range response.QuestionsUnanswered {
		_, err = f.WriteString(fmt.Sprintf(`
###  %v
%v
%v
`, v.Question, v.FullAnswer, v.ReasonUnanswered))
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}
