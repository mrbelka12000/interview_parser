<script setup>
import FileExplorer from './components/FileExplorer.vue'
import HelloWorld from './components/HelloWorld.vue'
import FileContent from './components/FileContent.vue'
import FileUpload from './components/FileUpload.vue'
import ApiKeyManager from './components/ApiKeyManager.vue'
import { ref } from 'vue'

const currentView = ref('explorer')
const selectedFilePath = ref('')

const switchView = (view) => {
  currentView.value = view
}

const showFileContent = (filePath) => {
  selectedFilePath.value = filePath
  currentView.value = 'fileContent'
}

const backToFileExplorer = () => {
  currentView.value = 'explorer'
  selectedFilePath.value = ''
}
</script>

<template>
  <div id="app">
    <header class="app-header">
      <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png"/>
      <nav class="nav-tabs">
        <button 
          @click="switchView('explorer')" 
          :class="{ active: currentView === 'explorer' }"
          class="nav-tab"
        >
          📂 File Explorer
        </button>
        <button 
          @click="switchView('upload')" 
          :class="{ active: currentView === 'upload' }"
          class="nav-tab"
        >
          🎤 Upload & Transcribe
        </button>
        <button 
          @click="switchView('apikey')" 
          :class="{ active: currentView === 'apikey' }"
          class="nav-tab"
        >
          🔑 API Key
        </button>
        <button 
          @click="switchView('greet')" 
          :class="{ active: currentView === 'greet' }"
          class="nav-tab"
        >
          👋 Greet
        </button>
      </nav>
    </header>
    
    <main class="app-main">
      <!-- Breadcrumb navigation -->
      <div v-if="currentView === 'fileContent'" class="breadcrumb">
        <button @click="backToFileExplorer" class="breadcrumb-link">
          ← Back to File Explorer
        </button>
        <span class="breadcrumb-separator">/</span>
        <span class="breadcrumb-current">File Content</span>
      </div>
      
      <FileExplorer 
        v-if="currentView === 'explorer'" 
        @file-selected="showFileContent" 
      />
      <FileUpload 
        v-else-if="currentView === 'upload'" 
      />
      <ApiKeyManager 
        v-else-if="currentView === 'apikey'" 
      />
      <FileContent 
        v-else-if="currentView === 'fileContent'" 
        :file-path="selectedFilePath" 
        @back="backToFileExplorer" 
      />
      <HelloWorld v-else-if="currentView === 'greet'" />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background-color: #f5f5f5;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  text-align: center;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

#logo {
  display: block;
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background-position: center;
  background-repeat: no-repeat;
  background-size: contain;
}

.nav-tabs {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 10px;
}

.nav-tab {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
  padding: 10px 20px;
  border-radius: 25px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.nav-tab:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
  transform: translateY(-2px);
}

.nav-tab.active {
  background: white;
  color: #667eea;
  border-color: white;
  box-shadow: 0 4px 15px rgba(0,0,0,0.2);
}

.app-main {
  flex: 1;
  padding: 0;
  background: white;
}

.breadcrumb {
  background: #f8f9fa;
  padding: 15px 20px;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  align-items: center;
  font-size: 14px;
}

.breadcrumb-link {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
  padding: 0;
  transition: color 0.3s;
}

.breadcrumb-link:hover {
  color: #0056b3;
  text-decoration: underline;
}

.breadcrumb-separator {
  margin: 0 10px;
  color: #6c757d;
}

.breadcrumb-current {
  color: #495057;
  font-weight: 500;
}

@media (max-width: 768px) {
  .app-header {
    padding: 15px;
  }
  
  #logo {
    width: 60px;
    height: 60px;
  }
  
  .nav-tabs {
    flex-direction: column;
    align-items: center;
    gap: 8px;
  }
  
  .nav-tab {
    width: 200px;
    text-align: center;
  }
}
</style>
