<template>
  <div class="api-key-manager">
    <div class="api-key-container">
      <h2>🔑 OpenAI API Key Management</h2>
      <p class="description">
        Manage your OpenAI API key for transcription and analysis services.
        The key is securely stored in a local database and encrypted.
      </p>

      <div class="current-key-section">
        <h3>Current API Key</h3>
        <div v-if="loading" class="loading">
          <div class="loading-spinner"></div>
          <p>Loading current API key...</p>
        </div>
        
        <div v-else-if="currentKey" class="key-display">
          <div class="key-info">
            <span class="key-label">Status:</span>
            <span class="key-status valid">✅ Valid</span>
          </div>
          <div class="key-info">
            <span class="key-label">Key:</span>
            <span class="key-value">{{ maskApiKey(currentKey) }}</span>
          </div>
          <div class="key-info">
            <span class="key-label">Last Updated:</span>
            <span class="key-date">{{ lastUpdated || 'Unknown' }}</span>
          </div>
        </div>
        
        <div v-else class="no-key">
          <div class="no-key-icon">🔒</div>
          <p>No API key configured</p>
          <p class="no-key-desc">Please add an OpenAI API key to use transcription services</p>
        </div>
      </div>

      <div class="new-key-section">
        <h3>{{ currentKey ? 'Update API Key' : 'Add New API Key' }}</h3>
        <div class="key-input-group">
          <label for="apiKey">OpenAI API Key:</label>
          <div class="input-wrapper">
            <input
              id="apiKey"
              v-model="newApiKey"
              :type="showKey ? 'text' : 'password'"
              placeholder="sk-..."
              class="api-key-input"
              :class="{ 'error': error }"
            />
            <button 
              @click="showKey = !showKey" 
              class="toggle-visibility"
              :title="showKey ? 'Hide key' : 'Show key'"
            >
              {{ showKey ? '👁️‍🗨️' : '👁️' }}
            </button>
          </div>
          <div class="input-help">
            Your API key starts with 'sk-' and can be found in your OpenAI dashboard
          </div>
        </div>

        <div v-if="error" class="error-message">
          <div class="error-icon">❌</div>
          <span>{{ error }}</span>
        </div>

        <div v-if="success" class="success-message">
          <div class="success-icon">✅</div>
          <span>{{ success }}</span>
        </div>

        <div class="action-buttons">
          <button
            @click="saveApiKey"
            :disabled="!newApiKey.trim() || isSaving"
            class="save-button"
          >
            <span v-if="!isSaving">{{ currentKey ? 'Update Key' : 'Save Key' }}</span>
            <span v-else>💾 Saving...</span>
          </button>
          
          <button
            v-if="currentKey"
            @click="confirmDelete = true"
            class="delete-button"
            :disabled="isDeleting"
          >
            <span v-if="!isDeleting">🗑️ Remove Key</span>
            <span v-else>🗑️ Removing...</span>
          </button>
        </div>
      </div>

      <!-- Delete Confirmation Modal -->
      <div v-if="confirmDelete" class="modal-overlay" @click="confirmDelete = false">
        <div class="modal" @click.stop>
          <h3>⚠️ Confirm Deletion</h3>
          <p>Are you sure you want to remove the OpenAI API key?</p>
          <p class="warning">This will disable transcription services until a new key is added.</p>
          <div class="modal-actions">
            <button @click="confirmDelete = false" class="cancel-button">Cancel</button>
            <button @click="deleteApiKey" class="confirm-delete-button">Delete Key</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetOpenAIAPIKey, SaveOpenAIAPIKey, DeleteOpenAIAPIKey } from '../../wailsjs/go/wails_app/App'

const loading = ref(true)
const currentKey = ref('')
const lastUpdated = ref('')
const newApiKey = ref('')
const showKey = ref(false)
const isSaving = ref(false)
const isDeleting = ref(false)
const error = ref('')
const success = ref('')
const confirmDelete = ref(false)

// Load current API key on component mount
const loadCurrentKey = async () => {
  try {
    loading.value = true
    error.value = ''
    success.value = ''
    
    const result = await GetOpenAIAPIKey()
    if (result.success) {
      currentKey.value = result.apiKey || ''
      lastUpdated.value = result.lastUpdated || new Date().toLocaleString()
    } else {
      // No key found is not an error
      currentKey.value = ''
      lastUpdated.value = ''
    }
  } catch (err) {
    console.error('Error loading API key:', err)
    error.value = `Failed to load API key: ${err.message || err.toString()}`
  } finally {
    loading.value = false
  }
}

// Save new API key
const saveApiKey = async () => {
  if (!newApiKey.value.trim()) {
    error.value = 'API key cannot be empty'
    return
  }

  // Basic validation for OpenAI API key format
  if (!newApiKey.value.startsWith('sk-')) {
    error.value = 'Invalid API key format. OpenAI API keys start with "sk-"'
    return
  }

  try {
    isSaving.value = true
    error.value = ''
    success.value = ''
    
    const result = await SaveOpenAIAPIKey(newApiKey.value.trim())
    
    if (result.success) {
      success.value = 'API key saved successfully!'
      currentKey.value = newApiKey.value.trim()
      lastUpdated.value = new Date().toLocaleString()
      newApiKey.value = ''
      showKey.value = false
      
      // Clear success message after 3 seconds
      setTimeout(() => {
        success.value = ''
      }, 3000)
    } else {
      error.value = result.message || 'Failed to save API key'
    }
  } catch (err) {
    console.error('Error saving API key:', err)
    error.value = `Failed to save API key: ${err.message || err.toString()}`
  } finally {
    isSaving.value = false
  }
}

// Delete API key
const deleteApiKey = async () => {
  try {
    isDeleting.value = true
    error.value = ''
    confirmDelete.value = false
    
    const result = await DeleteOpenAIAPIKey()
    
    if (result.success) {
      success.value = 'API key removed successfully!'
      currentKey.value = ''
      lastUpdated.value = ''
      
      // Clear success message after 3 seconds
      setTimeout(() => {
        success.value = ''
      }, 3000)
    } else {
      error.value = result.message || 'Failed to remove API key'
    }
  } catch (err) {
    console.error('Error deleting API key:', err)
    error.value = `Failed to remove API key: ${err.message || err.toString()}`
  } finally {
    isDeleting.value = false
  }
}

// Mask API key for display
const maskApiKey = (key) => {
  if (!key || key.length < 8) return key
  return key.substring(0, 7) + '*'.repeat(Math.max(0, key.length - 7))
}

// Load key on component mount
onMounted(() => {
  loadCurrentKey()
})
</script>

<style scoped>
.api-key-manager {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.api-key-container {
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

.current-key-section {
  margin-bottom: 30px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.current-key-section h3 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 18px;
}

.loading {
  display: flex;
  align-items: center;
  gap: 15px;
  color: #6c757d;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.key-display {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.key-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.key-label {
  font-weight: 600;
  color: #495057;
  min-width: 120px;
}

.key-status.valid {
  color: #28a745;
  font-weight: 500;
}

.key-value {
  font-family: 'Courier New', monospace;
  background: #e9ecef;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 450px; /* adjust if you want */
  display: inline-block;
}

.key-date {
  color: #6c757d;
  font-size: 14px;
}

.no-key {
  text-align: center;
  padding: 20px;
  color: #6c757d;
}

.no-key-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.no-key p {
  margin: 5px 0;
}

.no-key-desc {
  font-size: 14px;
  color: #868e96;
}

.new-key-section h3 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 18px;
}

.key-input-group {
  margin-bottom: 20px;
}

.key-input-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #495057;
}

.input-wrapper {
  display: flex;
  gap: 10px;
  margin-bottom: 8px;
}

.api-key-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #dee2e6;
  border-radius: 6px;
  font-size: 14px;
  font-family: 'Courier New', monospace;
  transition: border-color 0.3s;
}

.api-key-input:focus {
  outline: none;
  border-color: #007bff;
}

.api-key-input.error {
  border-color: #dc3545;
}

.toggle-visibility {
  padding: 8px 12px;
  background: #f8f9fa;
  border: 2px solid #dee2e6;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.toggle-visibility:hover {
  background: #e9ecef;
}

.input-help {
  font-size: 12px;
  color: #6c757d;
  margin-top: 5px;
}

.error-message, .success-message {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 6px;
  margin-bottom: 20px;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.success-message {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.action-buttons {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.save-button, .delete-button {
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.save-button {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
}

.save-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(40, 167, 69, 0.3);
}

.save-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.delete-button {
  background: linear-gradient(135deg, #dc3545 0%, #c82333 100%);
  color: white;
}

.delete-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(220, 53, 69, 0.3);
}

.delete-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Modal Styles */
.modal-overlay {
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

.modal {
  background: white;
  border-radius: 12px;
  padding: 30px;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
}

.modal h3 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 18px;
}

.modal p {
  margin: 0 0 10px 0;
  color: #666;
  line-height: 1.5;
}

.modal .warning {
  color: #dc3545;
  font-weight: 500;
}

.modal-actions {
  display: flex;
  gap: 15px;
  margin-top: 25px;
  justify-content: flex-end;
}

.cancel-button, .confirm-delete-button {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.cancel-button {
  background: #6c757d;
  color: white;
}

.cancel-button:hover {
  background: #5a6268;
}

.confirm-delete-button {
  background: #dc3545;
  color: white;
}

.confirm-delete-button:hover {
  background: #c82333;
}

@media (max-width: 600px) {
  .api-key-manager {
    padding: 10px;
  }

  .api-key-container {
    padding: 20px;
  }

  .input-wrapper {
    flex-direction: column;
    gap: 10px;
  }

  .action-buttons {
    flex-direction: column;
  }

  .modal-actions {
    flex-direction: column;
  }
}

.api-key-manager {
  color: #000;
}
</style>
