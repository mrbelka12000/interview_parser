package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/mrbelka12000/interview_parser/internal/client"
	"github.com/mrbelka12000/interview_parser/internal/config"
)

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeStart      MessageType = "start"
	MessageTypeResponse   MessageType = "response"
	MessageTypeAudio      MessageType = "audio"
	MessageTypeTranscribe MessageType = "transcribe"
	MessageTypeError      MessageType = "error"
	MessageTypeEnd        MessageType = "end"
	MessageTypeAnalytics  MessageType = "analytics"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type          MessageType `json:"type"`
	Data          interface{} `json:"data"`
	Timestamp     time.Time   `json:"timestamp"`
	AnalyticsData interface{} `json:"analytics_data"`
}

// StartMessage represents interview start request
type StartMessage struct {
	CV             string `json:"cv"`
	VacancyInfo    string `json:"vacancy_info"`
	Specialization string `json:"specialization"`
	Level          string `json:"level"`
	Meta           string `json:"meta"`
	QuestionsCount int    `json:"questions_count"`
}

// ResponseMessage represents a response (user answer or AI question)
type ResponseMessage struct {
	Text       string    `json:"text"`
	Timestamp  time.Time `json:"timestamp"`
	IsFromAI   bool      `json:"is_from_ai"`
	IsQuestion bool      `json:"is_question,omitempty"`
}

// AudioMessage represents audio data for transcription
type AudioMessage struct {
	AudioData []byte `json:"audio_data"`
	Format    string `json:"format"`
}

// AnalyticsMessage represents analytics data
type AnalyticsMessage struct {
	QuestionIndex int     `json:"question_index"`
	Question      string  `json:"question"`
	Answer        string  `json:"answer"`
	Category      string  `json:"category"`
	Accuracy      float64 `json:"accuracy"`
	Feedback      string  `json:"feedback"`
}

// InterviewSession manages a single mock interview session
type InterviewSession struct {
	conn           *websocket.Conn
	aiClient       *client.Client
	isActive       bool
	currentIndex   int
	questions      []client.GeneratedQuestion
	answers        []string
	context        string
	cv             string
	vacancyInfo    string
	specialization string
	level          string
	meta           string
	startTime      time.Time
	mutex          sync.RWMutex
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	currentSession *InterviewSession
	sessionMutex   sync.Mutex
)

func (s *InterviewSession) handleConnection() {
	defer s.close()

	// Send welcome message
	welcomeMsg := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      "Welcome to your mock interview! I'm generating personalized questions based on your CV and the vacancy information.",
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		Timestamp: time.Now(),
	}
	s.sendMessage(welcomeMsg)

	for {
		var msg WSMessage
		err := s.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.Timestamp = time.Now()

		switch msg.Type {
		case MessageTypeStart:
			go s.handleStartMessage(msg)
		case MessageTypeResponse:
			go s.handleResponseMessage(msg)
		case MessageTypeAudio:
			go s.handleAudioMessage(msg)
		case MessageTypeTranscribe:
			go s.handleTranscribeMessage(msg)
		case MessageTypeEnd:
			go s.handleEndMessage()
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (s *InterviewSession) close() {
	if !s.isActive {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.isActive = false
	if s.conn != nil {
		s.conn.Close()
	}

	log.Printf("Mock interview session closed")
}

func (s *InterviewSession) handleStartMessage(msg WSMessage) {
	body, _ := json.Marshal(msg.Data)
	var startMsg StartMessage
	if err := json.Unmarshal(body, &startMsg); err != nil {
		s.sendError("Invalid start message format")
		return
	}

	s.mutex.Lock()
	s.cv = startMsg.CV
	s.vacancyInfo = startMsg.VacancyInfo
	s.specialization = startMsg.Specialization
	s.level = startMsg.Level
	s.meta = startMsg.Meta
	s.startTime = time.Now()
	s.mutex.Unlock()

	// Generate interview questions based on provided data
	go s.generateQuestions(startMsg)
}

func (s *InterviewSession) generateQuestions(startMsg StartMessage) {
	if s.aiClient == nil {
		s.sendError("AI client not initialized")
		return
	}

	// Send generating message
	generatingMsg := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      "🤔 Generating personalized interview questions based on your CV and the vacancy...",
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		Timestamp: time.Now(),
	}
	s.sendMessage(generatingMsg)

	// Create mock interview request
	req := client.MockInterviewRequest{
		CV:             startMsg.CV,
		VacancyInfo:    startMsg.VacancyInfo,
		Specialization: startMsg.Specialization,
		Level:          startMsg.Level,
		Meta:           startMsg.Meta,
		QuestionsCount: startMsg.QuestionsCount,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := s.aiClient.GetMockInterviewQuestions(ctx, req)
	if err != nil {
		s.sendError(fmt.Sprintf("Failed to generate questions: %v", err))
		return
	}

	s.mutex.Lock()
	s.questions = response.GeneratedQuestions
	s.answers = make([]string, len(response.GeneratedQuestions))
	s.currentIndex = 0
	s.mutex.Unlock()

	// Send vacancy summary
	summaryMsg := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      fmt.Sprintf("📋 Vacancy Summary: %s", response.VacancySummary),
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		Timestamp: time.Now(),
	}
	s.sendMessage(summaryMsg)

	// Brief delay then send first question
	time.Sleep(2 * time.Second)
	s.sendNextQuestion()
}

func (s *InterviewSession) sendNextQuestion() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.currentIndex >= len(s.questions) {
		// Interview completed
		s.endInterview()
		return
	}

	question := s.questions[s.currentIndex]
	response := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:       fmt.Sprintf("Question %d (%s): %s", s.currentIndex+1, question.Category, question.Question),
			Timestamp:  time.Now(),
			IsFromAI:   true,
			IsQuestion: true,
		},
		Timestamp: time.Now(),
	}

	s.sendMessage(response)
}

func (s *InterviewSession) handleResponseMessage(msg WSMessage) {
	body, _ := json.Marshal(msg.Data)

	var responseMsg ResponseMessage
	if err := json.Unmarshal(body, &responseMsg); err != nil {
		s.sendError("Invalid response message format")
		return
	}

	// Store the user's response
	s.mutex.Lock()
	if s.currentIndex < len(s.answers) {
		s.answers[s.currentIndex] = responseMsg.Text
	}
	currentIdx := s.currentIndex
	s.mutex.Unlock()

	log.Printf("User response for question %d: %s", currentIdx+1, responseMsg.Text)

	// Send user's response back to chat (so it appears on the right side)
	userResponse := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      responseMsg.Text,
			Timestamp: time.Now(),
			IsFromAI:  false,
		},
		Timestamp: time.Now(),
	}
	s.sendMessage(userResponse)

	// Move to next question after a brief delay
	time.Sleep(1 * time.Second)

	s.mutex.Lock()
	s.currentIndex++
	s.mutex.Unlock()

	// Send next question
	time.Sleep(1 * time.Second)
	s.sendNextQuestion()
}

func (s *InterviewSession) handleAudioMessage(msg WSMessage) {
	body, _ := json.Marshal(msg.Data)

	var audioMsg AudioMessage
	if err := json.Unmarshal(body, &audioMsg); err != nil {
		s.sendError("Invalid audio message format")
		return
	}

	// Transcribe audio
	go s.transcribeAudio(audioMsg.AudioData)
}

func (s *InterviewSession) handleTranscribeMessage(msg WSMessage) {
	body, _ := json.Marshal(msg.Data)
	var audioMsg AudioMessage
	if err := json.Unmarshal(body, &audioMsg); err != nil {
		s.sendError("Invalid transcribe message format")
		return
	}

	go s.transcribeAudio(audioMsg.AudioData)
}

func (s *InterviewSession) transcribeAudio(audioData []byte) {
	if s.aiClient == nil {
		s.sendError("AI client not initialized")
		return
	}

	// Send transcribing message
	transcribingMsg := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      "🎙️ Transcribing your audio response...",
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		Timestamp: time.Now(),
	}
	s.sendMessage(transcribingMsg)

	// In a real implementation, you'd use your existing transcription functionality
	// For now, simulate transcription
	time.Sleep(2 * time.Second)
	transcribedText := "[Audio transcription would appear here - this would be the actual transcribed text from the audio]"

	// Send transcribed text as user response
	response := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      transcribedText,
			Timestamp: time.Now(),
			IsFromAI:  false,
		},
		Timestamp: time.Now(),
	}

	s.sendMessage(response)

	// Store the answer and move to next question
	s.mutex.Lock()
	if s.currentIndex < len(s.answers) {
		s.answers[s.currentIndex] = transcribedText
	}
	s.currentIndex++
	s.mutex.Unlock()

	// Send next question after processing
	time.Sleep(2 * time.Second)
	s.sendNextQuestion()
}

func (s *InterviewSession) handleEndMessage() {
	s.endInterview()
}

func (s *InterviewSession) endInterview() {
	// Generate final analytics
	s.generateFinalAnalytics()

	endMsg := WSMessage{
		Type: MessageTypeEnd,
		Data: ResponseMessage{
			Text:      "Thank you for completing the mock interview! I'm analyzing your responses and will provide detailed feedback.",
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		Timestamp: time.Now(),
	}

	s.sendMessage(endMsg)

	s.close()
}

func (s *InterviewSession) generateFinalAnalytics() {

	questions := make([]string, len(s.questions))
	for i, question := range s.questions {
		questions[i] = question.Question
	}
	resp, err := s.aiClient.AnalyzeMockInterview(context.Background(), client.AnalyzeMockInterviewRequest{
		CV:             s.cv,
		VacancyInfo:    s.vacancyInfo,
		Specialization: s.specialization,
		Level:          s.level,
		Meta:           s.meta,
		Questions:      questions,
		Answers:        s.answers,
	})
	if err != nil {
		s.sendError("Failed to analyze mock interview")
		return
	}

	summaryMsg := WSMessage{
		Type: MessageTypeAnalytics,
		Data: ResponseMessage{
			Text:      resp.CandidateSummary,
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		AnalyticsData: resp,
		Timestamp:     time.Now(),
	}
	s.sendMessage(summaryMsg)
}

func (s *InterviewSession) sendMessage(msg WSMessage) {
	if !s.isActive || s.conn == nil {
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if err := s.conn.WriteJSON(msg); err != nil {
		log.Printf("Error sending message: %v", err)
		s.close()
	}
}

func (s *InterviewSession) sendError(message string) {
	errorMsg := WSMessage{
		Type:      MessageTypeError,
		Data:      map[string]string{"error": message},
		Timestamp: time.Now(),
	}

	s.sendMessage(errorMsg)
}

func RunServer(cfg *config.Config, aiClient *client.Client) error {
	// Initialize AI client for the global session

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		sessionMutex.Lock()
		defer sessionMutex.Unlock()

		// Close existing session if active
		if currentSession != nil && currentSession.isActive {
			currentSession.close()
		}

		// Upgrade HTTP connection to a WebSocket connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}

		// Create new interview session with AI client
		currentSession = &InterviewSession{
			conn:         conn,
			aiClient:     aiClient,
			questions:    make([]client.GeneratedQuestion, 0),
			answers:      make([]string, 0),
			currentIndex: 0,
			isActive:     true,
		}

		log.Printf("New mock interview session started with AI client")
		go currentSession.handleConnection()
	})

	log.Printf("Mock Interview WebSocket server starting on port %d", cfg.WSServerPort)
	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.WSServerPort), nil)
}
