# Kicad-init
A simple way to initialize KiCad EDA projects.


## Getting Started
1. Install using directions [below](#installation)
2. Run and Provide new directory name
3. Select which tasks you want
4. let it do it's thing
5. You did it!

## Features
- Creates KiCad Files automatically
- Sets up Git
- Creates GitHub Repos
- Starts JOURNAL.md
- Adds .gitignore

## Installation
On macOS or Linux, install the latest release (verified against the release checksum):

```sh
curl -fsSL https://raw.githubusercontent.com/avycado13/kicad-init/main/install.sh | sh
```

This installs to `~/.local/bin`; add that directory to your `PATH` if needed. To choose another location, download the script and run `INSTALL_DIR=/your/bin sh install.sh`. Windows users can download the binary from [GitHub Releases](https://github.com/avycado13/kicad-init/releases).

The tool also needs `git` and the `gh` CLI for the corresponding project setup tasks.
