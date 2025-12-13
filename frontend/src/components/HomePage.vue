<template>
  <div class="home-page">
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">🎯 Interview Parser & Analyzer</h1>
        <p class="hero-subtitle">
          Your intelligent companion for interview preparation, recording, and analysis. 
          Transform the way you handle interviews with AI-powered insights and real-time feedback.
        </p>
        <div class="hero-features">
          <div class="feature-badge">
            <span class="feature-icon">🤖</span>
            <span>AI-Powered Analysis</span>
          </div>
          <div class="feature-badge">
            <span class="feature-icon">🎙️</span>
            <span>Real-time Recording</span>
          </div>
          <div class="feature-badge">
            <span class="feature-icon">📊</span>
            <span>Detailed Analytics</span>
          </div>
        </div>
      </div>
    </div>

    <div class="pages-grid">
      <div 
        v-for="page in pages" 
        :key="page.id"
        class="page-card"
        @click="navigateToPage(page.view)"
      >
        <div class="card-header">
          <div class="page-icon">{{ page.icon }}</div>
          <h3 class="page-title">{{ page.title }}</h3>
        </div>
        
        <div class="card-content">
          <p class="page-description">{{ page.description }}</p>
          
          <div class="page-features">
            <div 
              v-for="feature in page.features" 
              :key="feature"
              class="feature-item"
            >
              <span class="feature-check">✓</span>
              <span>{{ feature }}</span>
            </div>
          </div>
        </div>
        
        <div class="card-footer">
          <button class="navigate-button">
            {{ page.actionText }}
            <span class="arrow">→</span>
          </button>
        </div>
      </div>
    </div>

    <div class="workflow-section">
      <h2 class="section-title">🔄 How It Works</h2>
      <div class="workflow-steps">
        <div class="step">
          <div class="step-number">1</div>
          <div class="step-content">
            <h4>Record or Upload</h4>
            <p>Capture audio directly or upload existing interview files in various formats</p>
          </div>
        </div>
        <div class="step-connector">→</div>
        <div class="step">
          <div class="step-number">2</div>
          <div class="step-content">
            <h4>AI Transcription</h4>
            <p>Advanced AI converts speech to text with high accuracy</p>
          </div>
        </div>
        <div class="step-connector">→</div>
        <div class="step">
          <div class="step-number">3</div>
          <div class="step-content">
            <h4>Smart Analysis</h4>
            <p>Extract questions, evaluate answers, and provide detailed feedback</p>
          </div>
        </div>
        <div class="step-connector">→</div>
        <div class="step">
          <div class="step-number">4</div>
          <div class="step-content">
            <h4>View Analytics</h4>
            <p>Track performance trends and improve your interview skills</p>
          </div>
        </div>
      </div>
    </div>

    <div class="tips-section">
      <h2 class="section-title">💡 Pro Tips</h2>
      <div class="tips-grid">
        <div class="tip-card">
          <div class="tip-icon">🎯</div>
          <h4>Mock Interviews</h4>
          <p>Practice with AI-generated questions tailored to your CV and target position</p>
        </div>
        <div class="tip-card">
          <div class="tip-icon">📈</div>
          <h4>Track Progress</h4>
          <p>Monitor your improvement over time with detailed analytics and performance metrics</p>
        </div>
        <div class="tip-card">
          <div class="tip-icon">🔧</div>
          <h4>Customize Settings</h4>
          <p>Configure API keys and recording preferences for optimal performance</p>
        </div>
        <div class="tip-card">
          <div class="tip-icon">📱</div>
          <h4>Multi-format Support</h4>
          <p>Works with MP3, WAV, M4A, MP4, MOV, AVI and more formats</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineEmits } from 'vue'

const emit = defineEmits(['navigate'])

const pages = [
  {
    id: 'upload',
    view: 'upload',
    icon: '📤',
    title: 'Upload & Transcribe',
    description: 'Transform your audio and video interviews into text with cutting-edge AI transcription. Perfect for processing recorded interviews from any source.',
    actionText: 'Upload Files',
    features: [
      'Drag & drop interface',
      'Multiple format support (MP3, WAV, M4A, MP4, MOV, AVI)',
      'Real-time processing progress',
      'Automatic question extraction',
      'Instant transcription and analysis'
    ]
  },
  {
    id: 'record',
    view: 'record',
    icon: '🎙️',
    title: 'Record Audio',
    description: 'Capture crystal-clear audio directly from your microphone. Choose between interview practice mode or meeting recording with intelligent analysis.',
    actionText: 'Start Recording',
    features: [
      'High-quality audio recording',
      'Device selection and management',
      'Real-time recording status',
      'Interview or meeting modes',
      'Automatic processing after recording'
    ]
  },
  {
    id: 'interview',
    view: 'interview',
    icon: '🎯',
    title: 'Mock Interview',
    description: 'Practice your interview skills with an AI interviewer that generates personalized questions based on your CV and target position.',
    actionText: 'Start Practice',
    features: [
      'Personalized question generation',
      'Real-time conversation with AI',
      'Instant feedback and scoring',
      'Text and audio responses',
      'Comprehensive performance analysis'
    ]
  },
  {
    id: 'interviews',
    view: 'interviews',
    icon: '📋',
    title: 'Interview Management',
    description: 'Organize, review, and manage all your interview data in one place. Edit questions, track accuracy, and monitor your progress over time.',
    actionText: 'Manage Interviews',
    features: [
      'View all interview history',
      'Edit questions and answers',
      'Track accuracy scores',
      'Filter by date range',
      'Export interview data'
    ]
  },
  {
    id: 'calls',
    view: 'calls',
    icon: '📞',
    title: 'Calls & Meetings',
    description: 'Transform your business meetings and calls into actionable insights. Extract tasks, decisions, and next steps automatically.',
    actionText: 'View Calls',
    features: [
      'Meeting transcription and analysis',
      'Action item extraction',
      'Decision tracking',
      'Participant insights',
      'Summary generation'
    ]
  },
  {
    id: 'analytics',
    view: 'analytics',
    icon: '📊',
    title: 'Analytics Dashboard',
    description: 'Gain deep insights into your interview performance with comprehensive analytics, trends, and detailed metrics to track your improvement.',
    actionText: 'View Analytics',
    features: [
      'Performance trends over time',
      'Accuracy score analysis',
      'Question category breakdown',
      'Best and worst interview highlights',
      'Custom date range filtering'
    ]
  },
  {
    id: 'explorer',
    view: 'explorer',
    icon: '📂',
    title: 'File Explorer',
    description: 'Browse and manage all your processed files, transcripts, and analysis results in an intuitive file management interface.',
    actionText: 'Browse Files',
    features: [
      'Organized file structure',
      'Quick file preview',
      'Easy navigation',
      'Search functionality',
      'File type filtering'
    ]
  },
  {
    id: 'apikey',
    view: 'apikey',
    icon: '🔑',
    title: 'API Key Management',
    description: 'Securely manage your API keys for transcription and analysis services. Keep your credentials safe and organized.',
    actionText: 'Manage Keys',
    features: [
      'Secure credential storage',
      'Multiple provider support',
      'Easy key rotation',
      'Connection testing',
      'Service status monitoring'
    ]
  }
]

const navigateToPage = (view) => {
  emit('navigate', view)
}
</script>

<style scoped>
.home-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
}

/* Hero Section */
.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  padding: 60px 40px;
  margin-bottom: 40px;
  text-align: center;
  color: white;
  box-shadow: 0 10px 40px rgba(102, 126, 234, 0.3);
}

.hero-title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 20px;
  background: linear-gradient(45deg, #ffffff, #f0f4ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.3rem;
  line-height: 1.6;
  margin-bottom: 40px;
  opacity: 0.9;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
}

.hero-features {
  display: flex;
  justify-content: center;
  gap: 30px;
  flex-wrap: wrap;
}

.feature-badge {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.2);
  padding: 10px 20px;
  border-radius: 25px;
  backdrop-filter: blur(10px);
  font-weight: 500;
}

.feature-icon {
  font-size: 1.2rem;
}

/* Pages Grid */
.pages-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 25px;
  margin-bottom: 60px;
}

.page-card {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid #e9ecef;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.page-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  border-color: #667eea;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.page-icon {
  font-size: 2.5rem;
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 15px;
  flex-shrink: 0;
}

.page-title {
  font-size: 1.4rem;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.card-content {
  flex: 1;
  margin-bottom: 25px;
}

.page-description {
  color: #666;
  line-height: 1.6;
  margin-bottom: 20px;
  font-size: 1rem;
}

.page-features {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.9rem;
  color: #555;
}

.feature-check {
  color: #28a745;
  font-weight: bold;
  flex-shrink: 0;
}

.card-footer {
  margin-top: auto;
}

.navigate-button {
  width: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 10px;
  padding: 12px 20px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.navigate-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.arrow {
  transition: transform 0.3s ease;
}

.page-card:hover .arrow {
  transform: translateX(5px);
}

/* Workflow Section */
.workflow-section {
  margin-bottom: 60px;
}

.section-title {
  text-align: center;
  font-size: 2.2rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 40px;
}

.workflow-steps {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 200px;
}

.step-number {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  font-weight: bold;
  margin-bottom: 15px;
}

.step-content h4 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.step-content p {
  color: #666;
  font-size: 0.9rem;
  line-height: 1.4;
}

.step-connector {
  font-size: 1.5rem;
  color: #667eea;
  font-weight: bold;
}

/* Tips Section */
.tips-section {
  margin-bottom: 40px;
}

.tips-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 25px;
}

.tip-card {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 25px;
  text-align: center;
  border: 1px solid #e9ecef;
  transition: all 0.3s ease;
}

.tip-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.tip-icon {
  font-size: 2rem;
  margin-bottom: 15px;
}

.tip-card h4 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
}

.tip-card p {
  color: #666;
  font-size: 0.9rem;
  line-height: 1.4;
}

/* Responsive Design */
@media (max-width: 768px) {
  .hero-title {
    font-size: 2rem;
  }
  
  .hero-subtitle {
    font-size: 1.1rem;
  }
  
  .hero-features {
    gap: 15px;
  }
  
  .pages-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .workflow-steps {
    flex-direction: column;
    gap: 30px;
  }
  
  .step-connector {
    transform: rotate(90deg);
  }
  
  .tips-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .hero-section {
    padding: 40px 20px;
  }
  
  .page-card {
    padding: 20px;
  }
  
  .tip-card {
    padding: 20px;
  }
}
</style>
