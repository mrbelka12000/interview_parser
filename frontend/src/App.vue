<script setup>
import FileExplorer from './components/FileExplorer.vue'
import HelloWorld from './components/HelloWorld.vue'
import { ref } from 'vue'

const currentView = ref('explorer')

const switchView = (view) => {
  currentView.value = view
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
          @click="switchView('greet')" 
          :class="{ active: currentView === 'greet' }"
          class="nav-tab"
        >
          👋 Greet
        </button>
      </nav>
    </header>
    
    <main class="app-main">
      <FileExplorer v-if="currentView === 'explorer'" />
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
