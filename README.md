# git-monitor
Monitor Git repositories for new commits.

On each run the remote is fetched and compared to the local state. If anything
has changed (new commits) an information line for this repository is given.

All repositories are fetched without a local checkout into 
`~/.git-monitor/<repo>`. 

## Usage

Add a new repositor for monitoring:
```shell
git-monitor https://github.com/Kriechi/git-monitor.git
```

Check if there are new commits:
```shell
git-monitor
```
