<template>
  <div class="analytics-container">
    <div class="analytics-header">
      <h2>Interview Analytics</h2>
      <div class="filters">
        <div class="filter-group">
          <label>Date From:</label>
          <input type="date" v-model="filters.dateFrom" @change="refreshAnalytics" />
        </div>
        <div class="filter-group">
          <label>Date To:</label>
          <input type="date" v-model="filters.dateTo" @change="refreshAnalytics" />
        </div>
        <div class="filter-group">
          <label>Min Accuracy:</label>
          <input type="number" v-model="filters.minAccuracy" min="0" max="1" step="0.1" @change="refreshAnalytics" />
        </div>
        <div class="filter-group">
          <label>Max Accuracy:</label>
          <input type="number" v-model="filters.maxAccuracy" min="0" max="1" step="0.1" @change="refreshAnalytics" />
        </div>
        <button @click="refreshAnalytics" class="btn-primary">Refresh</button>
      </div>
    </div>

    <!-- Global Analytics -->
    <div class="analytics-section" v-if="globalAnalytics">
      <h3>Global Statistics</h3>
      <div class="stats-grid">
        <div class="stat-card">
          <h4>Total Interviews</h4>
          <span class="stat-value">{{ globalAnalytics.totalInterviews }}</span>
        </div>
        <div class="stat-card">
          <h4>Total Questions</h4>
          <span class="stat-value">{{ globalAnalytics.totalQuestions }}</span>
        </div>
        <div class="stat-card">
          <h4>Answered Rate</h4>
          <span class="stat-value">{{ globalAnalytics.globalAnsweredPercent.toFixed(1) }}%</span>
        </div>
        <div class="stat-card">
          <h4>Average Accuracy</h4>
          <span class="stat-value">{{ (globalAnalytics.globalAverageAccuracy * 100).toFixed(1) }}%</span>
        </div>
      </div>

      <div class="best-worst-section">
        <div class="interview-card best">
          <h4>Best Interview</h4>
          <p>{{ globalAnalytics.bestInterviewPath }}</p>
          <span>Score: {{ globalAnalytics.bestInterviewScore.toFixed(1) }}</span>
        </div>
        <div class="interview-card worst">
          <h4>Worst Interview</h4>
          <p>{{ globalAnalytics.worstInterviewPath }}</p>
          <span>Score: {{ globalAnalytics.worstInterviewScore.toFixed(1) }}</span>
        </div>
      </div>
    </div>

    <!-- Individual Interview Analytics -->
    <div class="analytics-section">
      <h3>Individual Interview Analytics</h3>
      <div class="interview-list" v-if="interviewAnalytics.length > 0">
        <div 
          v-for="analytics in interviewAnalytics" 
          :key="analytics.id" 
          class="interview-analytics-card"
        >
          <div class="interview-header">
            <button @click="showDetails(analytics)" class="btn-secondary">Details</button>
          </div>
          
          <div class="interview-stats">
            <div class="progress-section">
              <div class="progress-item">
                <span>Answered: {{ analytics.answeredQuestions }}/{{ analytics.totalQuestions }}</span>
                <div class="progress-bar">
                  <div 
                    class="progress-fill answered" 
                    :style="{ width: analytics.answeredPercentage + '%' }"
                  ></div>
                </div>
              </div>
              <div class="progress-item">
                <span>Unanswered: {{ analytics.unansweredQuestions }}/{{ analytics.totalQuestions }}</span>
                <div class="progress-bar">
                  <div 
                    class="progress-fill unanswered" 
                    :style="{ width: analytics.unansweredPercentage + '%' }"
                  ></div>
                </div>
              </div>
            </div>

            <div class="accuracy-section">
              <div class="accuracy-item">
                <span>Avg Accuracy: {{ (analytics.averageAccuracy * 100).toFixed(1) }}%</span>
              </div>
              <div class="confidence-breakdown">
                <span class="high-conf">High: {{ analytics.highConfidenceQuestions }}</span>
                <span class="med-conf">Med: {{ analytics.mediumConfidenceQuestions }}</span>
                <span class="low-conf">Low: {{ analytics.lowConfidenceQuestions }}</span>
              </div>
            </div>

            <div class="date-info">
              <small>Created: {{ formatDate(analytics.createdAt) }}</small>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="no-data">
        <p>No interview analytics found. Process some interviews to see analytics here.</p>
      </div>
    </div>

    <!-- Detail Modal -->
    <div v-if="selectedAnalytics" class="modal-overlay" @click="closeDetails">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Interview Details</h3>
          <button @click="closeDetails" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="detail-grid">
            <div class="detail-item">
              <label>Interview Path:</label>
              <span>{{ selectedAnalytics.interviewPath }}</span>
            </div>
            <div class="detail-item">
              <label>Total Questions:</label>
              <span>{{ selectedAnalytics.totalQuestions }}</span>
            </div>
            <div class="detail-item">
              <label>Answered Questions:</label>
              <span>{{ selectedAnalytics.answeredQuestions }} ({{ selectedAnalytics.answeredPercentage.toFixed(1) }}%)</span>
            </div>
            <div class="detail-item">
              <label>Unanswered Questions:</label>
              <span>{{ selectedAnalytics.unansweredQuestions }} ({{ selectedAnalytics.unansweredPercentage.toFixed(1) }}%)</span>
            </div>
            <div class="detail-item">
              <label>Average Accuracy:</label>
              <span>{{ (selectedAnalytics.averageAccuracy * 100).toFixed(1) }}%</span>
            </div>
            <div class="detail-item">
              <label>Average Answered Accuracy:</label>
              <span>{{ selectedAnalytics.averageAnsweredAccuracy ? (selectedAnalytics.averageAnsweredAccuracy * 100).toFixed(1) + '%' : 'N/A' }}</span>
            </div>
            <div class="detail-item">
              <label>Confidence Breakdown:</label>
              <div class="confidence-detail">
                <span class="high-conf">High: {{ selectedAnalytics.highConfidenceQuestions }}</span>
                <span class="med-conf">Medium: {{ selectedAnalytics.mediumConfidenceQuestions }}</span>
                <span class="low-conf">Low: {{ selectedAnalytics.lowConfidenceQuestions }}</span>
              </div>
            </div>
            <div class="detail-item">
              <label>Questions with Reason:</label>
              <span>{{ selectedAnalytics.questionsWithReason }}</span>
            </div>
            <div class="detail-item">
              <label>Analysis Path:</label>
              <span>{{ selectedAnalytics.analysisPath }}</span>
            </div>
            <div class="detail-item">
              <label>Created:</label>
              <span>{{ formatDate(selectedAnalytics.createdAt) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { GetGlobalAnalyticsAPI, GetAllInterviewAnalyticsAPI } from '../../wailsjs/go/app/App'

export default {
  name: 'Analytics',
  data() {
    return {
      globalAnalytics: null,
      interviewAnalytics: [],
      selectedAnalytics: null,
      filters: {
        dateFrom: '',
        dateTo: '',
        minAccuracy: 0,
        maxAccuracy: 0
      },
      loading: false
    }
  },
  mounted() {
    this.refreshAnalytics()
  },
  methods: {
    async refreshAnalytics() {
      this.loading = true
      try {
        // Get global analytics
        let globalResult = await GetGlobalAnalyticsAPI(
          this.filters.dateFrom,
          this.filters.dateTo,
          this.filters.minAccuracy,
          this.filters.maxAccuracy
        )
        this.globalAnalytics = globalResult || null

        // Get all interview analytics
        let interviewResult = await GetAllInterviewAnalyticsAPI(
          this.filters.dateFrom,
          this.filters.dateTo,
          this.filters.minAccuracy,
          this.filters.maxAccuracy
        )
        // Ensure result is an array, fallback to empty array if null/undefined
        this.interviewAnalytics = Array.isArray(interviewResult) ? interviewResult : []
      } catch (error) {
        console.error('Error fetching analytics:', error)
        this.globalAnalytics = null
        this.interviewAnalytics = []
        this.$emit('error', 'Failed to load analytics data')
      } finally {
        this.loading = false
      }
    },
    showDetails(analytics) {
      this.selectedAnalytics = analytics
    },
    closeDetails() {
      this.selectedAnalytics = null
    },
    getFileName(path) {
      if (!path) return 'Unknown File'
      return path.split('/').pop() || path
    },
    formatDate(dateString) {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleString()
    }
  }
}
</script>

<style scoped>
.analytics-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.analytics-header {
  margin-bottom: 30px;
}

.analytics-header h2 {
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

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.analytics-section {
  margin-bottom: 40px;
}

.analytics-section h3 {
  margin-bottom: 20px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
  border-left: 4px solid #007bff;
}

.stat-card h4 {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
  font-weight: 600;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #007bff;
}

.best-worst-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.interview-card {
  padding: 20px;
  border-radius: 8px;
  color: white;
}

.interview-card.best {
  background: linear-gradient(135deg, #28a745, #20c997);
}

.interview-card.worst {
  background: linear-gradient(135deg, #dc3545, #fd7e14);
}

.interview-card h4 {
  margin: 0 0 10px 0;
}

.interview-card p {
  margin: 0 0 10px 0;
  opacity: 0.9;
  font-size: 14px;
}

.interview-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.interview-analytics-card {
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

.progress-section {
  margin-bottom: 15px;
}

.progress-item {
  margin-bottom: 10px;
}

.progress-item span {
  display: block;
  margin-bottom: 5px;
  font-size: 14px;
  color: #666;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background-color: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.progress-fill.answered {
  background-color: #28a745;
}

.progress-fill.unanswered {
  background-color: #dc3545;
}

.accuracy-section {
  margin-bottom: 15px;
}

.accuracy-item {
  margin-bottom: 10px;
  font-size: 14px;
  color: #666;
}

.confidence-breakdown {
  display: flex;
  gap: 15px;
  font-size: 12px;
}

.high-conf {
  color: #28a745;
}

.med-conf {
  color: #ffc107;
}

.low-conf {
  color: #dc3545;
}

.date-info {
  font-size: 12px;
  color: #999;
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

.detail-grid {
  display: grid;
  gap: 15px;
}

.detail-item {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: 10px;
  align-items: start;
}

.detail-item label {
  font-weight: 600;
  color: #666;
}

.confidence-detail {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .best-worst-section {
    grid-template-columns: 1fr;
  }
  
  .detail-item {
    grid-template-columns: 1fr;
    gap: 5px;
  }
}

.analytics-container,
.analytics-container * {
  color: black !important;
}


</style>
