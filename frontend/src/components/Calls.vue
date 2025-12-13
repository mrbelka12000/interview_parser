<template>
  <div class="calls-container">
    <div class="calls-header">
      <h2>Call Management</h2>
      <div class="filters">
        <div class="filter-group">
          <label>Date From:</label>
          <input type="date" v-model="filters.dateFrom" @change="refreshCalls" />
        </div>
        <div class="filter-group">
          <label>Date To:</label>
          <input type="date" v-model="filters.dateTo" @change="refreshCalls" />
        </div>
        <button @click="refreshCalls" class="btn-primary">Refresh</button>
        <button @click="showCreateModal = true" class="btn-success">New Call</button>
      </div>
    </div>

    <!-- Calls List -->
    <div class="calls-section">
      <h3>All Calls</h3>
      <div class="call-list" v-if="calls.length > 0">
        <div 
          v-for="call in calls" 
          :key="call.ID"
          class="call-card"
        >
          <div class="call-header">
            <h4>Call #{{ call.ID }}</h4>
            <div class="call-actions">
              <button @click="viewCall(call)" class="btn-secondary">View</button>
              <button @click="editCall(call)" class="btn-warning">Edit</button>
              <button @click="deleteCall(call.ID)" class="btn-danger">Delete</button>
            </div>
          </div>
          
          <div class="call-info">
            <div class="info-row">
              <span class="label">Created:</span>
              <span class="value">{{ formatDate(call.created_at) }}</span>
            </div>
            <div class="info-row">
              <span class="label">Last Updated:</span>
              <span class="value">{{ formatDate(call.updated_at) }}</span>
            </div>
            <div class="info-row">
              <span class="label">Transcript:</span>
              <span class="value">{{ getTranscriptPreview(call.transcript) }}</span>
            </div>
            <div class="info-row">
              <span class="label">Has Analysis:</span>
              <span class="value">{{ call.analysis ? 'Yes' : 'No' }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="no-data">
        <p>No calls found. Create your first call or process some audio files to get started.</p>
      </div>
    </div>

    <!-- View Call Modal -->
    <div v-if="selectedCall" class="modal-overlay" @click="closeViewModal">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <h3>Call #{{ selectedCall.id }} Details</h3>
          <button @click="closeViewModal" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="call-meta">
            <div class="meta-item">
              <label>Created:</label>
              <span>{{ formatDate(selectedCall.created_at) }}</span>
            </div>
            <div class="meta-item">
              <label>Last Updated:</label>
              <span>{{ formatDate(selectedCall.updated_at) }}</span>
            </div>
            <div class="meta-item">
              <label>Has Analysis:</label>
              <span>{{ selectedCall.analysis ? 'Yes' : 'No' }}</span>
            </div>
          </div>
          
          <h4>Transcript</h4>
          <div class="transcript-content">
            {{ selectedCall.transcript }}
          </div>
          
          <h4 v-if="selectedCall.analysis">Analysis</h4>
          <div v-if="selectedCall.analysis" class="analysis-content">
            <pre>{{ formatAnalysis(selectedCall.analysis) }}</pre>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Call Modal -->
    <div v-if="showCreateModal || editingCall" class="modal-overlay" @click="closeCreateModal">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <h3>{{ editingCall ? 'Edit Call' : 'Create New Call' }}</h3>
          <button @click="closeCreateModal" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveCall">
            <div class="form-group">
              <label>Call ID (edit only):</label>
              <input type="number" v-model="formData.id" :disabled="!editingCall" />
            </div>
            
            <div class="form-group">
              <label>Transcript:</label>
              <textarea v-model="formData.transcript" required rows="10" placeholder="Enter the call transcript here..."></textarea>
            </div>
            
            <div class="form-group">
              <label>Analysis (JSON format, optional):</label>
              <textarea v-model="formData.analysisText" rows="8" placeholder="Enter analysis in JSON format or leave empty..."></textarea>
            </div>
            
            <div class="form-actions">
              <button type="submit" class="btn-primary">
                {{ editingCall ? 'Update Call' : 'Create Call' }}
              </button>
              <button type="button" @click="closeCreateModal" class="btn-secondary">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { GetAllCallsAPI, GetCallsByDateRangeAPI, GetCallAPI, SaveCallAPI, UpdateCallAPI, DeleteCallAPI } from '../../wailsjs/go/wails_app/App'

export default {
  name: 'Calls',
  data() {
    return {
      calls: [],
      selectedCall: null,
      editingCall: null,
      showCreateModal: false,
      filters: {
        dateFrom: '',
        dateTo: ''
      },
      formData: {
        id: null,
        transcript: '',
        analysisText: ''
      },
      loading: false
    }
  },
  mounted() {
    this.refreshCalls()
  },
  methods: {
    async refreshCalls() {
      this.loading = true
      try {
        let result
        if (this.filters.dateFrom || this.filters.dateTo) {
          result = await GetCallsByDateRangeAPI(this.filters.dateFrom, this.filters.dateTo)
        } else {
          result = await GetAllCallsAPI(100, 0) // Default limit and offset
        }
        // Ensure result is an array, fallback to empty array if null/undefined
        this.calls = Array.isArray(result) ? result : []
      } catch (error) {
        console.error('Error fetching calls:', error)
        this.calls = [] // Ensure calls is an array even on error
        this.$emit('error', 'Failed to load calls')
      } finally {
        this.loading = false
      }
    },
    
    async viewCall(call) {
      try {
        this.selectedCall = await GetCallAPI(call.ID)
      } catch (error) {
        console.error('Error fetching call details:', error)
        this.$emit('error', 'Failed to load call details')
      }
    },
    
    editCall(call) {
      this.editingCall = call
      this.formData = {
        id: call.ID,
        transcript: call.transcript,
        analysisText: call.analysis ? JSON.stringify(JSON.parse(call.analysis), null, 2) : ''
      }
    },
    
    async deleteCall(id) {
      if (!confirm('Are you sure you want to delete this call?')) {
        return
      }
      
      try {
        await DeleteCallAPI(id)
        this.refreshCalls()
      } catch (error) {
        console.error('Error deleting call:', error)
        this.$emit('error', 'Failed to delete call')
      }
    },
    
    async saveCall() {
      try {
        let analysis = null
        if (this.formData.analysisText.trim()) {
          try {
            analysis = JSON.parse(this.formData.analysisText)
          } catch (e) {
            this.$emit('error', 'Invalid JSON format in analysis field')
            return
          }
        }

        if (this.editingCall) {
          await UpdateCallAPI(this.formData.id, this.formData.transcript, analysis)
        } else {
          const callData = {
            id: 0, // Will be set by database
            transcript: this.formData.transcript,
            analysis: analysis ? JSON.stringify(analysis) : null
          }
          const createdCall = await SaveCallAPI(callData)
          // Add the created call to the local list to avoid full refresh
          if (createdCall && createdCall.id) {
            this.calls.unshift(createdCall)
          }
        }
        
        this.closeCreateModal()
        if (!this.editingCall) {
          // Only refresh if we didn't successfully add the call locally
          this.refreshCalls()
        }
      } catch (error) {
        console.error('Error saving call:', error)
        this.$emit('error', 'Failed to save call')
      }
    },
    
    closeViewModal() {
      this.selectedCall = null
    },
    
    closeCreateModal() {
      this.showCreateModal = false
      this.editingCall = null
      this.formData = {
        id: null,
        transcript: '',
        analysisText: ''
      }
    },
    
    getTranscriptPreview(transcript) {
      if (!transcript) return 'No transcript'
      const preview = transcript.substring(0, 100)
      return preview.length < transcript.length ? preview + '...' : preview
    },
    
    formatAnalysis(analysis) {
      try {
        return JSON.stringify(JSON.parse(analysis), null, 2)
      } catch (e) {
        return analysis
      }
    },
    
    formatDate(dateString) {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleDateString('en-GB', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      })
    }
  }
}
</script>

<style scoped>
.calls-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.calls-header {
  margin-bottom: 30px;
}

.calls-header h2 {
  margin-bottom: 20px;
  color: #333;
}

.filters {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  align-items: end;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.filter-group label {
  font-size: 12px;
  font-weight: 600;
  color: #666;
}

.filter-group input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.btn-primary {
  background-color: #007bff;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.btn-success {
  background-color: #28a745;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-warning {
  background-color: #ffc107;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.btn-danger {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.calls-section {
  margin-bottom: 40px;
}

.calls-section h3 {
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.call-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.call-card {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.call-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.call-header h4 {
  margin: 0;
  color: #333;
}

.call-actions {
  display: flex;
  gap: 8px;
}

.call-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 10px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  align-items: start;
}

.info-row .label {
  color: #666;
  font-weight: 500;
  flex-shrink: 0;
  margin-right: 10px;
}

.info-row .value {
  color: #333;
  font-weight: 600;
  text-align: right;
  word-break: break-word;
}

.no-data {
  text-align: center;
  padding: 40px;
  color: #666;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 700px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content.large {
  max-width: 900px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.modal-body {
  padding: 20px;
}

.call-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
  margin-bottom: 30px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
}

.meta-item label {
  font-weight: 600;
  color: #666;
}

.transcript-content, .analysis-content {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 20px;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  max-height: 300px;
  overflow-y: auto;
}

.analysis-content pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

/* Form Styles */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #666;
}

.form-group input, .form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  font-family: inherit;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .call-info {
    grid-template-columns: 1fr;
  }
  
  .call-actions {
    flex-direction: column;
    gap: 5px;
  }
  
  .modal-content {
    width: 95%;
  }
  
  .info-row {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .info-row .value {
    text-align: left;
    margin-top: 5px;
  }
}

.calls-container,
.calls-container * {
  color: black !important;
}
</style>
