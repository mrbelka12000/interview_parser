package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/openai/openai-go"
)

var (
	promptAnalyze = `
Ты анализируешь транскрипт собеседования (интервьюер ↔ кандидат).

ЦЕЛЬ: извлечь ВСЕ реальные смысловые вопросы интервьюера и классифицировать
каждый как отвеченный или НЕотвеченный кандидатом.

ФОРМАТ ВЫВОДА — строго ТОЛЬКО valid JSON согласно структуре в конце промпта.

----------------------------------------------------------------------
                  ЖЁСТКИЕ ПРАВИЛА (не нарушать)
----------------------------------------------------------------------

1. Вопросом считается ЛЮБОЕ предложение интервьюера со знаком вопроса "?"
   которое связано с предметом технического обсуждения.

2. ВОПРОСЫ, КОТОРЫЕ НУЖНО ПОЛНОСТЬЮ ИГНОРИРОВАТЬ (НЕ добавлять никуда):

   - Проверка связи/звука:
     «Слышно?», «Нормально слышно?», «Алло?»
   - Организационные вопросы:
     «Все понятно?», «Понятно?», «Продолжаем?»,
     «Дальше поехали?», «Идем дальше?»
   - Вопросы о знании:
     «Ты это знал(а)?», «Знала?», «Это простое?»
   - Вопросы “есть ли у тебя вопросы?”:
     «Какие-то вопросы есть?», «Есть вопросы?»
   - Вопросы-комментарии:
     «Ok?», «Так?» если они не содержат смысловой нагрузки.
   - Риторические или эмоциональные:
     «Да?», «Верно?», «Понятно?»
   - Уточнения без нового смысла:
     «То есть как?», «А если так?» — ЕСЛИ они являются частью предыдущего вопроса.

3. Если вопрос НЕ подпадает под исключения — он ОБЯЗАН быть классифицирован.

4. Вопрос считается ОТВЕЧЁННЫМ ТОЛЬКО если:

   - ответ дал ИМЕННО КАНДИДАТ (answerer="candidate"), И
   - ответ относится по смыслу к вопросу, И
   - accuracy >= 0.7.

5. ЕСЛИ:
   - ответ дал интервьюер, ИЛИ
   - кандидат ушёл в сторону, ИЛИ
   - сказал «не знаю», «не помню», «затрудняюсь», ИЛИ
   - accuracy < 0.7, ИЛИ
   - модель не уверена, был ли ответ
   → то вопрос обязан попасть в questions_unanswered.

6. В массив answered Могут попадать ТОЛЬКО вопросы,
   где answerer="candidate".

----------------------------------------------------------------------
                   ПРАВИЛА ОЦЕНКИ ACCURACY
----------------------------------------------------------------------

- 1.0 → полный, точный ответ, по сути, без ошибок
- 0.7 –0.9 → хороший, но неполный ответ
- 0.7 → минимальный порог, чтобы считать вопрос отвечённым
- 0.3–0.7 → ответ частичный, с ошибками → unanswered
- 0.0–0.2 → «не знаю», уход в сторону, неправильный ответ → unanswered

----------------------------------------------------------------------
                     СТРУКТУРА JSON (НЕ МЕНЯТЬ)
----------------------------------------------------------------------

{
  "answered": [
    {
      "question": "string",
      "answer": "краткое содержание ответа кандидата",
      "full_answer": "полный текст ответа кандидата (обрезать до 300 символов)",
      "accuracy": 0.0,
      "questioner": "interviewer",
      "answerer": "candidate"
    }
  ],
  "unanswered": [
    {
      "question": "string",
      "questioner": "interviewer",
      "full_answer": "что ответил кандидат",
      "reason": "string (например: 'no clear answer from candidate', 'accuracy < 0.6', 'candidate avoided question')"
    }
  ]
}

`

	transcriptHeader = `
----------------------------------------------------------------------
                     ТРАНСКРИПТ ДЛЯ АНАЛИЗА
----------------------------------------------------------------------

	%v
`
)

type AnalyzeResponse struct {
	QuestionsAnswered []struct {
		Question      string  `json:"question"`
		AnswerSummary string  `json:"answer"`
		FullAnswer    string  `json:"full_answer"`
		Accuracy      float64 `json:"accuracy"`
		Questioner    string  `json:"questioner"`
		Answerer      string  `json:"answerer"`
	} `json:"answered"`
	QuestionsUnanswered []struct {
		Question         string `json:"question"`
		Questioner       string `json:"questioner"`
		FullAnswer       string `json:"full_answer"`
		ReasonUnanswered string `json:"reason"`
	} `json:"unanswered"`
}

func (c *Client) AnalyzeTranscript(ctx context.Context, text string) (out AnalyzeResponse, err error) {
	now := time.Now()
	log.Printf("[i] Analyzing transcript")

	res, err := c.cl.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(promptAnalyze),
			openai.UserMessage(fmt.Sprintf(transcriptHeader, text)),
		},
		Model: c.cfg.GPTClassifyQuestionsModel,
	})
	if err != nil {
		return out, fmt.Errorf("failed to analyze transcript: %w", err)
	}

	if err = json.Unmarshal([]byte(res.Choices[0].Message.Content), &out); err != nil {
		return out, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.Printf("[i] Finished analyzing transcript, seconds spent: %v", time.Since(now).Seconds())
	return out, nil
}
