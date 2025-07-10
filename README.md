# clify - Command-line Assistant

![Demo](demo.svg)

A CLI tool that tells you how to run command-line programs on Linux, macOS, and
Windows command-line using human language.

## Features

- **Caching**: Stores answers locally to avoid repeated API calls for the same questions
- **Autocomplete**: Suggests previously asked questions as you type
- **Cross-platform**: Works on Linux, macOS, and Windows. Automatically detects your OS and provides relevant command-line solutions.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install aktagon/tap/clify
```

### From Source

```bash
go build -o clify
```

### Using make

```bash
make build
make install
```

## Usage

### Interactive Mode

```bash
clify
```

This starts an interactive session where you can ask questions and get
autocomplete suggestions based on your history.

### Direct Query Mode

```bash
clify "create a git branch from a tag"
```

This starts an interactive session with the query "create a git branch from a tag".

**Example output:**

```
Create a new Git branch based on a specific tag, allowing you to start development from that tagged point in history

🟡 git checkout -b new-branch-name tag-name
    Create and switch to a new branch based on the specified tag

🟢 git branch new-branch-name tag-name
    Create a new branch from a tag without switching to it

🟢 git tag --list
    List all available tags to see which one you want to branch from
```

## Configuration

Run the setup command to configure your API key:

```bash
clify setup
```

This will prompt you to enter your Anthropic API key and create a configuration file.

The app stores configuration in `~/.clify/config.yaml`:

```yaml
api_key: "your-api-key-here"
cache_file: "~/.clify/cache.json"
model: "claude-3-sonnet-20240229"
```

### Environment Variables

- `ANTHROPIC_API_KEY` - Optional. Overrides the API key from config file if set

## Cache

The tool automatically caches responses in `~/.clify/cache.json` to avoid
repeated API calls for the same questions.

## Dependencies

- Go 1.21+
- github.com/aktagon/llmkit
- github.com/charmbracelet/bubbletea

## Development

To build from source:

```bash
git clone https://github.com/aktagon/clify
cd clify
go build -o clify
```

To run tests:

```bash
go test ./...
```
