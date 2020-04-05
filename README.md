# git-monitor
Monitor Git repositories for new commits.

On each run the remote is fetched and compared to the local state. If anything
has changed (new commits) an information line for this repository is given.

All repositories are fetched without a local checkout into
`~/.git-monitor/<repo>`.

## Install

```console
go get github.com/kriechi/git-monitor
```

## Usage

```
git-monitor manages a list of git repositories and can check them for new
commits or changes on branches. It works against local and remote repositories.

git-monitor is a CLI application which can tell you if a repository has changes
since the last time you checked it. Think of it as "apt-get update" which tells
you that repository X on branch Y has new commits.

Usage:
  git-monitor [flags]
  git-monitor [command]

Available Commands:
  add         Add a new repository with a local clone
  check       Check all monitored repositories for changes, or only the ones passed as argument
  help        Help about any command
  list        List all monitored repositories
  remove      Remove a monitored repository

Flags:
      --config string     config file (default is $HOME/.git-monitor.yaml)
  -h, --help              help for git-monitor
      --repo_dir string   directory where to store local repositories
  -v, --verbose           enable verbose output

Use "git-monitor [command] --help" for more information about a command.
```
