package client

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
	promptAnalyzeEN = `
You are analyzing an interview transcript (interviewer ↔ candidate).

GOAL: extract ALL meaningful real questions asked by the interviewer and classify
each one as answered or UNanswered by the candidate.

OUTPUT FORMAT — strictly ONLY valid JSON according to the structure at the end of this prompt.

----------------------------------------------------------------------
                           HARD RULES (do not violate)
----------------------------------------------------------------------

1. A question is ANY sentence from the interviewer that contains a question mark "?"
   and is related to the subject of the technical discussion.

2. QUESTIONS THAT MUST BE COMPLETELY IGNORED (do NOT add them anywhere):

   - Connection/audio checks:
     “Do you hear me?”, “Is the sound okay?”, “Hello?”
   - Organizational questions:
     “Is everything clear?”, “Clear?”, “Shall we continue?”,
     “Move on?”, “Shall we go further?”
   - Knowledge-check questions:
     “Did you know this?”, “You knew?”, “Is this simple?”
   - “Do you have any questions?” type:
     “Any questions?”, “Do you have questions?”
   - Comment-questions:
     “Ok?”, “Like this?” — if they do not add semantic meaning.
   - Rhetorical or emotional questions:
     “Right?”, “Yes?”, “Clear?”
   - Clarifications without new meaning:
     “Meaning what?”, “What if this?” — IF they are part of the previous question.

3. If a question does NOT fall under exclusions — it MUST be classified.

4. A question is considered ANSWERED ONLY if:

   - the answer was given SPECIFICALLY by the CANDIDATE (answerer="candidate"), AND
   - the answer is semantically related to the question, AND
   - accuracy >= 0.7.

5. IF:
   - the interviewer gave the answer, OR
   - the candidate deviated from the topic, OR
   - said “don’t know”, “don’t remember”, “not sure”, OR
   - accuracy < 0.7, OR
   - the model is unsure whether the question was answered
   → then the question MUST be placed into questions_unanswered.

6. Only questions where answerer="candidate" may appear in the answered array.

----------------------------------------------------------------------
                     ACCURACY EVALUATION RULES
----------------------------------------------------------------------

- 1.0 → full, correct, precise answer, with no errors
- 0.7–0.9 → good but incomplete answer
- 0.7 → minimal threshold to consider a question answered
- 0.3–0.7 → partial or flawed answer → unanswered
- 0.0–0.2 → “don’t know”, irrelevant, or incorrect answer → unanswered

----------------------------------------------------------------------
                     JSON STRUCTURE (DO NOT CHANGE)
----------------------------------------------------------------------

{
  "answered": [
    {
      "question": "string",
      "answer": "short summary of the candidate's answer",
      "full_answer": "full text of the candidate's answer (trim to 300 characters)",
      "accuracy": 0.0,
      "questioner": "interviewer",
      "answerer": "candidate"
    }
  ],
  "unanswered": [
    {
      "question": "string",
      "questioner": "interviewer",
      "full_answer": "what the candidate said",
      "reason": "string (e.g. 'no clear answer from candidate', 'accuracy < 0.6', 'candidate avoided question')"
    }
  ]
}
`

	transcriptHeaderEN = `
----------------------------------------------------------------------
					TRANSCRIPT FOR ANALYSIS
----------------------------------------------------------------------

	%v
`
)

func getAnalyzePrompt(lang string) string {
	switch lang {
	case "en":
		return promptAnalyzeEN
	default:
		return promptAnalyze
	}
}

func getTranscriptPrompt(lang string) string {
	switch lang {
	case "en":
		return transcriptHeaderEN
	default:
		return transcriptHeader
	}
}
