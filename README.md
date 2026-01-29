<div align="center">

```
'##::: ##::'######:::::::::::'########:'##::::'##:'####:
 ###:: ##:'##... ##::::::::::... ##..:: ##:::: ##:. ##::
 ####: ##: ##:::..:::::::::::::: ##:::: ##:::: ##:: ##::
 ## ## ##:. ######::'#######:::: ##:::: ##:::: ##:: ##::
 ##. ####::..... ##:........:::: ##:::: ##:::: ##:: ##::
 ##:. ###:'##::: ##::::::::::::: ##:::: ##:::: ##:: ##::
 ##::. ##:. ######:::::::::::::: ##::::. #######::'####:
..::::..:::......:::::::::::::::..::::::.......:::....::
```

A beautiful terminal interface for searching NixOS packages in real-time.

**[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Development](#development)**

[![Built with Bubbletea](https://img.shields.io/badge/Built%20with-Bubbletea-5B8C5A?style=flat-square)](https://github.com/charmbracelet/bubbletea)
[![Go 1.25+](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)

</div>

---

## Features

- 🔍 **Real-time search** with instant results as you type
- ⚡ **Powered by official NixOS backend** - same data as search.nixos.org
- ⌨️ **Vim-style keybindings** - modal interface with Insert/Normal modes
- 🎨 **Beautiful Catppuccin theme** - easy on the eyes
- 📦 **Rich package details** - version, description, programs, platform support
- 📋 **One-click install commands** - copy any of 4 installation methods
- 🚀 **Fast & responsive** - debounced search, intelligent scrolling

## Installation

### From Source

**Prerequisites:** Go 1.25+

```bash
# Clone the repository
git clone https://github.com/yourusername/ns-tui.git
cd ns-tui

# Build with Go
go build -o bin/ns-tui ./cmd/ns-tui

# Or use Task (recommended)
task build
```

### With Nix

```bash
nix build
# or run directly without installing
nix run
```

## Usage

```bash
./bin/ns-tui
```

### Quick Start Guide

1. **Type to search** - Start in Insert mode, results appear as you type
2. **Navigate** - `Esc` → Normal mode, then `j`/`k` to move through results
3. **View package** - Press `Enter` to see full details
4. **Copy install command** - Use `Tab` to select method, `Enter` to copy
5. **Exit** - Press `q` to quit anytime

### Keybindings

<table>
<tr><th>Insert Mode</th><th>Normal Mode</th><th>Detail View</th></tr>
<tr><td>

| Key | Action |
|-----|--------|
| Type | Search packages |
| `↑` `↓` | Navigate results |
| `Esc` | → Normal mode |
| `q` | Quit |

</td><td>

| Key | Action |
|-----|--------|
| `i` | → Insert mode |
| `j` `k` | Move down/up |
| `g` `G` | Top/Bottom |
| `Enter` | View details |
| `q` | Quit |

</td><td>

| Key | Action |
|-----|--------|
| `Tab` | Cycle install methods |
| `Enter` | Copy command |
| `Esc` `b` | ← Back |
| `q` | Quit |

</td></tr>
</table>

## How It Works

ns-tui queries the official NixOS Elasticsearch backend (same as [search.nixos.org](https://search.nixos.org)) to provide real-time package search with relevance-based sorting.

**Built with:**
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework using The Elm Architecture
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

**Architecture:**
- Modal interface: Insert → Normal → Detail modes (vim-inspired)
- 300ms debounced search for optimal performance
- Responsive layout adapts to terminal size
- Catppuccin color scheme for comfort

## Development

### Project Structure

```
cmd/ns-tui/    # Application entry point
internal/
  ├── api/     # NixOS search backend client
  ├── models/  # Data structures (Package, Mode)
  ├── styles/  # Lipgloss theme & styles
  └── ui/      # Bubbletea components (Model, Update, View)
```

### Tasks

```bash
task build    # Build binary
task run      # Build and run
task fmt      # Format code
task lint     # Run linter
task clean    # Clean artifacts
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Code quality:** Run `task fmt` and `task lint` before committing.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with ❤️ using:
- [Charmbracelet](https://github.com/charmbracelet) TUI libraries
- [NixOS](https://nixos.org) package search API
- [Catppuccin](https://github.com/catppuccin/catppuccin) color palette

---

<div align="center">

**Star ⭐ this repo if you find it useful!**

</div>
