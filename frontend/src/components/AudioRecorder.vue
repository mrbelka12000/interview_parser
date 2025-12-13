<template>
  <div class="audio-recorder">
    <div class="recorder-container">
      <h2>🎙️ Real-time Audio Recording</h2>
      <p class="description">
        Record audio directly from your microphone and process it for transcription and analysis.
        Choose between interview analysis or call/meeting analysis modes.
      </p>

      <!-- Recording Controls -->
      <div class="recording-controls">
        <div class="main-controls">
          <button
            @click="toggleRecording"
            :class="['record-button', { recording: isRecording, stopping: isStopping }]"
            :disabled="isProcessing || isStopping"
          >
            <span v-if="!isRecording && !isStopping" class="record-icon">🎤</span>
            <span v-else-if="isRecording" class="record-icon">⏹️</span>
            <span v-else class="record-icon">⏳</span>

            <span v-if="!isRecording && !isStopping" class="record-text">Start Recording</span>
            <span v-else-if="isRecording" class="record-text">Stop Recording</span>
            <span v-else class="record-text">Stopping...</span>
          </button>

          <button
            @click="resetRecording"
            :disabled="isRecording || isStopping || isProcessing"
            class="reset-button"
          >
            🔄 Reset
          </button>
        </div>

        <!-- Recording Status -->
        <div v-if="isRecording || recordingStatus.duration > 0" class="recording-status">
          <div class="status-indicator">
            <div :class="['status-dot', { active: isRecording }]"></div>
            <span class="status-text">
              {{ isRecording ? 'Recording' : 'Recorded' }}: {{ formatDuration(recordingStatus.duration) }}
            </span>
          </div>

          <div v-if="recordingStatus.dataSize > 0" class="status-info">
            <span>Size: {{ formatFileSize(recordingStatus.dataSize) }}</span>
          </div>
        </div>
      </div>

      <!-- Device Selection -->
      <div class="device-selection">
        <h3>🎛️ Audio Devices</h3>
        
        <div class="device-row">
          <div class="device-group">
            <label for="input-device-select">Input Device (Microphone):</label>
            <select
              id="input-device-select"
              v-model="selectedInputDevice"
              @change="onInputDeviceChange"
              :disabled="isRecording || isStopping || isProcessing"
              class="device-select"
            >
              <option value="">Default Input Device</option>
              <option
                v-for="device in inputDevices"
                :key="device.id"
                :value="device.id"
              >
                {{ device.name }} {{ device.isDefault ? '(Default)' : '' }}
              </option>
            </select>
          </div>
        </div>

        <button
          @click="refreshDevices"
          :disabled="isRecording || isStopping || isProcessing"
          class="refresh-devices-btn"
        >
          🔄 Refresh Devices
        </button>
      </div>

      <!-- Recording Options -->
      <div class="recording-options">
        <div class="option-group">
          <label for="mode-select">Processing Mode:</label>
          <select
            id="mode-select"
            v-model="processingMode"
            :disabled="isRecording || isStopping || isProcessing"
            class="mode-select"
          >
            <option value="interview">👥 Interview Analysis</option>
            <option value="call">📞 Call/Meeting Analysis</option>
          </select>
          <small v-if="processingMode === 'call'">
            Call mode analyzes daily meetings and creates actionable plans with tasks, open questions, and next steps.
          </small>
          <small v-else>
            Interview mode extracts questions and evaluates candidate responses.
          </small>
        </div>

        <div class="option-group">
          <label for="filename-input">Filename (optional):</label>
          <input
            id="filename-input"
            v-model="customFilename"
            type="text"
            :placeholder="processingMode === 'call' ? 'meeting_recording' : 'interview_recording'"
            :disabled="isRecording || isStopping || isProcessing"
            class="filename-input"
          />
          <small>.wav will be added automatically</small>
        </div>

        <div class="option-group">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="autoProcess"
              :disabled="isRecording || isStopping || isProcessing"
            />
            Automatically process recording after save
          </label>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="action-buttons">
        <button
          @click="saveRecording"
          :disabled="!hasRecording || isRecording || isStopping || isProcessing"
          class="save-button"
        >
          💾 Save Recording
        </button>

        <button
          @click="saveAndProcess"
          :disabled="!hasRecording || isRecording || isStopping || isProcessing"
          class="process-button"
        >
          🚀 Save & Process
        </button>
      </div>

      <!-- Progress Indicator -->
      <div v-if="isProcessing" class="progress-container">
        <h3>🔄 Processing Recording...</h3>
        <div class="progress-bar-container">
          <div class="progress-bar" :style="{ width: progressPercentage + '%' }"></div>
        </div>
        <div class="progress-info">
          <span class="progress-percentage">{{ progressPercentage }}%</span>
          <span class="progress-stage">{{ progressStage }}</span>
        </div>
        <div class="progress-details">
          <p>{{ progressDetails }}</p>
        </div>
      </div>

      <!-- Results -->
      <div v-if="recordingResult" class="result">
        <h3>Recording Result:</h3>
        <div :class="['result-message', recordingResult.success ? 'success' : 'error']">
          {{ recordingResult.message }}
        </div>

        <div v-if="recordingResult.success && recordingResult.filePath" class="file-info">
          <p><strong>📁 File:</strong> {{ getFileName(recordingResult.filePath) }}</p>
          <p><strong>⏱️ Duration:</strong> {{ formatDuration(recordingResult.duration) }}</p>
          <p><strong>💾 Size:</strong> {{ formatFileSize(recordingResult.dataSize) }}</p>
          <p><strong>📂 Path:</strong> {{ recordingResult.filePath }}</p>
        </div>
      </div>

      <!-- Transcription Results -->
      <div v-if="transcriptionResult" class="result">
        <h3>🎉 Processing Complete!</h3>
        <div :class="['result-message', transcriptionResult.success ? 'success' : 'error']">
          {{ transcriptionResult.message }}
        </div>

        <div v-if="transcriptionResult.success" class="output-files">
          <h4>✅ Generated Files:</h4>
          <div class="file-links">
            <div v-if="processingMode === 'interview'" class="file-link">
              <strong>📝 Transcript:</strong>
              <a href="#" @click.prevent="openFile(transcriptionResult.transcriptPath)">
                {{ getFileName(transcriptionResult.transcriptPath) }}
              </a>
            </div>
            <div class="file-link">
              <strong>📊 Analysis:</strong>
              <a href="#" @click.prevent="openFile(transcriptionResult.analysisPath)">
                {{ getFileName(transcriptionResult.analysisPath) }}
              </a>
            </div>
          </div>
          <div class="completion-notice">
            <p v-if="processingMode === 'call'">🎉 Your meeting has been successfully transcribed and analyzed!</p>
            <p v-else>🎉 Your interview has been successfully transcribed and analyzed!</p>
          </div>
        </div>
      </div>

      <!-- File Content Modal -->
      <div v-if="showFileContent && fileToOpen" class="file-content-overlay" @click="closeFileContent">
        <div class="file-content-modal" @click.stop>
          <div class="modal-header">
            <h3>{{ getFileName(fileToOpen) }}</h3>
            <button @click="closeFileContent" class="close-btn">×</button>
          </div>
          <FileContent :file-path="fileToOpen" @back="closeFileContent" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import {
  StartAudioRecording,
  StopAudioRecording,
  SaveRecording,
  SaveAndProcessRecording,
  SaveAndProcessRecordingForCall,
  GetRecordingStatus,
  GetInputDevices,
  SetAudioInputDevice,
} from '../../wailsjs/go/wails_app/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import FileContent from './FileContent.vue'

const isRecording = ref(false)
const isStopping = ref(false)
const isProcessing = ref(false)
const customFilename = ref('')
const autoProcess = ref(true)
const processingMode = ref('interview')
const recordingResult = ref(null)
const transcriptionResult = ref(null)
const recordingStatus = ref({ duration: 0, dataSize: 0 })
const progressPercentage = ref(0)
const progressStage = ref('')
const progressDetails = ref('')
const fileToOpen = ref(null)
const showFileContent = ref(false)
const volume = ref(1.0) // Default volume is 100%

// Device selection state
const inputDevices = ref([])
const selectedInputDevice = ref('')

let statusInterval = null

const hasRecording = computed(() => {
  return recordingStatus.value.duration > 0 || recordingStatus.value.dataSize > 0
})

const getFileName = (filePath) => {
  if (!filePath) return ''
  return filePath.split(/[\\/]/).pop() || filePath
}

const formatDuration = (seconds) => {
  if (!seconds || seconds <= 0) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const toggleRecording = async () => {
  if (isRecording.value) {
    await stopRecording()
  } else {
    await startRecording()
  }
}

const startRecording = async () => {
  try {
    const result = await StartAudioRecording()

    if (result.success) {
      isRecording.value = true
      recordingResult.value = null
      transcriptionResult.value = null

      // Start status polling
      startStatusPolling()

      console.log('Recording started')
    } else {
      recordingResult.value = result
      console.error('Failed to start recording:', result.message)
    }
  } catch (error) {
    console.error('Error starting recording:', error)
    recordingResult.value = {
      success: false,
      message: `Error: ${error.message || error.toString()}`
    }
  }
}

const stopRecording = async () => {
  isStopping.value = true

  try {
    const result = await StopAudioRecording()

    if (result.success) {
      isRecording.value = false
      recordingStatus.value.duration = result.duration || 0
      recordingStatus.value.dataSize = result.dataSize || 0

      // Stop status polling
      stopStatusPolling()

      recordingResult.value = {
        success: true,
        message: 'Recording stopped successfully',
        duration: result.duration,
        dataSize: result.dataSize
      }

      console.log('Recording stopped:', result)

      // Auto-process if enabled
      if (autoProcess.value) {
        setTimeout(() => {
          saveAndProcess()
        }, 1000)
      }
    } else {
      recordingResult.value = result
      console.error('Failed to stop recording:', result.message)
    }
  } catch (error) {
    console.error('Error stopping recording:', error)
    recordingResult.value = {
      success: false,
      message: `Error: ${error.message || error.toString()}`
    }
  } finally {
    isStopping.value = false
  }
}

const resetRecording = async () => {
  if (isRecording.value) {
    await stopRecording()
  }

  recordingResult.value = null
  transcriptionResult.value = null
  recordingStatus.value = { duration: 0, dataSize: 0 }
  customFilename.value = ''
}

const saveRecording = async () => {
  isProcessing.value = true

  try {
    const result = await SaveRecording(customFilename.value)
    recordingResult.value = result

    if (result.success) {
      console.log('Recording saved:', result)
    } else {
      console.error('Failed to save recording:', result.message)
    }
  } catch (error) {
    console.error('Error saving recording:', error)
    recordingResult.value = {
      success: false,
      message: `Error: ${error.message || error.toString()}`
    }
  } finally {
    isProcessing.value = false
  }
}

const saveAndProcess = async () => {
  isProcessing.value = true

  // Reset progress
  progressPercentage.value = 0
  progressStage.value = 'Initializing...'
  progressDetails.value = 'Preparing to process your recording...'

  try {
    let result
    if (processingMode.value === 'call') {
      result = await SaveAndProcessRecordingForCall(customFilename.value)
    } else {
      result = await SaveAndProcessRecording(customFilename.value)
    }
    
    transcriptionResult.value = result

    if (result.success) {
      console.log('Recording processed:', result)
    } else {
      console.error('Failed to process recording:', result.message)
    }
  } catch (error) {
    console.error('Error processing recording:', error)
    transcriptionResult.value = {
      success: false,
      message: `Error: ${error.message || error.toString()}`
    }
  } finally {
    isProcessing.value = false
  }
}

const startStatusPolling = () => {
  statusInterval = setInterval(async () => {
    try {
      const result = await GetRecordingStatus()
      if (result.success) {
        recordingStatus.value.duration = result.duration || 0
        recordingStatus.value.dataSize = result.dataSize || 0
      }
    } catch (error) {
      console.error('Error getting recording status:', error)
    }
  }, 500) // Poll every 500ms
}

const stopStatusPolling = () => {
  if (statusInterval) {
    clearInterval(statusInterval)
    statusInterval = null
  }
}

// Handle progress events from backend
const handleProgressEvent = (data) => {
  console.log('Progress event received:', data)
  if (data && data.percentage !== undefined) {
    progressPercentage.value = data.percentage
    progressStage.value = data.stage || 'Processing...'
    progressDetails.value = data.details || 'Working on your recording...'
  }
}

const openFile = async (filePath) => {
  try {
    console.log('Opening file:', filePath)
    fileToOpen.value = filePath
    showFileContent.value = true
  } catch (error) {
    console.error('Error opening file:', error)
  }
}

const closeFileContent = () => {
  showFileContent.value = false
  fileToOpen.value = null
}

const setVolume = async () => {
  try {
    const result = await SetAudioVolume(volume.value)
    if (result.success) {
      console.log('Volume set to:', result.volume)
    } else {
      console.error('Failed to set volume:', result.message)
    }
  } catch (error) {
    console.error('Error setting volume:', error)
  }
}

// Device management functions
const loadDevices = async () => {
  try {
    // Load input devices
    const inputResult = await GetInputDevices()
    if (inputResult.success) {
      inputDevices.value = inputResult.devices
      console.log('Input devices loaded:', inputDevices.value)
    } else {
      console.error('Failed to load input devices:', inputResult.message)
    }
  } catch (error) {
    console.error('Error loading devices:', error)
  }
}

const refreshDevices = async () => {
  if (isRecording.value || isStopping.value || isProcessing.value) {
    return
  }

  await loadDevices()
}

const onInputDeviceChange = async () => {
  if (isRecording.value || isStopping.value || isProcessing.value) {
    return
  }

  if (!selectedInputDevice.value) {
    // Reset to default device
    console.log('Using default input device')
    return
  }

  try {
    const result = await SetAudioInputDevice(selectedInputDevice.value)
    if (result.success) {
      console.log('Input device changed successfully')
    } else {
      console.error('Failed to change input device:', result.message)
      // Reset selection on failure
      selectedInputDevice.value = ''
    }
  } catch (error) {
    console.error('Error changing input device:', error)
    selectedInputDevice.value = ''
  }
}


const initializeVolume = async () => {
  try {
    const result = await GetAudioVolume()
    if (result.success) {
      volume.value = result.volume
      console.log('Initial volume:', result.volume)
    } else {
      console.error('Failed to get volume:', result.message)
    }
  } catch (error) {
    console.error('Error getting volume:', error)
  }
}

onMounted(async () => {
  EventsOn('progress', handleProgressEvent)
  await initializeVolume()
  await loadDevices()
})

onUnmounted(() => {
  stopStatusPolling()
})
</script>

<style scoped>
.audio-recorder {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.recorder-container {
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

/* Recording Controls */
.recording-controls {
  margin-bottom: 30px;
}

.main-controls {
  display: flex;
  gap: 20px;
  justify-content: center;
  margin-bottom: 20px;
}

.record-button {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  color: white;
  border: none;
  border-radius: 50px;
  padding: 20px 40px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 10px;
  box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3);
}

.record-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(231, 76, 60, 0.4);
}

.record-button.recording {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  animation: pulse 2s infinite;
}

.record-button.stopping {
  background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%);
}

.record-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.record-icon {
  font-size: 24px;
}

.record-text {
  font-size: 16px;
}

.reset-button {
  background: linear-gradient(135deg, #95a5a6 0%, #7f8c8d 100%);
  color: white;
  border: none;
  border-radius: 50px;
  padding: 15px 30px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.reset-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(149, 165, 166, 0.3);
}

.reset-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Recording Status */
.recording-status {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  text-align: center;
}

.status-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 10px;
}

.status-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #dc3545;
  transition: background 0.3s ease;
}

.status-dot.active {
  background: #28a745;
  animation: blink 1s infinite;
}

.status-text {
  font-weight: 600;
  color: #333;
}

.status-info {
  font-size: 14px;
  color: #666;
}

/* Device Selection */
.device-selection {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 30px;
}

.device-selection h3 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 18px;
  font-weight: 600;
}

.device-row {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.device-group {
  flex: 1;
}

.device-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.device-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

.device-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.device-select:disabled {
  background: #e9ecef;
  color: #6c757d;
  cursor: not-allowed;
}

.refresh-devices-btn {
  background: linear-gradient(135deg, #6c757d 0%, #495057 100%);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.refresh-devices-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(108, 117, 125, 0.3);
}

.refresh-devices-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Recording Options */
.recording-options {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 30px;
}

.option-group {
  margin-bottom: 15px;
}

.option-group:last-child {
  margin-bottom: 0;
}

.option-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.filename-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.filename-input:focus {
  outline: none;
  border-color: #667eea;
}

.option-group small {
  display: block;
  margin-top: 5px;
  color: #666;
  font-size: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #333;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

.volume-slider {
  width: 100%;
  height: 6px;
  border-radius: 3px;
  background: #ddd;
  outline: none;
  -webkit-appearance: none;
  appearance: none;
  cursor: pointer;
}

.volume-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.volume-slider::-webkit-slider-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.volume-slider::-moz-range-thumb {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  cursor: pointer;
  border: none;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.volume-slider::-moz-range-thumb:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.volume-slider:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.volume-slider:disabled::-webkit-slider-thumb {
  cursor: not-allowed;
  transform: scale(1);
}

.volume-slider:disabled::-moz-range-thumb {
  cursor: not-allowed;
  transform: scale(1);
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 15px;
  margin-bottom: 30px;
}

.save-button, .process-button {
  flex: 1;
  padding: 15px 20px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.save-button {
  background: linear-gradient(135deg, #17a2b8 0%, #138496 100%);
  color: white;
}

.save-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(23, 162, 184, 0.3);
}

.process-button {
  background: linear-gradient(135deg, #28a745 0%, #218838 100%);
  color: white;
}

.process-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(40, 167, 69, 0.3);
}

.save-button:disabled, .process-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Progress Bar (same as FileUpload) */
.progress-container {
  margin-top: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #dee2e6;
}

.progress-container h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 18px;
  text-align: center;
}

.progress-bar-container {
  width: 100%;
  height: 20px;
  background: #e9ecef;
  border-radius: 10px;
  overflow: hidden;
  margin-bottom: 15px;
  position: relative;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  transition: width 0.3s ease;
  min-width: 2%;
  position: relative;
  overflow: hidden;
}

.progress-bar::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.3),
    transparent
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.progress-percentage {
  font-size: 18px;
  font-weight: 600;
  color: #667eea;
}

.progress-stage {
  font-size: 14px;
  color: #6c757d;
  font-weight: 500;
}

.progress-details {
  text-align: center;
}

.progress-details p {
  margin: 0;
  font-size: 14px;
  color: #495057;
  font-style: italic;
}

/* Results */
.result {
  margin-top: 30px;
  padding: 20px;
  border-radius: 8px;
  background: #f8f9fa;
}

.result-message {
  padding: 12px;
  border-radius: 6px;
  font-weight: 500;
}

.result-message.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.result-message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.file-info {
  margin-top: 15px;
  padding: 15px;
  background: white;
  border-radius: 6px;
}

.file-info p {
  margin: 5px 0;
  color: #333;
}

.file-info strong {
  color: #495057;
}

.output-files {
  margin-top: 15px;
}

.file-links {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.file-link {
  padding: 10px;
  background: white;
  border-radius: 4px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.file-link a:hover {
  text-decoration: underline;
}

.completion-notice {
  margin-top: 20px;
  padding: 15px;
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border-radius: 8px;
  text-align: center;
  box-shadow: 0 4px 15px rgba(40, 167, 69, 0.2);
}

.completion-notice p {
  margin: 5px 0;
  font-weight: 500;
}

/* Animations */
@keyframes pulse {
  0% {
    box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3);
  }
  50% {
    box-shadow: 0 4px 25px rgba(231, 76, 60, 0.6);
  }
  100% {
    box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3);
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

/* File Content Modal (same as FileUpload) */
.file-content-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.file-content-modal {
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  max-width: 90vw;
  max-height: 90vh;
  width: 1200px;
  height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #dee2e6;
  background: #f8f9fa;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #6c757d;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.3s;
}

.close-btn:hover {
  background: #e9ecef;
  color: #495057;
}

/* Responsive */
@media (max-width: 768px) {
  .audio-recorder {
    padding: 10px;
  }

  .recorder-container {
    padding: 20px;
  }

  .main-controls {
    flex-direction: column;
    align-items: center;
  }

  .device-row {
    flex-direction: column;
    gap: 15px;
  }

  .action-buttons {
    flex-direction: column;
  }

  .file-content-modal {
    width: 95vw;
    height: 85vh;
    margin: 10px;
  }
}
.audio-recorder {
   max-width: 800px;
   margin: 0 auto;
   padding: 20px;
   color: #000; /* << make all text inside black */
 }
</style>
