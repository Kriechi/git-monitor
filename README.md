# git-monitor
Monitor Git repositories for new commits.

On each run the remote is fetched and compared to the local state. If anything
has changed (new commits) an information line for this repository is given.

A common activities of most developers and IT enthusiasts is a daily `apt-get
update` (or similar) to check for new software package updates. Some even
perform a `brew update` or `npm update` to always get the bleeding edge of new
releases. Some software is only available on GitHub, Gitlab, or other git-based
hosting platforms without direct integration into package manager. `git-monitor`
helps you to say on top on things and get notified about new commits and changes
in repositories and their branches.

All repositories are fetched without a local checkout into
`~/.git-monitor/<repo>` (unless `repo_dir` is configured otherwise).

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

# Configuration

You can use this YAML configuration file and save it to `~/.git-monitor.yaml`
(or use `--config` to override):
```yaml
verbose: true
repo_dir: ~/.my_repos_to_track
ignored_branches:
  - requires-io-master
  - some-other-branch-to-ignore
```
