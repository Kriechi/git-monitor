package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

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

func init() {
	rootCmd.AddCommand(checkCmd)
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

	var resultList [][]string
	for a := 1; a <= len(repos); a++ {
		result := <-results
		if len(result.ChangedBranches) > 0 {
			for _, branch := range result.ChangedBranches {
				resultList = append(resultList, []string{result.Name, branch, "http://"})
			}
		}
	}

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