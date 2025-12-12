<template>
  <div class="mock-interview">
    <div class="interview-container">
      <h2>🎯 Mock Interview Session</h2>
      <p class="description">
        Practice your interview skills with AI-generated questions tailored to your CV and target position.
      </p>

      <!-- Interview Setup Form -->
      <div v-if="!interviewStarted" class="setup-form">
        <div class="form-group">
          <label for="cv">Your CV:</label>
          <textarea
            id="cv"
            v-model="interviewSetup.cv"
            placeholder="Paste your CV or resume here..."
            class="form-textarea"
            rows="6"
          ></textarea>
        </div>

        <div class="form-group">
          <label for="vacancy-info">Vacancy Information:</label>
          <textarea
            id="vacancy-info"
            v-model="interviewSetup.vacancyInfo"
            placeholder="Paste the job description or vacancy information here..."
            class="form-textarea"
            rows="6"
          ></textarea>
        </div>

        <div class="form-group">
          <label for="specialization">Specialization:</label>
          <input
            id="specialization"
            v-model="interviewSetup.specialization"
            type="text"
            placeholder="e.g., Backend Development, Data Science, DevOps"
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label for="level">Experience Level:</label>
          <select
            id="level"
            v-model="interviewSetup.level"
            class="form-select"
          >
            <option value="Junior">Junior (0-2 years)</option>
            <option value="Middle">Middle (2-5 years)</option>
            <option value="Senior">Senior (5-10 years)</option>
            <option value="Lead">Lead/Principal (10+ years)</option>
          </select>
        </div>

        <div class="form-group">
          <label for="questions-count">Number of Questions:</label>
          <select
            id="questions-count"
            v-model="interviewSetup.questionsCount"
            class="form-select"
          >
            <option :value="5">5 Questions</option>
            <option :value="10">10 Questions</option>
            <option :value="15">15 Questions</option>
          </select>
        </div>

        <div class="form-group">
          <label for="meta">Additional Context:</label>
          <textarea
            id="meta"
            v-model="interviewSetup.meta"
            placeholder="Any specific technologies, projects, or areas you'd like to focus on..."
            class="form-textarea"
            rows="4"
          ></textarea>
        </div>

        <button
          @click="startInterview"
          :disabled="!isSetupValid || isConnecting"
          class="start-button"
        >
          <span v-if="!isConnecting">🚀 Start Interview</span>
          <span v-else>⏳ Connecting...</span>
        </button>
      </div>

      <!-- Interview Interface -->
      <div v-if="interviewStarted" class="interview-interface">
        <!-- Connection Status -->
        <div class="connection-status">
          <div :class="['status-dot', { connected: isConnected, disconnected: !isConnected }]"></div>
          <span class="status-text">
            {{ isConnected ? 'Connected' : 'Disconnected' }}
          </span>
          <button
            @click="endInterview"
            class="end-interview-btn"
            v-if="isConnected"
          >
            🛑 End Interview
          </button>
        </div>


        <!-- Chat Messages -->
        <div class="chat-container" ref="chatContainer">
          <div
              v-for="message in messages"
              :key="message.timestamp"
              :class="[
      'message',
      {
        'ai-message': message.is_from_ai,
        'user-message': !message.is_from_ai,
        'question-message': message.isQuestion
      }
    ]"
          >
            <div class="message-content">
              <div class="message-header">
        <span class="sender">
          {{ message.is_from_ai ? '🤖 Interviewer' : '👤 You' }}
        </span>
                <span class="timestamp">{{ formatTime(message.timestamp) }}</span>
              </div>
              <div class="message-text">{{ message.text }}</div>
            </div>
          </div>

          <!-- Typing Indicator -->
          <div v-if="isTyping" class="typing-indicator">
            <div class="typing-dots">
              <span></span><span></span><span></span>
            </div>
            <span class="typing-text">Interviewer is typing...</span>
          </div>
        </div>

        <div v-if="analysisResult" class="analytics-panel">
          <h3>📊 Mock Interview Analysis</h3>

          <!-- Summary -->
          <div class="analytics-card">
            <div class="analytics-header">
      <span class="category">
        Level: {{ analysisResult.evaluation_level }}
      </span>
              <span :class="['accuracy', getAccuracyClass(analysisResult.final_score.average_accuracy)]">
        {{ Math.round(analysisResult.final_score.average_accuracy * 100) }}%
      </span>
            </div>

            <div class="analytics-content">
              <p class="feedback">{{ analysisResult.candidate_summary }}</p>
              <p><strong>Verdict:</strong> {{ analysisResult.final_score.verdict }}</p>
              <p>{{ analysisResult.final_score.verdict_reason }}</p>
            </div>
          </div>

          <!-- 🔥 SINGLE ITERATION -->
          <div class="analytics-grid">
            <div
                v-for="(q, idx) in analysisResult.questions_evaluation"
                :key="idx"
                class="analytics-card"
            >
              <div class="analytics-header">
                <span class="question-number">{{idx+1}}: {{ q.question }}</span>
                <span :class="['accuracy', getAccuracyClass(q.accuracy)]">
          {{ Math.round(q.accuracy * 100) }}%
        </span>
              </div>

              <div class="analytics-content">
                <p class="question">{{ q.question }}</p>
                <p class="answer">{{ q.answer }}</p>

                <p class="feedback">
                  <strong>Assessment:</strong> {{ q.assessment }}
                </p>

                <p v-if="q.reason_unanswered" class="feedback">
                  <strong>Reason:</strong> {{ q.reason_unanswered }}
                </p>

                <p v-if="q.what_was_expected" class="feedback">
                  <strong>Expected:</strong> {{ q.what_was_expected }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Input Area -->
        <div class="input-area">
          <div class="input-controls">
            <textarea
              v-model="userInput"
              @keydown.enter.prevent="sendMessage"
              @keydown.shift.enter.prevent="userInput += '\n'"
              placeholder="Type your answer here... (Enter to send, Shift+Enter for new line)"
              class="message-input"
              :disabled="!isConnected || interviewEnded"
              rows="3"
            ></textarea>
            
            <div class="input-actions">
              <button
                @click="toggleAudioRecording"
                :class="['audio-btn', { recording: isRecording }]"
                :disabled="!isConnected || interviewEnded"
                title="Record audio answer"
              >
                <span v-if="!isRecording">🎤</span>
                <span v-else>⏹️</span>
              </button>
              
              <button
                @click="sendMessage"
                :disabled="!userInput.trim() || !isConnected || interviewEnded"
                class="send-button"
              >
                📤 Send
              </button>
            </div>
          </div>
          
          <!-- Audio Recording Status -->
          <div v-if="isRecording" class="audio-status">
            <div class="recording-indicator">
              <div class="recording-dot"></div>
              <span>Recording... {{ recordingDuration }}s</span>
            </div>
            <button @click="stopAudioRecording" class="stop-recording-btn">
              Stop Recording
            </button>
          </div>
        </div>
      </div>

      <!-- Interview Summary -->
      <div v-if="interviewEnded && interviewSummary" class="interview-summary">
        <h3>📊 Interview Summary</h3>
        <div class="summary-content">
          <p><strong>Specialization:</strong> {{ interviewSetup.specialization }}</p>
          <p><strong>Level:</strong> {{ interviewSetup.level }}</p>
          <p><strong>Questions Asked:</strong> {{ messages.filter(m => m.isQuestion).length }}</p>
          <p><strong>Duration:</strong> {{ interviewDuration }}</p>

          <div v-if="analytics.length > 0" class="final-analytics">
            <h4>📈 Performance Overview</h4>
            <p><strong>Average Accuracy:</strong> {{ Math.round(averageAccuracy * 100) }}%</p>
            <p><strong>Questions Answered:</strong> {{ analytics.length }}/{{ interviewSetup.questionsCount }}</p>
          </div>
          <div class="summary-actions">
            <button @click="restartInterview" class="restart-button">
              🔄 Start New Interview
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { GetWebSocketURL } from '../../wailsjs/go/app/App'

const interviewSetup = ref({
  cv: '',
  vacancyInfo: '',
  specialization: '',
  level: 'Middle',
  meta: '',
  questionsCount: 10
})

const interviewStarted = ref(false)
const isConnecting = ref(false)
const isConnected = ref(false)
const interviewEnded = ref(false)
const interviewSummary = ref(null)
const interviewDuration = ref('')
const interviewStartTime = ref(null)

const messages = ref([])
const userInput = ref('')
const isTyping = ref(false)
const analytics = ref([])

const ws = ref(null)
const chatContainer = ref(null)

// Audio recording state
const isRecording = ref(false)
const recordingDuration = ref(0)
const recordingTimer = ref(null)
const audioChunks = ref([])

const mediaRecorder = ref(null)
const mediaStream = ref(null)

const analysisResult = ref(null)

const averageAccuracy = computed(() => {
  if (analytics.value.length === 0) return 0
  const sum = analytics.value.reduce((acc, curr) => acc + curr.accuracy, 0)
  return sum / analytics.value.length
})

const isSetupValid = computed(() => {
  return interviewSetup.value.cv.trim() !== '' && 
         interviewSetup.value.vacancyInfo.trim() !== '' &&
         interviewSetup.value.specialization.trim() !== '' &&
         interviewSetup.value.level !== '' &&
         interviewSetup.value.questionsCount > 0
})

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const getAccuracyClass = (accuracy) => {
  if (accuracy >= 0.8) return 'high'
  if (accuracy >= 0.6) return 'medium'
  return 'low'
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

const startInterview = async () => {
  isConnecting.value = true
  
  try {
    // Get WebSocket URL from backend
    const wsUrl = await GetWebSocketURL()
    ws.value = new WebSocket(wsUrl)
    
    ws.value.onopen = () => {
      console.log('WebSocket connected')
      isConnected.value = true
      isConnecting.value = false
      interviewStarted.value = true
      interviewStartTime.value = new Date()
      
      // Send interview setup with new structure
      ws.value.send(JSON.stringify({
        type: 'start',
        data: {
          cv: interviewSetup.value.cv,
          vacancy_info: interviewSetup.value.vacancyInfo,
          specialization: interviewSetup.value.specialization,
          level: interviewSetup.value.level,
          meta: interviewSetup.value.meta,
          questions_count: interviewSetup.value.questionsCount
        }
      }))
    }
    
    ws.value.onmessage = (event) => {
      const message = JSON.parse(event.data)
      handleMessage(message)
    }
    
    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
      isConnecting.value = false
      isConnected.value = false
    }
    
    ws.value.onclose = () => {
      console.log('WebSocket disconnected')
      isConnected.value = false
      if (interviewStarted.value && !interviewEnded.value) {
        endInterview()
      }
    }
    
  } catch (error) {
    console.error('Failed to start interview:', error)
    isConnecting.value = false
  }
}

const handleMessage = (message) => {
  isTyping.value = false
  
  switch (message.type) {
    case 'response':
      if (message.data) {
        messages.value.push({
          text: message.data.text,
          is_from_ai: message.data.is_from_ai,
          isQuestion: message.data.isQuestion,
          timestamp: new Date().toISOString()
        })
        scrollToBottom()
      }
      break

    case 'analytics': {
      const data = message.analytics_data
      if (data?.questions_evaluation) {
        analysisResult.value = data
        scrollToBottom()
      }
      break
    }
      
    case 'error':
      console.error('Server error:', message.data.error)
      messages.value.push({
        text: `Error: ${message.data.error}`,
        is_from_ai: true,
        timestamp: new Date().toISOString()
      })
      scrollToBottom()
      break
      
    case 'end':
      handleInterviewEnd()
      break
      
    default:
      console.log('Unknown message type:', message.type)
  }
}

const sendMessage = () => {
  if (!userInput.value.trim() || !isConnected.value || interviewEnded.value) {
    return
  }
  
  const messageText = userInput.value.trim()
  userInput.value = ''
  
  // Add user message to chat
  messages.value.push({
    text: messageText,
    is_from_ai: false,
    timestamp: new Date().toISOString()
  })
  
  // Send to WebSocket
  ws.value.send(JSON.stringify({
    type: 'response',
    data: {
      text: messageText,
      timestamp: new Date().toISOString(),
      is_from_ai: false
    }
  }))
  
  scrollToBottom()
  
  // Show typing indicator
  isTyping.value = true
}

const endInterview = () => {
  if (ws.value) {
    ws.value.send(JSON.stringify({ type: 'end' }))
  }
  handleInterviewEnd()
}

const handleInterviewEnd = () => {
  interviewEnded.value = true
  
  if (interviewStartTime.value) {
    const duration = Math.floor((new Date() - interviewStartTime.value) / 1000)
    const minutes = Math.floor(duration / 60)
    const seconds = duration % 60
    interviewDuration.value = `${minutes}m ${seconds}s`
  }
  
  interviewSummary.value = {
    specialization: interviewSetup.value.specialization,
    level: interviewSetup.value.level,
    questionsAsked: interviewSetup.value.questionsCount,
    duration: interviewDuration.value
  }

  isConnected.value = false
}

const restartInterview = () => {
  // Reset all state
  interviewStarted.value = false
  isConnecting.value = false
  isConnected.value = false
  interviewEnded.value = false
  interviewSummary.value = null
  interviewDuration.value = ''
  interviewStartTime.value = null
  messages.value = []
  userInput.value = ''
  isTyping.value = false
  analytics.value = []
  
  // Reset interview setup
  interviewSetup.value = {
    cv: '',
    vacancyInfo: '',
    specialization: '',
    level: 'Middle',
    meta: '',
    questionsCount: 10
  }
}

// Audio recording functions
const toggleAudioRecording = () => {
  if (isRecording.value) {
    stopAudioRecording()
  } else {
    startAudioRecording()
  }
}

const startAudioRecording = async () => {
  try {
    mediaStream.value = await navigator.mediaDevices.getUserMedia({ audio: true })
    mediaRecorder.value = new MediaRecorder(mediaStream.value)

    audioChunks.value = []
    recordingDuration.value = 0

    mediaRecorder.value.ondataavailable = (event) => {
      audioChunks.value.push(event.data)
    }

    mediaRecorder.value.onstop = () => {
      const audioBlob = new Blob(audioChunks.value, { type: 'audio/webm' })
      sendAudioData(audioBlob)

      // Stop all tracks
      mediaStream.value?.getTracks().forEach(t => t.stop())
      mediaStream.value = null
      mediaRecorder.value = null
    }

    mediaRecorder.value.start()
    isRecording.value = true

    recordingTimer.value = setInterval(() => {
      recordingDuration.value++
    }, 1000)
  } catch (error) {
    console.error('Error starting audio recording:', error)
  }
}

const stopAudioRecording = () => {
  if (recordingTimer.value) {
    clearInterval(recordingTimer.value)
    recordingTimer.value = null
  }

  if (mediaRecorder.value && mediaRecorder.value.state !== 'inactive') {
    mediaRecorder.value.stop()
  }

  isRecording.value = false
  recordingDuration.value = 0
}

const sendAudioData = async (audioBlob) => {
  try {
    const reader = new FileReader()
    reader.onloadend = () => {
      const arrayBuffer = reader.result
      
      // Send audio data to WebSocket
      ws.value.send(JSON.stringify({
        type: 'audio',
        data: {
          audio_data: Array.from(new Uint8Array(arrayBuffer)),
          format: 'webm'
        }
      }))
      
      // Show typing indicator while transcribing
      isTyping.value = true
    }
    reader.readAsArrayBuffer(audioBlob)
    
  } catch (error) {
    console.error('Error sending audio data:', error)
  }
}

onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
  if (recordingTimer.value) {
    clearInterval(recordingTimer.value)
  }
})
</script>

<style scoped>
.mock-interview {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.interview-container {
  background: white;
  border-radius: 12px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

h2 {
  color: #333;
  margin-bottom: 10px;
  font-size: 24px;
}

.description {
  color: #666;
  margin-bottom: 30px;
  line-height: 1.6;
}

/* Setup Form */
.setup-form {
  max-width: 600px;
  margin: 0 auto;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #333;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 12px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: #667eea;
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
  font-family: inherit;
}

.start-button {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 15px 30px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.start-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.start-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Interview Interface */
.connection-status {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 20px;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #dc3545;
}

.status-dot.connected {
  background: #28a745;
}

.status-dot.disconnected {
  background: #dc3545;
}

.status-text {
  font-weight: 600;
  color: #333;
}

.end-interview-btn {
  margin-left: auto;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
}

/* Analytics Panel */
.analytics-panel {
  margin-bottom: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 2px solid #e9ecef;
}

.analytics-panel h3 {
  margin: 0 0 15px 0;
  color: #333;
}

.analytics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;
  max-height: 300px;
  padding-top: 10px;
  overflow-y: auto;
}

.analytics-card {
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.analytics-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.question-number {
  font-weight: 600;
  color: #667eea;
}

.category {
  background: #667eea;
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.accuracy {
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.accuracy.high {
  background: #28a745;
  color: white;
}

.accuracy.medium {
  background: #ffc107;
  color: #212529;
}

.accuracy.low {
  background: #dc3545;
  color: white;
}

.analytics-content {
  font-size: 12px;
}

.analytics-content .question {
  font-weight: 600;
  margin-bottom: 5px;
}

.analytics-content .answer {
  color: #666;
  margin-bottom: 5px;
}

.analytics-content .feedback {
  color: #28a745;
  font-style: italic;
}

/* Chat Container */
.chat-container {
  height: 400px;
  overflow-y: auto;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  background: #fafbfc;
}

.message {
  margin-bottom: 20px;
  display: flex;
}

.ai-message {
  justify-content: flex-start;
}

.user-message {
  justify-content: flex-end;
}

.message-content {
  max-width: 70%;
  background: white;
  padding: 15px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.ai-message .message-content {
  background: #667eea;
  color: white;
  border-top-left-radius: 4px;
}

.user-message .message-content {
  background: #28a745;
  color: white;
  border-top-right-radius: 4px;
}

.question-message .message-content {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-left: 4px solid #ffd700;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 12px;
  opacity: 0.8;
}

.sender {
  font-weight: 600;
}

.timestamp {
  opacity: 0.7;
}

.message-text {
  line-height: 1.5;
  white-space: pre-wrap;
}

/* Typing Indicator */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  max-width: 200px;
}

.typing-dots {
  display: flex;
  gap: 4px;
}

.typing-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #667eea;
  animation: typing 1.4s infinite ease-in-out;
}

.typing-dots span:nth-child(1) {
  animation-delay: -0.32s;
}

.typing-dots span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes typing {
  0%, 80%, 100% {
    transform: scale(0);
    opacity: 0.5;
  }
  40% {
    transform: scale(1);
    opacity: 1;
  }
}

.typing-text {
  font-size: 14px;
  color: #666;
}

/* Input Area */
.input-area {
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  padding: 15px;
  background: white;
}

.input-controls {
  display: flex;
  gap: 10px;
  align-items: flex-end;
}

.message-input {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  resize: none;
  font-family: inherit;
}

.message-input:focus {
  outline: none;
  border-color: #667eea;
}

.input-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.audio-btn {
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px;
  cursor: pointer;
  font-size: 18px;
  transition: background 0.3s;
}

.audio-btn:hover:not(:disabled) {
  background: #5a6268;
}

.audio-btn.recording {
  background: #dc3545;
  animation: pulse 2s infinite;
}

.send-button {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 20px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: background 0.3s;
}

.send-button:hover:not(:disabled) {
  background: #5a6db8;
}

.send-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Audio Recording Status */
.audio-status {
  margin-top: 10px;
  padding: 10px;
  background: #fff3cd;
  border-radius: 6px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.recording-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.recording-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #dc3545;
  animation: blink 1s infinite;
}

.stop-recording-btn {
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 6px 12px;
  font-size: 12px;
  cursor: pointer;
}

/* Interview Summary */
.interview-summary {
  margin-top: 30px;
  padding: 20px;
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border-radius: 12px;
  text-align: center;
}

.interview-summary h3 {
  margin: 0 0 15px 0;
  font-size: 20px;
}

.summary-content p {
  margin: 8px 0;
  font-size: 16px;
}

.final-analytics {
  margin: 20px 0;
  padding: 15px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.final-analytics h4 {
  margin: 0 0 10px 0;
}

.summary-actions {
  margin-top: 20px;
}

.restart-button {
  background: white;
  color: #28a745;
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.restart-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(255, 255, 255, 0.3);
}

/* Animations */
@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(220, 53, 69, 0.7);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(220, 53, 69, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(220, 53, 69, 0);
  }
}

@keyframes blink {
  0%, 50% {
    opacity: 1;
  }
  25%, 75% {
    opacity: 0.3;
  }
}

/* Responsive */
@media (max-width: 768px) {
  .mock-interview {
    padding: 10px;
  }

  .interview-container {
    padding: 20px;
  }

  .message-content {
    max-width: 85%;
  }

  .input-controls {
    flex-direction: column;
  }

  .message-input {
    width: 100%;
  }

  .input-actions {
    width: 100%;
    justify-content: space-between;
  }

  .analytics-grid {
    grid-template-columns: 1fr;
  }
}
</style>
