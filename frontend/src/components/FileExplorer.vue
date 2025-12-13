<script setup>
import {onMounted, ref} from 'vue'
import {GetFiles, GetFilesInDirectory} from '../../wailsjs/go/wails_app/App'

const emit = defineEmits(['file-selected', 'directory-opened'])

let files = ref([])
let selectedFile = ref(null)
let currentDirectory = ref('')
let loading = ref(false)
let error = ref('')
let navigationHistory = ref([])

// File icon mapping
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

// Load files from backend
const loadFiles = async (directoryPath = null) => {
  loading.value = true
  error.value = ''
  
  try {
    let result
    if (directoryPath) {
      result = await GetFilesInDirectory(directoryPath)
      currentDirectory.value = directoryPath
    } else {
      result = await GetFiles()
      currentDirectory.value = ''
      navigationHistory.value = []
    }
    
    // Add parent directory option if not in root
    if (directoryPath) {
      const parentPath = directoryPath.split('/').slice(0, -1).join('/')
      if (parentPath) {
        files.value = [
          {
            name: '..',
            path: parentPath,
            isDir: true,
            size: 0,
            extension: ''
          },
          ...result.sort((a, b) => {
            // Directories first, then files
            if (a.isDir && !b.isDir) return -1
            if (!a.isDir && b.isDir) return 1
            // Then alphabetically
            return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
          })
        ]
      } else {
        files.value = [
          {
            name: '..',
            path: '',
            isDir: true,
            size: 0,
            extension: ''
          },
          ...result.sort((a, b) => {
            // Directories first, then files
            if (a.isDir && !b.isDir) return -1
            if (!a.isDir && b.isDir) return 1
            // Then alphabetically
            return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
          })
        ]
      }
    } else {
      files.value = result.sort((a, b) => {
        // Directories first, then files
        if (a.isDir && !b.isDir) return -1
        if (!a.isDir && b.isDir) return 1
        // Then alphabetically
        return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
      })
    }
  } catch (err) {
    error.value = `Failed to load files: ${err}`
    console.error('Error loading files:', err)
  } finally {
    loading.value = false
  }
}

// Handle file click
const handleFileClick = (file) => {
  selectedFile.value = file
  
  if (file.isDir) {
    // Navigate to directory
    if (file.name === '..') {
      // Go back to parent directory
      if (navigationHistory.value.length > 0) {
        const previousDir = navigationHistory.value.pop()
        loadFiles(previousDir)
      } else {
        loadFiles() // Back to root
      }
    } else {
      // Add current directory to history and navigate to subdirectory
      if (currentDirectory.value) {
        navigationHistory.value.push(currentDirectory.value)
      }
      loadFiles(file.path)
    }
  } else {
    // Navigate to file content
    emit('file-selected', file.path)
  }
}

// Refresh files
const refreshFiles = () => {
  loadFiles()
  selectedFile.value = null
}

// Load files on component mount
onMounted(() => {
  loadFiles()
})
</script>

<template>
  <div class="file-explorer">
    <div class="header">
      <h2>📂 File Explorer</h2>
      <button @click="refreshFiles" class="refresh-btn" :disabled="loading">
        🔄 Refresh
      </button>
    </div>
    
    <div class="content">
      <!-- Files List -->
      <div class="files-section">
        <h3>Files</h3>
        
        <div v-if="loading" class="loading">
          Loading files...
        </div>
        
        <div v-if="error" class="error">
          ❌ {{ error }}
        </div>
        
        <div v-if="!loading && !error" class="files-grid">
          <div
            v-for="file in files"
            :key="file.path"
            @click="handleFileClick(file)"
            class="file-item"
            :class="{ active: selectedFile?.path === file.path }"
          >
            <div class="file-icon">
              {{ getFileIcon(file) }}
            </div>
            <div class="file-info">
              <div class="file-name" :title="file.name">
                {{ file.name }}
              </div>
              <div class="file-meta">
                <span v-if="!file.isDir">{{ formatFileSize(file.size) }}</span>
                <span v-else>Directory</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.file-explorer {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #e0e0e0;
}

.header h2 {
  margin: 0;
  color: #333;
}

.refresh-btn {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.refresh-btn:hover:not(:disabled) {
  background: #45a049;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.content {
  display: grid;
  grid-template-columns: 1fr; /* was 1fr 1fr */
  gap: 20px;
  height: 600px;
}

.files-section, .details-section {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  border: 1px solid #dee2e6;
}

.files-section h3, .details-section h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #495057;
}

.loading, .error, .no-selection {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}

.error {
  color: #dc3545;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 10px;
  max-height: 500px;
  overflow-y: auto;
  padding: 5px;
}

.file-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.file-item:hover {
  background: #e3f2fd;
  border-color: #2196F3;
  transform: translateY(-2px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.file-item.active {
  background: #bbdefb;
  border-color: #1976D2;
  box-shadow: 0 2px 8px rgba(25, 118, 210, 0.3);
}

.file-icon {
  font-size: 32px;
  margin-bottom: 8px;
  line-height: 1;
}

.file-name {
  font-size: 12px;
  font-weight: 500;
  color: #333;
  word-break: break-word;
  line-height: 1.2;
  margin-bottom: 4px;
}

.file-meta {
  font-size: 10px;
  color: #6c757d;
}

.selected-file-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #dee2e6;
}

.large-icon {
  font-size: 48px;
  margin-right: 15px;
}

.selected-file-header h4 {
  margin: 0 0 5px 0;
  color: #333;
}
.details-black {
  color: #000;
}
.details-black strong {
  color: #000;
}

.selected-file-path {
  margin: 0;
  font-size: 12px;
  color: #6c757d;
  word-break: break-all;
}

.file-attributes {
  margin-bottom: 20px;
}

.attribute {
  margin-bottom: 8px;
  padding: 5px 0;
  border-bottom: 1px solid #f0f0f0;
}

.attribute strong {
  color: #495057;
  margin-right: 8px;
}

.instructions {
  background: #e7f3ff;
  border: 1px solid #b3d9ff;
  border-radius: 4px;
  padding: 15px;
  margin-top: 20px;
}

.instructions h5 {
  margin: 0 0 10px 0;
  color: #0066cc;
  font-size: 14px;
}

.instructions p {
  margin: 0;
  color: #004080;
  font-size: 13px;
}

.backend-response {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  padding: 15px;
}

.backend-response h5 {
  margin-top: 0;
  margin-bottom: 10px;
  color: #495057;
}

.backend-response pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #333;
  background: #f8f9fa;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

@media (max-width: 768px) {
  .content {
    grid-template-columns: 1fr;
    height: auto;
  }
  
  .files-grid {
    max-height: 300px;
  }
  
  .selected-file-header {
    flex-direction: column;
    text-align: center;
  }
  
  .large-icon {
    margin-right: 0;
    margin-bottom: 10px;
  }
}
</style>
