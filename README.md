# Interview Parser Desktop App

A powerful desktop application built with Wails (Go backend + Vue.js frontend) for transcribing and analyzing interview recordings. The application processes audio/video files to extract meaningful questions from interviewers and classifies them as answered or unanswered by candidates, providing detailed analytics about the interview.

## Overview

Interview Parser Desktop App provides a user-friendly graphical interface for processing interview recordings. It combines the power of Go's backend processing with Vue.js's modern frontend to create a seamless experience for interview analysis.

## Features

### Core Functionality
- **Audio/Video Transcription**: Uses OpenAI's advanced models to transcribe interview recordings
- **Question Analysis**: Intelligently extracts and categorizes questions from interviewers
- **Answer Classification**: Determines which questions were adequately answered by candidates
- **Parallel Processing**: Efficiently processes large files using configurable parallel workers
- **Chunk-based Processing**: Splits large recordings into manageable chunks for optimal processing
- **Detailed Analytics**: Generates comprehensive reports with question-answer pairs and accuracy scores

### Desktop Application Features
- **Intuitive GUI**: Modern Vue.js interface with tabbed navigation
- **File Explorer**: Built-in file browser for easy file management
- **Progress Tracking**: Real-time progress updates during processing
- **API Key Management**: Secure storage and management of OpenAI API keys
- **File Upload**: Drag-and-drop or file picker for media files
- **Content Viewer**: Built-in viewer for transcripts and analysis results

## Architecture

### Backend (Go)
- **Wails Framework**: Desktop application framework
- **OpenAI Integration**: Direct API integration for transcription and analysis
- **SQLite Database**: Local storage for API keys and configuration
- **FFmpeg Integration**: Audio/video processing and chunking
- **Parallel Processing**: Concurrent chunk processing for performance

### Frontend (Vue.js)
- **Vue 3**: Modern reactive framework
- **Component-based Architecture**: Modular UI components
- **Responsive Design**: Adapts to different screen sizes
- **Real-time Updates**: Event-driven communication with backend

## Requirements

### System Requirements
- **Operating System**: Windows, macOS, or Linux
- **Go**: 1.25 or higher
- **Node.js**: For frontend development
- **FFmpeg**: Must be installed and available in PATH
- **OpenAI API Key**: Valid key with access to transcription and analysis models

### Development Requirements
- **Wails v2**: Desktop application framework
- **Vue.js 3**: Frontend framework
- **Vite**: Build tool for frontend

## Installation

### Prerequisites

1. **Install Go** (1.25+):
   ```bash
   # macOS
   brew install go
   
   # Ubuntu/Debian
   sudo apt-get install golang-go
   
   # Windows
   # Download from https://golang.org/dl/
   ```

2. **Install Node.js** (for development):
   ```bash
   # macOS
   brew install node
   
   # Ubuntu/Debian
   sudo apt-get install nodejs npm
   
   # Windows
   # Download from https://nodejs.org/
   ```

3. **Install Wails**:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **Install FFmpeg**:
   ```bash
   # macOS
   brew install ffmpeg
   
   # Ubuntu/Debian
   sudo apt-get install ffmpeg
   
   # Windows
   # Download from https://ffmpeg.org/download.html
   ```

### Building from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/mrbelka12000/interview_parser.git
   cd interview_parser
   ```

2. **Install dependencies**:
   ```bash
   # Backend dependencies
   go mod download
   
   # Frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Build the application**:
   ```bash
   # Development build
   wails dev
   
   # Production build
   wails build
   ```

### Running the Application

1. **Development mode** (with hot reload):
   ```bash
   wails dev
   ```

2. **Production mode**:
   ```bash
    wails build -platform darwin/arm64

    APP="build/bin/interview_parser.app"
    cp /usr/local/bin/ffprobe "$APP/Contents/MacOS/ffprobe"
    cp /usr/local/bin/ffmpeg  "$APP/Contents/MacOS/ffmpeg"
   ```

## Usage

### First Time Setup

1. **Launch the application** after installation
2. **Configure API Key**:
   - Navigate to the "🔑 API Key" tab
   - Enter your OpenAI API key (should start with "sk-")
   - Click "Save API Key"
3. **Verify Setup**: The application will validate the API key automatically

### Processing Interview Recordings

1. **Upload File**:
   - Navigate to the "🎤 Upload & Transcribe" tab
   - Click "Pick File" or drag and drop a media file
   - Supported formats: MP3, WAV, M4A, MP4, MOV, AVI

2. **Configure Processing**:
   - Choose chunk size (default: 100 seconds)
   - Select number of parallel workers (default: CPU cores)
   - Choose whether to load existing chunks

3. **Start Processing**:
   - Click "Process File"
   - Monitor progress in real-time
   - Wait for completion

4. **View Results**:
   - Navigate to "📂 File Explorer" tab
   - Find generated files (transcript and analysis)
   - Click on files to view content

### File Management

1. **File Explorer**:
   - Browse your working directory
   - View file details (size, type, etc.)
   - Open files for content viewing

2. **Content Viewer**:
   - View transcript files (.txt)
   - View analysis files (.md)
   - Navigate back using breadcrumb trail

## Configuration

### Application Settings

The application stores configuration in:
- **Database**: `~/.interview_parser/local.db`
- **Output Directory**: `~/.interview_parser/`
- **Chunks Directory**: `~/.interview_parser/output/chunks/`

### Processing Parameters

- **Chunk Duration**: 100 seconds (configurable)
- **Parallel Workers**: Number of CPU cores (configurable)
- **Transcription Model**: `gpt-4o-transcribe` (configurable)
- **Analysis Model**: `o3` (configurable)
- **Language**: Russian (ru) or English (en)

## Output Files

### Transcript File
- **Location**: `~/.interview_parser/[filename]_transcript_[timestamp].txt`
- **Format**: Plain text transcribed content
- **Encoding**: UTF-8

### Analysis File
- **Location**: `~/.interview_parser/[filename]_analysis_[timestamp].md`
- **Format**: Markdown with structured analysis
- **Sections**:
  - Answered Questions with accuracy scores
  - Unanswered Questions with reasoning
  - Summary statistics

## Analysis Details

### Question Classification Rules

The system follows specific rules to ensure accurate analysis and supports multiple languages:

#### Language Support
- **Russian (ru)**: Default language, uses Russian prompts and analysis
- **English (en)**: Uses English prompts and analysis for English-language interviews

#### Questions That Are Ignored

**Russian interviews:**
- Sound checks ("Слышно?", "Нормально слышно?")
- Organizational questions ("Все понятно?", "Продолжаем?")
- Knowledge checks ("Ты это знал?")
- "Do you have questions?" type questions
- Rhetorical or emotional questions
- Clarifications without new meaning

**English interviews:**
- Sound checks ("Do you hear me?", "Is the sound okay?")
- Organizational questions ("Is everything clear?", "Shall we continue?")
- Knowledge-check questions ("Did you know this?")
- "Do you have any questions?" type questions
- Rhetorical or emotional questions
- Clarifications without new meaning

#### Answer Classification Criteria
A question is considered "answered" only if:
- The answer comes from the candidate
- The answer is relevant to the question
- The accuracy score is 0.7 or higher

#### Accuracy Scoring
- **1.0**: Complete, accurate answer without errors
- **0.7-0.9**: Good but incomplete answer
- **0.7**: Minimum threshold for answered questions
- **0.3-0.7**: Partial answer with errors → unanswered
- **0.0-0.2**: "Don't know", avoidance, wrong answer → unanswered

## Troubleshooting

### Common Issues

1. **FFmpeg not found**:
   - Install FFmpeg and ensure it's in your PATH
   - Verify installation with `ffmpeg -version`

2. **API Key Issues**:
   - Check that your API key starts with "sk-"
   - Verify your OpenAI project has access to the required models
   - Check your API key usage and limits

3. **Processing Errors**:
   - Reduce chunk size for large files
   - Decrease parallel workers if experiencing memory issues
   - Check file format compatibility

4. **Application Won't Start**:
   - Verify all dependencies are installed
   - Check system compatibility
   - Review build logs for errors

### Error Messages

- **"Your project does not have access to model"**:
  - Add the required model to your OpenAI project
  - Check your OpenAI project settings

- **"No API Key found"**:
  - Navigate to API Key tab and enter your key
  - Ensure the key format is correct (starts with "sk-")

- **"Failed to create chunk"**:
  - Check FFmpeg installation
  - Verify input file format
  - Check disk space

## Development

### Project Structure

```
interview_parser/
├── main.go                    # Application entry point
├── wails.json                 # Wails configuration
├── go.mod                     # Go module definition
├── internal/                  # Backend packages
│   ├── app/                   # Core application logic
│   │   └── app.go            # Main application struct
│   ├── client/                # OpenAI client integration
│   │   ├── client.go         # OpenAI client wrapper
│   │   ├── analyze.go        # Transcript analysis
│   │   ├── transcribe.go     # Audio transcription
│   │   ├── prompts.go        # AI prompts
│   │   └── validate.go       # API key validation
│   ├── config/                # Configuration management
│   │   └── config.go         # Configuration parsing
│   ├── parser/                # Media processing
│   │   ├── parser.go         # Parser interface
│   │   ├── converter.go      # Audio/video conversion
│   │   ├── formatter.go      # Text formatting
│   │   └── saver.go          # File output operations
│   └── db.go                 # Database operations
├── frontend/                  # Vue.js frontend
│   ├── src/
│   │   ├── App.vue           # Main application component
│   │   ├── main.js           # Frontend entry point
│   │   ├── components/       # Vue components
│   │   │   ├── ApiKeyManager.vue    # API key management
│   │   │   ├── FileContent.vue      # File content viewer
│   │   │   ├── FileExplorer.vue     # File browser
│   │   │   ├── FileUpload.vue       # File upload interface
│   │   │   └── HelloWorld.vue       # Example component
│   │   ├── assets/           # Static assets
│   │   │   ├── images/       # Image files
│   │   │   └── fonts/        # Font files
│   │   └── style.css         # Global styles
│   ├── package.json          # Frontend dependencies
│   └── vite.config.js        # Vite configuration
└── build/                     # Build output
    ├── appicon.png            # Application icon
    └── platform-specific/     # Platform-specific files
```

### Development Workflow

1. **Frontend Development**:
   ```bash
   cd frontend
   npm run dev    # Start development server
   ```

2. **Backend Development**:
   ```bash
   wails dev      # Start full application with hot reload
   ```

3. **Building for Production**:
   ```bash
   wails build    # Build for current platform
   wails build -platform darwin/amd64 -webview2 embed  # Cross-platform
   ```

### Adding New Features

1. **Backend Changes**:
   - Modify structs in `internal/app/app.go`
   - Update Wails bindings if adding new methods
   - Add database schema changes if needed

2. **Frontend Changes**:
   - Create new components in `frontend/src/components/`
   - Update routing in `frontend/src/App.vue`
   - Add new API calls in component scripts

3. **API Integration**:
   - Update client methods in `internal/client/`
   - Add validation logic as needed
   - Update TypeScript bindings

## Contributing

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**:
   ```bash
   git commit -m 'Add some amazing feature'
   ```
4. **Push to the branch**:
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Development Guidelines

- Follow Go conventions for backend code
- Use Vue 3 Composition API for frontend components
- Add TypeScript definitions for new API methods
- Update documentation for new features
- Test changes across different platforms

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:

- **GitHub Issues**: Open an issue on the repository
- **Documentation**: Check this README and inline comments
- **Community**: Join discussions in GitHub Discussions

### Contact

- **Email**: beka.teka11@gmail.com
- **GitHub**: https://github.com/mrbelka12000/interview_parser

## Acknowledgments

- **Wails Team**: For the excellent desktop application framework
- **OpenAI**: For providing powerful AI models
- **Vue.js Team**: For the reactive frontend framework
- **FFmpeg Team**: For robust media processing tools

---

**Interview Parser Desktop App** - Transform your interview recordings into actionable insights with the power of AI.
