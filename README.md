# devlog

A lightweight developer productivity CLI that automatically generates standup
reports from your Git commits and personal notes — all stored locally, never
sent anywhere.

---

## Features

- `devlog today` — see what you committed in the last 24 hours
- `devlog note` — save a quick work note
- `devlog plan` — save what you'll work on today
- `devlog blocker` — log, list and resolve blockers
- `devlog standup` — generate a standup report as text, Markdown, Slack or JSON, and copy it to your clipboard
- `devlog clear` — reset your notes, plans and/or blockers

Only **your** commits are included (matched by your git `user.email`), merge
commits are skipped, and on Mondays the report covers Friday.

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
  clear       Clear stored notes, plans and/or blockers
  help        Help about any command
  note        Save a manual work note
  plan        Save what you plan to work on (the standup's Today section)
  standup     Generate a daily standup report from commits and notes
  today       Show today's git commits
```

---

## Command Reference

A quick overview of everything devlog can do. Examples for each command are in
[Usage](#usage) below.

### Commands

| Command | What it does |
|---|---|
| `devlog today` | Show your commits from the last 24 hours |
| `devlog note <message>` | Save a note about something you did |
| `devlog plan <message>` | Save something you plan to work on today |
| `devlog plan` | List today's plans |
| `devlog blocker <message>` | Save a blocker |
| `devlog blocker list` | List open blockers with their numbers (same as `devlog blocker`) |
| `devlog blocker resolve <number>...` | Mark one or more blockers as resolved |
| `devlog standup` | Print your standup report |
| `devlog clear` | Delete saved notes, plans and/or blockers |
| `devlog help <command>` | Show help for any command |

Quotes around messages are optional: `devlog note fixed the login bug` works.

### Options

**`devlog standup`**

| Option | Short | What it does |
|---|---|---|
| `--format <name>` | `-f` | Output format: `text` (default), `markdown`, `slack` or `json` |
| `--copy` | `-c` | Also copy the report to your clipboard |
| `--author <pattern>` | | Show commits by someone else (name or email) |
| `--all-authors` | | Show commits from everyone, not just you |

**`devlog today`**

| Option | What it does |
|---|---|
| `--author <pattern>` | Show commits by someone else (name or email) |
| `--all-authors` | Show commits from everyone, not just you |

**`devlog plan`**

| Option | Short | What it does |
|---|---|---|
| `--tomorrow` | `-t` | Plan for (or list) the next workday instead of today. On Friday, that's Monday |

**`devlog blocker resolve`**

| Alias | Example |
|---|---|
| `done`, `rm` | `devlog blocker done 2` |

`devlog blocker list` can also be shortened to `devlog blocker ls`.

**`devlog clear`** (pick at least one)

| Option | What it deletes |
|---|---|
| `--notes` | All notes |
| `--plans` | All plans |
| `--blockers` | All blockers, open and resolved |
| `--all` | Everything above |

### How the standup report is built

| Section | Contents |
|---|---|
| **Yesterday** | Your commits and notes since the start of the previous workday (on Monday, from Friday) |
| **Today** | Your plans for today |
| **Blockers** | Every blocker you haven't resolved |

---

## Usage

### Show today's commits

Run this inside any Git repository:
```bash
devlog today
```

Only your own commits are shown. Use `--all-authors` to include everyone, or
`--author <pattern>` to pick someone else.

Output:
```
Today:
  • Fix authentication bug
  • Add Redis cache layer
  • Refactor payment service
```

---

### Save a note
Notes are things you did that don't show up as commits. Quotes are optional.
```bash
devlog note Investigated Redis timeout issue
```

Output:
```
Note saved.
```

---

### Plan your day
Plans fill the **Today** section of your standup.
```bash
devlog plan Finish Redis retry logic
devlog plan --tomorrow Review billing PR   # for the next workday (Friday → Monday)
devlog plan                                # list today's plans
```

---

### Log and resolve blockers
Blockers stay in every standup until you resolve them.
```bash
devlog blocker Waiting for AWS credentials from DevOps
devlog blocker list        # or just: devlog blocker
devlog blocker resolve 1
```

Output of `devlog blocker list`:
```
Open blockers:
  1. Waiting for AWS credentials from DevOps  (since Mon Nov 23)

Resolve one with: devlog blocker resolve <number>
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
  • Investigated Redis timeout issue

Today:
  • Finish Redis retry logic

Blockers:
  • Waiting for AWS credentials from DevOps
```

- **Yesterday** — your commits and notes since the start of the previous
  workday (on Monday, that's Friday onwards).
- **Today** — your plans for today.
- **Blockers** — every blocker you haven't resolved.

Run it outside a git repository and you still get your notes, plans and
blockers.

#### Paste it into Slack, Teams or a doc
```bash
devlog standup --format slack --copy    # short: devlog standup -f slack -c
```

| `--format` | Output |
|---|---|
| `text` (default) | The plain report above |
| `markdown` | `**Yesterday**` headings and `- ` lists |
| `slack` | `*Yesterday*` headings and `•` bullets |
| `json` | `{"yesterday": [...], "today": [...], "blockers": [...]}` |

`--copy` needs `pbcopy` (built into macOS), PowerShell (built into Windows), or
`wl-copy`, `xclip` or `xsel` on Linux.

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

Clear only plans:
```bash
devlog clear --plans
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
| `plans.json` | Your saved plans |
| `blockers.json` | Your blockers, including resolved ones |

To inspect your stored data at any time:

**macOS / Linux:**
```bash
cat ~/.devlog/notes.json
cat ~/.devlog/plans.json
cat ~/.devlog/blockers.json
```

**Windows (PowerShell):**
```powershell
type $env:USERPROFILE\.devlog\notes.json
type $env:USERPROFILE\.devlog\plans.json
type $env:USERPROFILE\.devlog\blockers.json
```

---

## Project Structure
```
devlog/
├── cmd/
│   ├── root.go        # Root command
│   ├── today.go       # devlog today
│   ├── gitopts.go     # --author / --all-authors flags
│   ├── note.go        # devlog note
│   ├── plan.go        # devlog plan
│   ├── blocker.go     # devlog blocker (list, resolve)
│   ├── standup.go     # devlog standup
│   └── clear.go       # devlog clear
│
├── internal/
│   ├── activity/
│   │   ├── model.go       # Unified Activity type
│   │   └── from_git.go    # Converts git commits to Activity
│   ├── clipboard/
│   │   └── clipboard.go   # Copies text with the OS clipboard tool
│   ├── git/
│   │   └── collector.go   # Runs git log and parses commits
│   ├── store/
│   │   └── store.go       # Atomic JSON read/write in ~/.devlog
│   ├── notes/
│   │   └── storage.go     # notes.json
│   ├── plans/
│   │   └── storage.go     # plans.json
│   ├── blockers/
│   │   └── storage.go     # blockers.json (open / resolved)
│   ├── workday/
│   │   └── workday.go     # Previous / next workday dates
│   └── report/
│       ├── generator.go   # Builds the standup report
│       └── render.go      # text / markdown / slack / json output
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
- [x] Markdown / Slack / JSON output and clipboard copy (`devlog standup -f slack -c`)
- [ ] Standups across multiple repositories
- [ ] `devlog config` — Set default author, team name, timezone
- [ ] Slack integration — Post standup directly to a Slack channel

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss what you would
like to change.

Run the tests before opening a PR:
```bash
go test ./...
```

---

## License

MIT