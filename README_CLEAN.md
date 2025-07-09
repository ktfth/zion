# 🚀 Zion CLI - Refactored

A clean, elegant command-line tool for AI-powered project scaffolding.

## ✨ Features

- **AI-Powered Generation**: Generate complete project structures using AI
- **Multiple Providers**: Support for Gemini and OpenAI
- **Clean Architecture**: Well-structured, maintainable codebase
- **Comprehensive Testing**: Full test coverage for reliability
- **Retry Logic**: Robust error handling and retry mechanisms
- **Simple Interface**: Clean, intuitive command-line interface

## 🏗️ Architecture

```
internal/
├── core/           # Core business logic
│   ├── interfaces.go
│   ├── project.go
│   └── generator.go
├── providers/      # AI provider implementations
│   ├── factory.go
│   ├── gemini.go
│   └── openai.go
cmd/               # CLI commands
├── root.go
├── scaffold.go
└── provider.go
```

## 📦 Installation

```bash
# Clone the repository
git clone https://github.com/ktfth/zion.git
cd zion

# Build the application
go build -o zion .

# Or install directly
go install .
```

## 🚀 Usage

### Basic Usage

```bash
# Generate a Go project
zion scaffold -l go -n my-api -d "REST API with authentication"

# Generate a Python project
zion scaffold -l python -n ml-project -d "Machine learning project"

# Generate a TypeScript project
zion scaffold -l typescript -n web-app -d "Modern web application"
```

### Provider Management

```bash
# List available providers
zion provider list

# Test a provider
zion provider test gemini
zion provider test openai
```

### Configuration

Set your API keys as environment variables:

```bash
# For Gemini
export GEMINI_API_KEY="your-gemini-api-key"

# For OpenAI
export OPENAI_API_KEY="your-openai-api-key"
```

### Advanced Usage

```bash
# Use specific provider
zion scaffold -l go -n my-project -d "Description" -p gemini -k "your-key"

# Use specific model
zion scaffold -l python -n project -d "Description" -p openai -m "gpt-4"

# Set retry attempts
zion scaffold -l typescript -n app -d "Description" -r 5
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 🛠️ Development

### Project Structure

- `internal/core/`: Core business logic and interfaces
- `internal/providers/`: AI provider implementations
- `cmd/`: CLI command definitions
- `main.go`: Application entry point

### Adding New Providers

1. Implement the `AIProvider` interface in `internal/providers/`
2. Add the provider to the factory in `internal/providers/factory.go`
3. Update the provider command in `cmd/provider.go`

### Code Standards

- Follow Go best practices
- Write comprehensive tests
- Use meaningful variable and function names
- Add documentation for public APIs
- Keep functions small and focused

## 📊 Differences from Original

### Removed
- Complex layered generation system
- Plugin system (can be re-added if needed)
- Evaluation system (can be re-added if needed)
- Contextual mode with llms.txt
- Multiple provider variants (OpenRouter, Claude)
- Excessive configuration options

### Improved
- Clean, modular architecture
- Comprehensive test coverage
- Simplified command interface
- Better error handling
- Cleaner code structure
- Focused functionality

### Maintained
- AI-powered project generation
- Multiple AI provider support
- Retry logic
- Command-line interface
- Core scaffolding functionality

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📝 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- Original Zion CLI for inspiration
- Cobra CLI framework
- AI providers (Google Gemini, OpenAI)
