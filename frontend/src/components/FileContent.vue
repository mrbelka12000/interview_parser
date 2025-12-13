<script setup>
import { ref, onMounted } from 'vue'
import { ReadFileContent } from '../../wailsjs/go/wails_app/App'

const props = defineProps({
  filePath: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['back'])

const fileContent = ref(null)
const loading = ref(true)
const error = ref('')

// File icon mapping (same as FileExplorer)
const fileIcons = {
  // Directories
  dir: '📁',

  // Code files
  '.js': '🟨',
  '.ts': '🔷',
  '.vue': '💚',
  '.go': '🐹',
  '.html': '🌐',
  '.css': '🎨',
  '.scss': '🎨',
  '.json': '📋',
  '.xml': '📋',
  '.yaml': '📋',
  '.yml': '📋',

  // Documents
  '.pdf': '📕',
  '.doc': '📘',
  '.docx': '📘',
  '.txt': '📄',
  '.md': '📝',
  '.rtf': '📄',

  // Images
  '.jpg': '🖼️',
  '.jpeg': '🖼️',
  '.png': '🖼️',
  '.gif': '🖼️',
  '.svg': '🎨',
  '.ico': '🖼️',
  '.bmp': '🖼️',

  // Video
  '.mp4': '🎬',
  '.avi': '🎬',
  '.mov': '🎬',
  '.wmv': '🎬',
  '.flv': '🎬',
  '.mkv': '🎬',

  // Audio
  '.mp3': '🎵',
  '.wav': '🎵',
  '.flac': '🎵',
  '.aac': '🎵',
  '.ogg': '🎵',

  // Archives
  '.zip': '📦',
  '.rar': '📦',
  '.tar': '📦',
  '.gz': '📦',
  '.7z': '📦',

  // Default
  'file': '📄'
}

// Get icon for file based on extension
const getFileIcon = (file) => {
  if (file.isDir) {
    return fileIcons.dir
  }
  return fileIcons[file.extension.toLowerCase()] || fileIcons.file
}

// Format file size
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Detect if content is likely binary
const isBinaryContent = (content) => {
  return false
  // Check for null bytes and other binary indicators
  // if (content.includes('\x00')) return true
  if (content.length > 0) {
    const nonTextChars = content.replace(/[\x09\x0A\x0D\x20-\x7E]/g, '').length
    const ratio = nonTextChars / content.length
    return ratio > 0.3 // If more than 30% non-text chars, consider it binary
  }
  return false
}

// Get syntax highlighting class based on file extension
const getSyntaxClass = (extension) => {
  const syntaxMap = {
    '.js': 'language-javascript',
    '.ts': 'language-typescript',
    '.vue': 'language-html',
    '.go': 'language-go',
    '.html': 'language-html',
    '.css': 'language-css',
    '.scss': 'language-scss',
    '.json': 'language-json',
    '.xml': 'language-xml',
    '.yaml': 'language-yaml',
    '.yml': 'language-yaml',
    '.md': 'language-markdown',
    '.py': 'language-python',
    '.java': 'language-java',
    '.cpp': 'language-cpp',
    '.c': 'language-c',
    '.h': 'language-c',
    '.php': 'language-php',
    '.rb': 'language-ruby',
    '.sh': 'language-bash',
    '.sql': 'language-sql'
  }
  return syntaxMap[extension.toLowerCase()] || ''
}

// Load file content
const loadFileContent = async () => {
  loading.value = true
  error.value = ''

  try {
    const result = await ReadFileContent(props.filePath)
    fileContent.value = result

    if (result.error) {
      error.value = result.error
    }
  } catch (err) {
    error.value = `Failed to load file content: ${err}`
    console.error('Error loading file content:', err)
  } finally {
    loading.value = false
  }
}

// Go back to file explorer
const goBack = () => {
  emit('back')
}

// Load content when component mounts
onMounted(() => {
  loadFileContent()
})
</script>

<template>
  <div class="file-content">
    <div class="header">
      <button @click="goBack" class="back-btn">
        ← Back to Files
      </button>

      <div v-if="fileContent && !fileContent.error" class="file-info">
        <span class="large-icon">{{ getFileIcon(fileContent) }}</span>
        <div>
          <h2>{{ fileContent.name }}</h2>
          <p class="file-path">{{ fileContent.path }}</p>
          <p class="file-meta">
            {{ formatFileSize(fileContent.size) }} • {{ fileContent.extension || 'No extension' }}
          </p>
        </div>
      </div>
    </div>

    <div class="content">
      <!-- Loading state -->
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>Loading file content...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="error">
        <div class="error-icon">❌</div>
        <h3>Error</h3>
        <p>{{ error }}</p>
        <button @click="loadFileContent" class="retry-btn">Retry</button>
      </div>

      <!-- File content display -->
      <div v-else-if="fileContent" class="content-display">
        <!-- Binary file warning -->
        <div v-if="isBinaryContent(fileContent.content)" class="binary-warning">
          <div class="warning-icon">⚠️</div>
          <h3>Binary File Detected</h3>
          <p>This file contains binary data and cannot be displayed as text.</p>
          <p>File size: {{ formatFileSize(fileContent.size) }}</p>
        </div>

        <!-- Text content -->
        <div v-else class="text-content">
          <div class="content-header">
            <h3>File Content</h3>
            <span class="line-count">
              {{ fileContent.content.split('\n').length }} lines
            </span>
          </div>

          <div class="code-container">
            <pre><code :class="getSyntaxClass(fileContent.extension)">{{ fileContent.content }}</code></pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.file-content {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 2px solid #e0e0e0;
}

.back-btn {
  background: #6c757d;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.back-btn:hover {
  background: #5a6268;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.large-icon {
  font-size: 48px;
  line-height: 1;
}

.file-info h2 {
  margin: 0 0 5px 0;
  color: #333;
  font-size: 24px;
}

.file-path {
  margin: 0 0 5px 0;
  font-size: 14px;
  color: #6c757d;
  word-break: break-all;
}

.file-meta {
  margin: 0;
  font-size: 12px;
  color: #6c757d;
}

.content {
  min-height: 400px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #6c757d;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error {
  text-align: center;
  padding: 40px 20px;
  color: #dc3545;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 8px;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.error h3 {
  margin: 0 0 10px 0;
  color: #721c24;
}

.retry-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  margin-top: 15px;
  transition: background-color 0.3s;
}

.retry-btn:hover {
  background: #c82333;
}

.binary-warning {
  text-align: center;
  padding: 40px 20px;
  color: #856404;
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 8px;
}

.warning-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.binary-warning h3 {
  margin: 0 0 10px 0;
  color: #856404;
}

.content-display {
  background: white;
  border-radius: 8px;
  border: 1px solid #dee2e6;
  overflow: hidden;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #f8f9fa;
  border-bottom: 1px solid #dee2e6;
}

.content-header h3 {
  margin: 0;
  color: #495057;
  font-size: 16px;
}

.line-count {
  font-size: 12px;
  color: #6c757d;
  background: #e9ecef;
  padding: 4px 8px;
  border-radius: 4px;
}

.code-container {
  max-height: 600px;
  overflow: auto;
}

.code-container pre {
  margin: 0;
  padding: 20px;
  background: #f8f9fa;
  font-family: 'Courier New', Consolas, Monaco, 'Andale Mono', monospace;
  font-size: 13px;
  line-height: 1.4;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: #333;
}

.code-container code {
  background: none;
  padding: 0;
  font-size: inherit;
  color: inherit;
}

/* Syntax highlighting placeholders */
.language-javascript { color: #d73a49; }
.language-typescript { color: #d73a49; }
.language-html { color: #22863a; }
.language-css { color: #0d7377; }
.language-json { color: #032f62; }
.language-markdown { color: #032f62; }

@media (max-width: 768px) {
  .file-content {
    padding: 15px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .file-info {
    align-items: flex-start;
  }

  .large-icon {
    font-size: 36px;
  }

  .file-info h2 {
    font-size: 20px;
  }

  .code-container {
    max-height: 400px;
  }

  .code-container pre {
    padding: 15px;
    font-size: 12px;
  }
}
</style>

