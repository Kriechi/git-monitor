package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check all monitored repositories for changes, or only the ones passed as argument",
	Long:  `This command checks every repository for changes by updating its git remote.`,
	Run:   runCheck,
}

func runCheck(cmd *cobra.Command, args []string) {
	var repos []string
	repoDir := viper.GetString("repo_dir")

	if len(args) == 0 {
		repos = getMonitoriedRepositories()
	} else {
		repos = append(repos, args...)
	}

	reposToCheck := make(chan RepoToCheck, len(repos))
	results := make(chan RepoToCheck, len(repos))

	numWorkers := runtime.NumCPU() * 16 // we mostly wait for I/O, so we can safely overcommit on CPU resources
	for w := 1; w <= numWorkers; w++ {
		go worker(reposToCheck, results)
	}
	for _, repo := range repos {
		reposToCheck <- RepoToCheck{
			Name: repo,
			Path: filepath.Join(repoDir, repo),
		}
	}
	close(reposToCheck)

	resultList := formatResults(repos, results)
	printCheckResults(resultList)
}

func formatResults(repos []string, results chan RepoToCheck) [][]string {
	var resultList [][]string
	for a := 1; a <= len(repos); a++ {
		result := <-results

		entries, err := formatResult(result)
		if err != nil {
			continue
		}
		resultList = append(resultList, entries...)
	}
	return resultList
}

func formatResult(result RepoToCheck) ([][]string, error) {
	var resultList [][]string
	if len(result.ChangedBranches) > 0 {
		gitRepo, err := git.PlainOpen(result.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: failed to open git repository:", result.Name, err)
			return nil, errors.New("failed to open git repository")
		}
		remotes, err := gitRepo.Remotes()
		if err != nil || len(remotes) == 0 {
			fmt.Fprintln(os.Stderr, "ERROR: failed to get remotes from git repository:", result.Name, err)
			return nil, errors.New("failed to get remotes from git repository")
		}

		baseURL := remotes[0].Config().URLs[0]

		ignoredBranches := viper.GetStringSlice("ignored_branches")
		for _, branch := range result.ChangedBranches {
			skip := false
			for _, b := range ignoredBranches {
				if branch == b {
					skip = true
					break
				}
			}
			if skip {
				continue
			}

			url := baseURL
			if strings.Contains(url, "github.com") {
				// https://github.com/<username>/<reponame>.git
				url = strings.TrimSuffix(url, ".git") + "/commits/" + branch
			} else if strings.Contains(url, "gitlab.com") {
				// https://gitlab.com/<username>/<reponame>/-/commits/master
				url = strings.TrimSuffix(url, ".git") + "/-/commits/" + branch
			} else if strings.Contains(url, "bitbucket.com") {
				// https://bitbucket.org/<username>/<reponame>/commits/branch/hg-crew-tip
				url = strings.TrimSuffix(url, ".git") + "/commits/branch/" + branch
			}
			if strings.HasPrefix(url, "git@") {
				// git@github.com:<username>/<reponame>.git
				url = strings.Replace(url, ":", "/", 1)
				url = "https://" + strings.TrimPrefix(url, "git@")
			}

			resultList = append(resultList, []string{result.Name, branch, url})
		}
	}
	return resultList, nil
}

func printCheckResults(resultList [][]string) {
	if len(resultList) == 0 {
		fmt.Println("Already up-to-date.")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Repository", "Changed Branch", "Remote URL"})
		table.SetBorder(false)
		table.AppendBulk(resultList)
		table.Render()
	}
}

// RepoToCheck contains input and output data
type RepoToCheck struct {
	Name            string
	Path            string
	ChangedBranches []string
}

func worker(reposToCheck <-chan RepoToCheck, results chan<- RepoToCheck) {
	for repo := range reposToCheck {
		results <- checkRepository(repo)
	}
}

func checkRepository(repo RepoToCheck) RepoToCheck {
	repo.ChangedBranches = make([]string, 0)

	cmd := exec.Command("git", "remote", "-v", "update")
	cmd.Dir = repo.Path
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: failed to check repository:", repo, err)
		return repo
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		// we are only interested in actually changed references:
		// e.g.:     ce5196f..9d842f5  master     -> origin/master
		r := regexp.MustCompile(`\s+[[:alnum:]]+\.\.[[:alnum:]]+\s+(?P<Branch>[[:graph:]]+)\s+`)
		m := r.FindStringSubmatch(line)
		if len(m) == 2 {
			repo.ChangedBranches = append(repo.ChangedBranches, m[1])
		}
	}
	return repo
}

/*
[0] ➜  hugo-academic git:(master) ✗ git remote -v update 2>&1
Fetching origin
POST git-upload-pack (122 bytes)
POST git-upload-pack (gzip 1248 to 677 bytes)
remote: Enumerating objects: 29, done.
remote: Counting objects: 100% (26/26), done.
remote: Compressing objects: 100% (1/1), done.
remote: Total 5 (delta 4), reused 5 (delta 4), pack-reused 0
Unpacking objects: 100% (5/5), 837 bytes | 93.00 KiB/s, done.
From https://github.com/gcushen/hugo-academic
   ce5196f..9d842f5  master     -> origin/master
[0] ➜  hugo-academic git:(master) ✗ git remote -v update 2>&1
Fetching origin
POST git-upload-pack (122 bytes)
From https://github.com/gcushen/hugo-academic
 = [up to date]      master     -> origin/master
*/
