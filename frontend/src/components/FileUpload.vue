<template>
  <div class="file-upload">
    <div class="upload-container">
      <h2>🎤 Upload Interview for Transcription</h2>
      <p class="description">
        Upload an audio or video file to transcribe and analyze the interview.
        The file will be processed in chunks using AI transcription and then analyzed for questions.
      </p>

      <!-- Single upload area: click = OpenFileDialog, drag&drop still works -->
      <div class="upload-area"
           :class="{ 'drag-over': isDragOver }"
           @dragover.prevent="isDragOver = true"
           @dragleave.prevent="isDragOver = false"
           @drop.prevent="handleDrop"
           @click="pickFile">

        <div class="upload-content">
          <div class="upload-icon">📁</div>
          <div class="upload-text">
            <p>Drag and drop your file here, or click to browse</p>
            <p class="file-types">Supported formats: MP3, WAV, M4A, MP4, MOV, AVI</p>
          </div>
        </div>
      </div>

      <div v-if="selectedFile" class="file-info">
        <h3>Selected File:</h3>
        <div class="file-details">
          <p><strong>Name:</strong> {{ selectedFile.name }}</p>
          <p v-if="selectedFile.size"><strong>Size:</strong> {{ formatFileSize(selectedFile.size) }}</p>
          <p v-if="selectedFile.type"><strong>Type:</strong> {{ selectedFile.type }}</p>
          <p v-if="selectedFilePath"><strong>Full path:</strong> {{ selectedFilePath }}</p>
        </div>
      </div>

      <div class="options">
        <label class="checkbox-label">
          <input
              type="checkbox"
              v-model="loadChunks"
          />
          Load existing chunks (if available)
        </label>
      </div>

      <button
          @click="uploadFile"
          :disabled="!selectedFilePath || isProcessing"
          class="upload-button"
      >
        <span v-if="!isProcessing">🚀 Process File</span>
        <span v-else>⏳ Processing... ({{ progressText }})</span>
      </button>

      <!-- Progress Bar with Percentage -->
      <div v-if="isProcessing" class="progress-container">
        <h3>🔄 Processing Progress</h3>
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

      <div v-if="result" class="result">
        <h3>Processing Result:</h3>
        <div :class="['result-message', result.success ? 'success' : 'error']">
          {{ result.message }}
        </div>

        <div v-if="result.success" class="output-files">
          <h4>✅ Processing Complete! Generated Files:</h4>
          <div class="file-links">
            <div class="file-link">
              <strong>📝 Transcript:</strong>
              <a href="#" @click.prevent="openFile(result.transcriptPath)">
                {{ getFileName(result.transcriptPath) }}
              </a>
            </div>
            <div class="file-link">
              <strong>📊 Analysis:</strong>
              <a href="#" @click.prevent="openFile(result.analysisPath)">
                {{ getFileName(result.analysisPath) }}
              </a>
            </div>
          </div>
          <div class="completion-notice">
            <p>🎉 Your interview has been successfully transcribed and analyzed!</p>
            <p>You can now view the transcript and analysis using the links above.</p>
          </div>
        </div>
      </div>

      <!-- File Content Modal/Overlay -->
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
import { ref, onMounted, onUnmounted } from 'vue'
import { ProcessFileForTranscription, PickFile } from '../../wailsjs/go/wails_app/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import FileContent from './FileContent.vue'

const isDragOver = ref(false)
const selectedFile = ref(null)          // for UI (name/size/type)
const selectedFilePath = ref('')        // full path for backend
const isProcessing = ref(false)
const loadChunks = ref(false)
const result = ref(null)
const progressText = ref('')
const progressPercentage = ref(0)
const progressStage = ref('')
const progressDetails = ref('')
const fileToOpen = ref(null)
const showFileContent = ref(false)

const getFileName = (filePath) => {
  if (!filePath) return ''
  return filePath.split(/[\\/]/).pop() || filePath
}

const formatFileSize = (bytes) => {
  if (!bytes) return ''
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Click on upload area → call Go's PickFile()
const pickFile = async () => {
  try {
    const path = await PickFile()   // generated by Wails

    console.log('PickFile() returned:', path)

    if (!path) {
      // user cancelled dialog
      return
    }

    selectedFilePath.value = path
    selectedFile.value = {
      name: getFileName(path),
      size: 0,
      type: 'Unknown',
      path,
    }
    result.value = null
  } catch (err) {
    console.error('PickFile error:', err)
  }
}

// Drag & drop: use file.path if Wails provides it, else just name
const handleDrop = (event) => {
  isDragOver.value = false
  event.preventDefault()

  const files = event.dataTransfer?.files
  if (!files || files.length === 0) return

  const file = files[0]

  selectedFile.value = file
  selectedFilePath.value = file.path || file.name
  result.value = null

  console.log('Dropped file:', {
    name: file.name,
    path: file.path,
    size: file.size,
    type: file.type,
  })
}

const uploadFile = async () => {
  if (!selectedFilePath.value) return

  isProcessing.value = true
  result.value = null
  progressText.value = 'Processing file...'

  // Reset progress
  progressPercentage.value = 0
  progressStage.value = 'Initializing...'
  progressDetails.value = 'Preparing to process your file...'

  try {
    console.log('Processing file:', selectedFilePath.value)
    
    // Actual processing with real progress from backend
    result.value = await ProcessFileForTranscription(
        selectedFilePath.value,
        loadChunks.value,
    )

    console.log('Backend result:', result.value)
  } catch (error) {
    progressPercentage.value = 0
    progressStage.value = 'Error'
    progressDetails.value = 'Processing failed'
    result.value = {
      success: false,
      message: `Error: ${error.message || error.toString()}`,
    }
  } finally {
    isProcessing.value = false
    progressText.value = ''
  }
}

// Handle progress events from backend
const handleProgressEvent = (data) => {
  console.log('Progress event received:', data)
  if (data && data.percentage !== undefined) {
    progressPercentage.value = data.percentage
    progressStage.value = data.stage || 'Processing...'
    progressDetails.value = data.details || 'Working on your file...'
    
    if (data.isComplete) {
      progressText.value = 'Complete!'
    }
  }
}

// Set up event listener for progress updates
onMounted(() => {
  EventsOn('progress', handleProgressEvent)
})

// Clean up event listener
onUnmounted(() => {
  // Note: Wails doesn't provide a way to remove event listeners in current version
  // This will be cleaned up when component is unmounted naturally
})

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
</script>

<style scoped>
.file-upload {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.upload-container {
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

.upload-area {
  border: 3px dashed #ddd;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.upload-area:hover {
  border-color: #667eea;
  background: #f0f4ff;
}

.upload-area.drag-over {
  border-color: #667eea;
  background: #e8efff;
  transform: scale(1.02);
}

.upload-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.upload-text p {
  margin: 5px 0;
}

.file-types {
  font-size: 14px;
  color: #888;
}

.file-info {
  margin-top: 20px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 6px;
}

.file-details p {
  margin: 5px 0;
}

.options {
  margin: 20px 0;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

/* make checkbox text black */
.checkbox-label {
  color: black;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

.upload-button {
  width: 100%;
  padding: 15px 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 20px;
}

.upload-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.upload-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Progress Bar Styles */
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

@media (max-width: 600px) {
  .file-upload {
    padding: 10px;
  }

  .upload-container {
    padding: 20px;
  }

  .upload-area {
    padding: 30px 15px;
  }

  .upload-icon {
    font-size: 36px;
  }
}

/* file-info text black */
.file-info,
.file-info h3,
.file-info p,
.file-info strong {
  color: #000 !important;
}

.checkbox-label input {
  color: inherit;
}

/* File Content Modal Styles */
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

.file-content-modal :deep(.file-content) {
  flex: 1;
  padding: 0;
  max-width: none;
  max-height: none;
}

.file-content-modal :deep(.content) {
  height: 100%;
  max-height: none;
}

@media (max-width: 768px) {
  .file-content-modal {
    width: 95vw;
    height: 85vh;
    margin: 10px;
  }
  
  .modal-header {
    padding: 15px;
  }
  
  .modal-header h3 {
    font-size: 16px;
  }
  
  .close-btn {
    font-size: 20px;
    width: 25px;
    height: 25px;
  }
}

.output-files {
  color: #000;
}
.output-files a {
  color: #0a58ca; /* optional link color */
}
</style>
