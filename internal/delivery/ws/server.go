package ws

import (
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
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// StartMessage represents interview start request
type StartMessage struct {
	JobTitle       string `json:"job_title"`
	Experience     string `json:"experience"`
	Specialization string `json:"specialization"`
	Context        string `json:"context"`
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

// InterviewSession manages a single mock interview session
type InterviewSession struct {
	conn         *websocket.Conn
	aiClient     *client.Client
	isActive     bool
	currentIndex int
	questions    []string
	context      string
	jobTitle     string
	mutex        sync.RWMutex
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
			Text:      "Welcome to your mock interview! Please provide the job details and context to begin.",
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
	s.context = startMsg.Context
	s.jobTitle = startMsg.JobTitle
	s.mutex.Unlock()

	// Generate interview questions based on context
	go s.generateQuestions(startMsg)
}

func (s *InterviewSession) generateQuestions(startMsg StartMessage) {
	if s.aiClient == nil {
		s.sendError("AI client not initialized")
		return
	}

	// For now, using placeholder questions
	// TODO: Integrate with actual GPT client when available
	questions := []string{
		"Can you tell me about your experience with this technology stack?",
		"Describe a challenging project you've worked on recently.",
		"How do you approach problem-solving when faced with unknown requirements?",
		"What's your experience with team collaboration and communication?",
		"Where do you see yourself professionally in the next 2-3 years?",
	}

	s.mutex.Lock()
	s.questions = questions
	s.currentIndex = 0
	s.mutex.Unlock()

	// Send first question
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
			Text:       fmt.Sprintf("Question %d: %s", s.currentIndex+1, question),
			Timestamp:  time.Now(),
			IsFromAI:   true,
			IsQuestion: true,
		},
		Timestamp: time.Now(),
	}

	s.sendMessage(response)
}

func (s *InterviewSession) handleResponseMessage(msg WSMessage) {
	var responseMsg ResponseMessage
	if err := json.Unmarshal(msg.Data.([]byte), &responseMsg); err != nil {
		s.sendError("Invalid response message format")
		return
	}

	// Store the user's response (you could save this for analysis)
	log.Printf("User response: %s", responseMsg.Text)

	// Move to next question
	s.mutex.Lock()
	s.currentIndex++
	s.mutex.Unlock()

	// Send next question after a brief delay
	time.Sleep(1 * time.Second)
	s.sendNextQuestion()
}

func (s *InterviewSession) handleAudioMessage(msg WSMessage) {
	var audioMsg AudioMessage
	if err := json.Unmarshal(msg.Data.([]byte), &audioMsg); err != nil {
		s.sendError("Invalid audio message format")
		return
	}

	// Transcribe audio
	go s.transcribeAudio(audioMsg.AudioData)
}

func (s *InterviewSession) handleTranscribeMessage(msg WSMessage) {
	var audioMsg AudioMessage
	if err := json.Unmarshal(msg.Data.([]byte), &audioMsg); err != nil {
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

	// For now, simulate transcription
	// In a real implementation, you'd use your existing transcription functionality
	transcribedText := "[Audio transcription would appear here]"

	// Send transcription confirmation
	response := WSMessage{
		Type: MessageTypeResponse,
		Data: ResponseMessage{
			Text:      fmt.Sprintf("🎤 Transcribed: %s", transcribedText),
			Timestamp: time.Now(),
			IsFromAI:  false,
		},
		Timestamp: time.Now(),
	}

	s.sendMessage(response)

	// Move to next question
	s.mutex.Lock()
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
	endMsg := WSMessage{
		Type: MessageTypeEnd,
		Data: ResponseMessage{
			Text:      "Thank you for completing the mock interview! Your responses have been recorded.",
			Timestamp: time.Now(),
			IsFromAI:  true,
		},
		Timestamp: time.Now(),
	}

	s.sendMessage(endMsg)

	// Close the session after a delay
	time.Sleep(2 * time.Second)
	s.close()
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

func RunServer(cfg *config.Config) error {
	// Initialize AI client for the global session
	aiClient := client.New(cfg)

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
			questions:    make([]string, 0),
			currentIndex: -1,
			isActive:     true,
		}

		log.Printf("New mock interview session started with AI client")
		go currentSession.handleConnection()
	})

	log.Printf("Mock Interview WebSocket server starting on port %d", cfg.WSServerPort)
	return http.ListenAndServe(fmt.Sprintf(":%v", cfg.WSServerPort), nil)
}
