package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all monitored repositories",
	Long:  `This command prints a list of all currently monitoried repositories.`,
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	repoDir := viper.GetString("repo_dir")
	var repoList [][]string
	repos := getMonitoriedRepositories()
	for _, repo := range repos {
		gitRepo, err := git.PlainOpen(filepath.Join(repoDir, repo))
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: failed to open git repository:", repo, err)
			continue
		}

		remotes, err := gitRepo.Remotes()
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: failed to get remotes from repository:", repo, err)
			continue
		}

		branchList := parseMonitoredBranches(repo)
		branches := strings.Join(branchList, ", ")

		for _, remote := range remotes {
			for _, url := range remote.Config().URLs {
				repoList = append(repoList, []string{repo, url, branches})
			}
		}
	}

	printListResults(repoList)
}

func parseMonitoredBranches(repo string) []string {
	repoDir := viper.GetString("repo_dir")

	var branchList []string
	branchFile := filepath.Join(repoDir, repo, ".git-monitor-branches")
	if _, err := os.Stat(branchFile); err == nil {
		byts, err := ioutil.ReadFile(branchFile)
		if err == nil {
			content := strings.Replace(string(byts), "\r\n", "\n", -1)
			for _, branch := range strings.Split(content, "\n") {
				branch = strings.TrimSpace(branch)
				if len(branch) > 0 {
					branchList = append(branchList, branch)
				}
			}
		}
	} else {
		branchList = []string{"master"}
	}
	return branchList
}

func printListResults(repoList [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repository", "Remote URL", "Monitored Branches"})
	table.SetBorder(false)
	table.AppendBulk(repoList)
	table.Render()
}

func getMonitoriedRepositories() []string {
	var repos []string
	repoDir := viper.GetString("repo_dir")
	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != repoDir {
			repos = append(repos, info.Name())
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: failed to list repositories:", err)
		os.Exit(1)
	}
	return repos
}
