# Clify

[![Build Status](https://github.com/aktagon/clify/workflows/CI/badge.svg)](https://github.com/aktagon/clify/actions)
[![Release](https://img.shields.io/github/v/release/aktagon/clify)](https://github.com/aktagon/clify/releases)

Natural language to shell commands. Powered by Claude.

![Demo](demo.svg)

## Install

```bash
brew tap aktagon/clify https://github.com/aktagon/clify
brew install clify
clify setup   # paste your Anthropic API key
```

From source:

```bash
git clone https://github.com/aktagon/clify && cd clify && go build -o clify
```

## Use

```bash
clify "find large files in current directory"
```

```
🟡 find . -type f -size +100M -exec ls -lh {} \;
    Find files larger than 100MB and show their sizes

🟢 du -ah . | sort -hr | head -20
    Show top 20 largest files and directories by size

🟢 ls -laSh | head -20
    List files sorted by size (largest first)
```

Interactive mode (autocomplete, history):

```bash
clify
```

## Configure

`~/.clify/config.yaml`:

```yaml
api_key: "sk-ant-..."
cache_file: "~/.clify/cache.json"
model: "<claude-model-id>"
```

`ANTHROPIC_API_KEY` overrides the config file.

## Behavior

- Caches responses locally. No duplicate API calls.
- Detects Linux, macOS, or Windows and adapts commands.
- Returns ranked alternatives, not a single guess.

## Build

Requires Go 1.21+.

```bash
go build -o clify   # or: make build
go test ./...
```

---

Built by [Aktagon](https://aktagon.com). Applied AI for regulated markets.
