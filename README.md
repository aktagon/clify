# Clify: Transform Natural Language into Terminal Commands

[![Build Status](https://github.com/aktagon/clify/workflows/CI/badge.svg)](https://github.com/aktagon/clify/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/aktagon/clify)](https://github.com/aktagon/clify)
[![Release](https://img.shields.io/github/v/release/aktagon/clify)](https://github.com/aktagon/clify/releases)

![Demo](demo.svg)

Ask "create a git branch from a tag" and get the exact commands you need. Clify converts natural language into executable terminal commands using AI.

## Quick Start

```bash
# Install
brew tap aktagon/clify https://github.com/aktagon/clify
brew install clify

# Configure with your Anthropic API key
clify setup

# Ask questions
clify "find large files in current directory"
```

**Output:**

```
🟡 find . -type f -size +100M -exec ls -lh {} \;
    Find files larger than 100MB and show their sizes

🟢 du -ah . | sort -hr | head -20
    Show top 20 largest files and directories by size

🟢 ls -laSh | head -20
    List files sorted by size (largest first)
```

## Why Choose Clify?

| Tool      | Strengths                              | Limitations           |
| --------- | -------------------------------------- | --------------------- |
| **Clify** | AI-powered, contextual, conversational | Requires API key      |
| `tldr`    | Fast, offline, community examples      | Limited coverage      |
| `man`     | Comprehensive documentation            | Complex, overwhelming |

Clify understands context, suggests alternatives, and adapts to your operating system automatically.

## Installation

### macOS and Linux

```bash
brew tap aktagon/clify https://github.com/aktagon/clify
brew install clify
```

### From Source

```bash
git clone https://github.com/aktagon/clify
cd clify
go build -o clify
```

### Build with Make

```bash
make build
make install
```

## Usage

### Interactive Mode

```bash
clify
```

Start a session with autocomplete and command history.

### Direct Queries

```bash
clify "compress multiple files into tar.gz"
clify "monitor network traffic on port 80"
clify "batch rename files with regex"
```

## Configuration

### Initial Setup

```bash
clify setup
```

Enter your [Anthropic API key](https://console.anthropic.com/) when prompted.

### Configuration File

Clify stores settings in `~/.clify/config.yaml`:

```yaml
api_key: "your-api-key-here"
cache_file: "~/.clify/cache.json"
model: "claude-3-sonnet-20240229"
```

### Environment Variables

Set `ANTHROPIC_API_KEY` to override the config file API key.

## Features

- **Smart Caching**: Stores responses locally to eliminate duplicate API calls
- **Autocomplete**: Suggests previous questions as you type
- **OS Detection**: Provides commands specific to Linux, macOS, or Windows
- **Multiple Solutions**: Shows alternative approaches for complex tasks

## Development

### Requirements

- Go 1.21+
- [Anthropic API key](https://console.anthropic.com/)

### Build

```bash
git clone https://github.com/aktagon/clify
cd clify
go build -o clify
# or make build
```

### Test

```bash
go test ./...
```

## Dependencies

- [github.com/aktagon/llmkit](https://github.com/aktagon/llmkit) - LLM integration
- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI
