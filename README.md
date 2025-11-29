# Interview Parser

A powerful Go-based tool for transcribing and analyzing interview recordings. The application processes audio/video files to extract meaningful questions from interviewers and classifies them as answered or unanswered by candidates, providing detailed analytics about the interview.

## Features

- **Audio/Video Transcription**: Uses OpenAI's advanced models to transcribe interview recordings
- **Question Analysis**: Intelligently extracts and categorizes questions from interviewers
- **Answer Classification**: Determines which questions were adequately answered by candidates
- **Parallel Processing**: Efficiently processes large files using configurable parallel workers
- **Chunk-based Processing**: Splits large recordings into manageable chunks for optimal processing
- **Detailed Analytics**: Generates comprehensive reports with question-answer pairs and accuracy scores

## Requirements

- Go 1.23.0 or higher
- [FFmpeg](https://ffmpeg.org/) and [FFprobe](https://ffmpeg.org/) installed and available in PATH
- Valid OpenAI API key with access to the specified models
- SQLite3 (for API key storage)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/mrbelka12000/interview_parser.git
cd interview_parser
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o interview-parser ./cmd
```

## Usage

### Basic Usage

```bash
# Process an interview recording
./interview-parser --input interview.mp4

# Use custom output file
./interview-parser --input interview.mp4 --output my_analysis.md

# Specify chunk size for large files
./interview-parser --input interview.mp4 --chunk_seconds 120
```

### Advanced Usage

```bash
# Use different AI models
./interview-parser --input interview.mp4 --gpt_transcribe_model whisper-1 --gpt_classify_questions_model gpt-4

# Load previously created chunks instead of creating new ones
./interview-parser --load_chunks --chunks_dir ./output/chunks

# Use multiple parallel workers for faster processing
./interview-parser --input interview.mp4 --parallel_workers 8
```

### Environment Variables

You can set your OpenAI API key using environment variables:

```bash
export OPENAI_API_KEY=your_api_key_here
./interview-parser --input interview.mp4
```

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Path to input audio/video file to transcribe and analyze | Required (unless using `--load_chunks`) |
| `--output` | Path to output file for analysis results | `analytics.md` |
| `--chunk_seconds` | Duration in seconds for each audio chunk during processing | `100` |
| `--gpt_transcribe_model` | OpenAI model to use for audio transcription | `gpt-4o-transcribe` |
| `--gpt_classify_questions_model` | OpenAI model to use for question analysis | `o3` |
| `--load_chunks` | Load previously created chunks instead of creating new ones | `false` |
| `--chunks_dir` | Directory to store/load audio chunks | `~/.interview_parser/output/chunks` |
| `--parallel_workers` | Number of parallel workers for processing | Number of CPU cores |
| `--openai_api_key` | OpenAI API key | Can be set via `OPENAI_API_KEY` env var |

## Output

The tool generates two files:

1. **Transcript File**: Saved to `~/.interview_parser/transcript.txt` - Contains the full transcribed text
2. **Analysis File**: Saved to the specified output path (default: `analytics.md`) - Contains categorized questions and answers

### Analysis Format

The analysis file contains:

#### Answered Questions
```
### [Question text]
[Answer summary]
[Accuracy score]
```

#### Unanswered Questions
```
### [Question text]
[Candidate's response]
[Reason for being classified as unanswered]
```

## How It Works

1. **Audio Processing**: The tool splits the input audio/video file into manageable chunks
2. **Transcription**: Each chunk is transcribed using OpenAI's transcription models
3. **Text Formatting**: The transcribed text is formatted for better readability
4. **Question Analysis**: The system analyzes the transcript to extract interviewer questions
5. **Classification**: Questions are classified as answered or unanswered based on candidate responses
6. **Report Generation**: A detailed report is generated with accuracy scores and reasoning

## Question Classification Rules

The system follows specific rules to ensure accurate analysis:

### Questions That Are Ignored
- Sound checks ("Слышно?", "Нормально слышно?")
- Organizational questions ("Все понятно?", "Продолжаем?")
- Knowledge checks ("Ты это знал?")
- "Do you have questions?" type questions
- Rhetorical or emotional questions
- Clarifications without new meaning

### Answer Classification Criteria
A question is considered "answered" only if:
- The answer comes from the candidate
- The answer is relevant to the question
- The accuracy score is 0.7 or higher

### Accuracy Scoring
- **1.0**: Complete, accurate answer without errors
- **0.7-0.9**: Good but incomplete answer
- **0.7**: Minimum threshold for answered questions
- **0.3-0.7**: Partial answer with errors → unanswered
- **0.0-0.2**: "Don't know", avoidance, wrong answer → unanswered

## Configuration

The application stores configuration and data in `~/.interview_parser/`:
- `local.db`: SQLite database for API key storage
- `output/chunks/`: Directory for audio chunks
- `transcript.txt`: Latest transcript file

## Troubleshooting

### Common Issues

1. **FFmpeg not found**: Ensure FFmpeg and FFprobe are installed and in your PATH
2. **API Key Issues**: Verify your OpenAI API key has access to the specified models
3. **Memory Issues**: Reduce `chunk_seconds` or `parallel_workers` for large files
4. **Model Access**: Check that your OpenAI project has access to the selected models

### Error Messages

- `Your project does not have access to model`: Add the required model to your OpenAI project
- `No OPENAI API key found`: Set the API key via flag or environment variable
- `config is nil`: Check your configuration parameters

## Development

### Project Structure

```
interview_parser/
├── cmd/                 # Main application entry point
├── client/              # OpenAI client integration
├── config/              # Configuration management
├── converter.go         # Audio processing utilities
├── db.go               # Database operations
├── formatter.go        # Text formatting
├── saver.go            # File output operations
└── go.mod              # Go module definition
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/mrbelka12000/interview_parser.git
cd interview_parser

# Install dependencies
go mod tidy

# Build for your current platform
go build -o interview-parser ./cmd

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o interview-parser-linux ./cmd
GOOS=windows GOARCH=amd64 go build -o interview-parser.exe ./cmd
GOOS=darwin GOARCH=amd64 go build -o interview-parser-mac ./cmd
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Open an issue on GitHub
- Check the troubleshooting section above
- Run `./interview-parser --help` for detailed usage information