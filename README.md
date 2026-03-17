# devlog

A lightweight developer productivity CLI that automatically generates standup
reports from your Git commits and personal notes — all stored locally, never
sent anywhere.

---

## Features

- `devlog today` — see what you committed today
- `devlog note` — save a quick work note
- `devlog blocker` — log a blocker for your standup
- `devlog standup` — generate a full standup report
- `devlog clear` — reset your notes and/or blockers

---

## Requirements

- [Go 1.22+](https://go.dev/dl/)
- [Git](https://git-scm.com/downloads) installed and available in your PATH
- macOS, Linux, or Windows

---

## Installation

### 1. Clone the repository
```bash
git clone https://github.com/khurshiduktamov/devlog.git
cd devlog
```

### 2. Add Go binaries to your PATH

Skip this step if you have already done it before.

**macOS (zsh):**
```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc && source ~/.zshrc
```

**macOS/Linux (bash):**
```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc && source ~/.bashrc
```

**Windows (PowerShell — run once):**
```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\go\bin", [EnvironmentVariableTarget]::User)
```
Then restart your terminal for the change to take effect.

### 3. Install the binary

**macOS / Linux:**
```bash
go install .
```

**Windows (PowerShell):**
```powershell
go install .
```

### 4. Verify the installation

**macOS / Linux:**
```bash
devlog --help
```

**Windows (PowerShell):**
```powershell
devlog --help
```

You should see:
```
devlog helps you track your work by collecting git commits and notes,
then generating a clean standup report.

Usage:
  devlog [command]

Available Commands:
  blocker     Save a blocker that will appear in your standup report
  clear       Clear stored notes and/or blockers
  help        Help about any command
  note        Save a manual work note
  standup     Generate a daily standup report from commits and notes
  today       Show today's git commits
```

---

## Usage

### Show today's commits

Run this inside any Git repository:
```bash
devlog today
```

Output:
```
Today:
  • Fix authentication bug
  • Add Redis cache layer
  • Refactor payment service
```

---

### Save a note
```bash
devlog note "Investigated Redis timeout issue"
```

Output:
```
Note saved.
```

---

### Log a blocker
```bash
devlog blocker "Waiting for AWS credentials from DevOps"
```

Output:
```
Blocker saved.
```

---

### Generate a standup report
```bash
devlog standup
```

Output:
```
Yesterday:
  • Implemented Stripe webhook retry
  • Fixed order validation bug

Today:
  • Investigated Redis timeout issue

Blockers:
  • Waiting for AWS credentials from DevOps
```

---

### Clear stored data

Clear only notes:
```bash
devlog clear --notes
```

Clear only blockers:
```bash
devlog clear --blockers
```

Clear everything:
```bash
devlog clear --all
```

---

## Local Storage

All data is stored **locally on your machine**. Nothing is ever sent online.

| Platform | Path |
|---|---|
| macOS / Linux | `~/.devlog/` |
| Windows | `C:\Users\yourname\.devlog\` |

| File | Contents |
|---|---|
| `notes.json` | Your saved notes |
| `blockers.json` | Your saved blockers |

To inspect your stored data at any time:

**macOS / Linux:**
```bash
cat ~/.devlog/notes.json
cat ~/.devlog/blockers.json
```

**Windows (PowerShell):**
```powershell
type $env:USERPROFILE\.devlog\notes.json
type $env:USERPROFILE\.devlog\blockers.json
```

---

## Project Structure
```
devlog/
├── cmd/
│   ├── root.go        # Root command
│   ├── today.go       # devlog today
│   ├── note.go        # devlog note
│   ├── blocker.go     # devlog blocker
│   ├── standup.go     # devlog standup
│   └── clear.go       # devlog clear
│
├── internal/
│   ├── activity/
│   │   ├── model.go       # Unified Activity type
│   │   └── from_git.go    # Converts git commits to Activity
│   ├── git/
│   │   └── collector.go   # Runs git log and parses commits
│   ├── notes/
│   │   └── storage.go     # Reads and writes notes.json
│   ├── blockers/
│   │   └── storage.go     # Reads and writes blockers.json
│   └── report/
│       └── generator.go   # Builds the standup report
│
├── main.go
├── go.mod
└── README.md
```

---

## Updating devlog

**macOS / Linux:**
```bash
git pull && go install .
```

**Windows (PowerShell):**
```powershell
git pull; go install .
```

---

## Roadmap

- [ ] `devlog summary` — AI-powered natural language standup via local Python service
- [ ] `devlog export` — Export standup report as Markdown or plain text file
- [ ] `devlog config` — Set default author, team name, timezone
- [ ] Slack integration — Post standup directly to a Slack channel

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss what you would
like to change.

---

## License

MIT