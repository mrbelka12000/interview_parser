<template>
  <div class="interviews-container">
    <div class="interviews-header">
      <h2>Interview Management</h2>
      <div class="filters">
        <div class="filter-group">
          <label>Date From:</label>
          <input type="date" v-model="filters.dateFrom" @change="refreshInterviews" />
        </div>
        <div class="filter-group">
          <label>Date To:</label>
          <input type="date" v-model="filters.dateTo" @change="refreshInterviews" />
        </div>
        <button @click="refreshInterviews" class="btn-primary">Refresh</button>
        <button @click="showCreateModal = true" class="btn-success">New Interview</button>
      </div>
    </div>

    <!-- Interviews List -->
    <div class="interviews-section">
      <h3>All Interviews</h3>
      <div class="interview-list" v-if="interviews.length > 0">
        <div 
          v-for="interview in interviews" 
          :key="interview.id" 
          class="interview-card"
        >
          <div class="interview-header">
            <h4>Interview #{{ interview.id }}</h4>
            <div class="interview-actions">
              <button @click="viewInterview(interview)" class="btn-secondary">View</button>
              <button @click="editInterview(interview)" class="btn-warning">Edit</button>
              <button @click="deleteInterview(interview.id)" class="btn-danger">Delete</button>
            </div>
          </div>
          
          <div class="interview-info">
            <div class="info-row">
              <span class="label">Questions:</span>
              <span class="value">{{ interview.qa.length }}</span>
            </div>
            <div class="info-row">
              <span class="label">Created:</span>
              <span class="value">{{ formatDate(interview.created_at) }}</span>
            </div>
            <div class="info-row">
              <span class="label">Last Updated:</span>
              <span class="value">{{ formatDate(interview.updated_at) }}</span>
            </div>
            <div class="info-row">
              <span class="label">Average Accuracy:</span>
              <span class="value">{{ getAverageAccuracy(interview.qa) }}%</span>
            </div>
            <div class="info-row">
              <span class="label">Answered:</span>
              <span class="value">{{ getAnsweredCount(interview.qa) }}/{{ interview.qa.length }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="no-data">
        <p>No interviews found. Create your first interview or process some audio files to get started.</p>
      </div>
    </div>

    <!-- View Interview Modal -->
    <div v-if="selectedInterview" class="modal-overlay" @click="closeViewModal">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <h3>Interview #{{ selectedInterview.id }} Details</h3>
          <button @click="closeViewModal" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="interview-meta">
            <div class="meta-item">
              <label>Created:</label>
              <span>{{ formatDate(selectedInterview.created_at) }}</span>
            </div>
            <div class="meta-item">
              <label>Last Updated:</label>
              <span>{{ formatDate(selectedInterview.updated_at) }}</span>
            </div>
            <div class="meta-item">
              <label>Total Questions:</label>
              <span>{{ selectedInterview.qa.length }}</span>
            </div>
            <div class="meta-item">
              <label>Average Accuracy:</label>
              <span>{{ getAverageAccuracy(selectedInterview.qa) }}%</span>
            </div>
          </div>
          
          <h4>Questions & Answers</h4>
          <div class="qa-list">
            <div v-for="(qa, index) in selectedInterview.qa" :key="qa.id" class="qa-item">
              <div class="qa-header">
                <span class="question-number">Q{{ index + 1 }}</span>
                <span class="accuracy-badge" :class="getAccuracyClass(qa.accuracy)">
                  {{ qa.accuracy.toFixed(1) }}%
                </span>
              </div>
              <div class="question-text">
                <strong>Question:</strong> {{ qa.question }}
              </div>
              <div class="answer-text">
                <strong>Answer:</strong> {{ qa.full_answer || 'No answer provided' }}
              </div>
              <div v-if="qa.reason_unanswered" class="reason-text">
                <strong>Reason Unanswered:</strong> {{ qa.reason_unanswered }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Interview Modal -->
    <div v-if="showCreateModal || editingInterview" class="modal-overlay" @click="closeCreateModal">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <h3>{{ editingInterview ? 'Edit Interview' : 'Create New Interview' }}</h3>
          <button @click="closeCreateModal" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveInterview">
            <div class="form-group">
              <label>Interview ID (edit only):</label>
              <input type="number" v-model="formData.id" :disabled="!editingInterview" />
            </div>
            
            <div class="qa-form-section">
              <h4>Questions & Answers</h4>
              <div v-for="(qa, index) in formData.qa" :key="index" class="qa-form-item">
                <div class="qa-form-header">
                  <span>Question {{ index + 1 }}</span>
                  <button type="button" @click="removeQuestion(index)" class="btn-danger small">Remove</button>
                </div>
                <div class="form-group">
                  <label>Question:</label>
                  <textarea v-model="qa.question" required rows="2"></textarea>
                </div>
                <div class="form-group">
                  <label>Answer:</label>
                  <textarea v-model="qa.full_answer" rows="3"></textarea>
                </div>
                <div class="form-group">
                  <label>Accuracy (0-100):</label>
                  <input type="number" v-model="qa.accuracy" min="0" max="100" step="0.1" required />
                </div>
                <div class="form-group">
                  <label>Reason Unanswered:</label>
                  <input type="text" v-model="qa.reason_unanswered" />
                </div>
              </div>
              
              <button type="button" @click="addQuestion" class="btn-secondary">Add Question</button>
            </div>
            
            <div class="form-actions">
              <button type="submit" class="btn-primary">
                {{ editingInterview ? 'Update Interview' : 'Create Interview' }}
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
import { GetAllInterviewsAPI, GetInterviewAPI, SaveInterviewAPI, UpdateInterviewAPI, DeleteInterviewAPI } from '../../wailsjs/go/app/App'

export default {
  name: 'Interviews',
  data() {
    return {
      interviews: [],
      selectedInterview: null,
      editingInterview: null,
      showCreateModal: false,
      filters: {
        dateFrom: '',
        dateTo: ''
      },
      formData: {
        id: null,
        qa: []
      },
      loading: false
    }
  },
  mounted() {
    this.refreshInterviews()
  },
  methods: {
    async refreshInterviews() {
      this.loading = true
      try {
        const result = await GetAllInterviewsAPI(this.filters.dateFrom, this.filters.dateTo)
        // Ensure result is an array, fallback to empty array if null/undefined
        this.interviews = Array.isArray(result) ? result : []
      } catch (error) {
        console.error('Error fetching interviews:', error)
        this.interviews = [] // Ensure interviews is an array even on error
        this.$emit('error', 'Failed to load interviews')
      } finally {
        this.loading = false
      }
    },

    async viewInterview(interview) {
      try {
        this.selectedInterview = await GetInterviewAPI(interview.id)
      } catch (error) {
        console.error('Error fetching interview details:', error)
        this.$emit('error', 'Failed to load interview details')
      }
    },
    
    editInterview(interview) {
      this.editingInterview = interview
      this.formData = {
        id: interview.id,
        qa: JSON.parse(JSON.stringify(interview.qa))
      }
    },
    
    async deleteInterview(id) {
      if (!confirm('Are you sure you want to delete this interview?')) {
        return
      }
      
      try {
        await DeleteInterviewAPI(id)
        this.refreshInterviews()
      } catch (error) {
        console.error('Error deleting interview:', error)
        this.$emit('error', 'Failed to delete interview')
      }
    },
    
    async saveInterview() {
      try {
        if (this.editingInterview) {
          const interviewData = { id: this.formData.id }
          await UpdateInterviewAPI(interviewData, this.formData.qa)
        } else {
          const interviewData = { qa: this.formData.qa }
          await SaveInterviewAPI(interviewData)
        }
        
        this.closeCreateModal()
        this.refreshInterviews()
      } catch (error) {
        console.error('Error saving interview:', error)
        this.$emit('error', 'Failed to save interview')
      }
    },
    
    closeViewModal() {
      this.selectedInterview = null
    },
    
    closeCreateModal() {
      this.showCreateModal = false
      this.editingInterview = null
      this.formData = {
        id: null,
        qa: []
      }
    },
    
    addQuestion() {
      this.formData.qa.push({
        question: '',
        full_answer: '',
        accuracy: 0,
        reason_unanswered: ''
      })
    },
    
    removeQuestion(index) {
      this.formData.qa.splice(index, 1)
    },
    
    getAverageAccuracy(qa) {
      if (!qa || qa.length === 0) return 0
      const sum = qa.reduce((acc, item) => acc + item.accuracy, 0)
      return (sum / qa.length).toFixed(1)
    },
    
    getAnsweredCount(qa) {
      if (!qa) return 0
      return qa.filter(item => item.full_answer && item.full_answer.trim() !== '').length
    },
    
    getAccuracyClass(accuracy) {
      if (accuracy >= 0.9) return 'high'
      if (accuracy >= 0.7) return 'medium'
      return 'low'
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
.interviews-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.interviews-header {
  margin-bottom: 30px;
}

.interviews-header h2 {
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

.btn-danger.small {
  padding: 4px 8px;
  font-size: 11px;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.interviews-section {
  margin-bottom: 40px;
}

.interviews-section h3 {
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.interview-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.interview-card {
  background: white;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.interview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.interview-header h4 {
  margin: 0;
  color: #333;
}

.interview-actions {
  display: flex;
  gap: 8px;
}

.interview-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 10px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.info-row .label {
  color: #666;
  font-weight: 500;
}

.info-row .value {
  color: #333;
  font-weight: 600;
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

.interview-meta {
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

.qa-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.qa-item {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
}

.qa-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.question-number {
  font-weight: bold;
  color: #007bff;
}

.accuracy-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  color: white;
}

.accuracy-badge.high {
  background-color: #28a745;
}

.accuracy-badge.medium {
  background-color: #ffc107;
}

.accuracy-badge.low {
  background-color: #dc3545;
}

.question-text, .answer-text, .reason-text {
  margin-bottom: 10px;
  font-size: 14px;
  line-height: 1.5;
}

/* Form Styles */
.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 600;
  color: #666;
}

.form-group input, .form-group textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.qa-form-section {
  margin-bottom: 20px;
}

.qa-form-item {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
}

.qa-form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-weight: bold;
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
  
  .interview-info {
    grid-template-columns: 1fr;
  }
  
  .interview-actions {
    flex-direction: column;
    gap: 5px;
  }
  
  .modal-content {
    width: 95%;
  }
}

.interviews-container,
.interviews-container * {
  color: black !important;
}
</style>
