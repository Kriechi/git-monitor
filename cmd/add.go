package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new repository with a local clone",
	Long:  `This command adds a repository by making a new shallow clone of the repository.`,
	Args:  cobra.RangeArgs(1, 2),
	Run:   runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	repoDir := viper.GetString("repo_dir")

	url := args[0]

	var targetDir string
	if len(args) == 2 {
		targetDir = args[1]
	} else {
		targetDir = extractHumanishRepoName(url)
	}

	path := filepath.Join(repoDir, targetDir)
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:          url,
		Depth:        1,
		NoCheckout:   true,
		SingleBranch: false,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to clone repository from %s as %s: %v\n", url, targetDir, err)
		os.Exit(1)
	}
	fmt.Printf("Added repository for monitoring from %s as %s.\n", url, targetDir)
}

func extractHumanishRepoName(url string) string {
	// based on https://github.com/git/git/blob/90bbd502d54fe920356fa9278055dc9c9bfe9a56/contrib/examples/git-clone.sh#L231-L232
	targetDir := url // just to get started
	targetDir = strings.TrimSpace(targetDir)
	targetDir = strings.TrimSuffix(targetDir, "/")
	targetDir = regexp.MustCompile(`:*/*\.git`).ReplaceAllString(targetDir, "")
	targetDir = regexp.MustCompile(`.*[/:]`).ReplaceAllString(targetDir, "")
	return targetDir
}
