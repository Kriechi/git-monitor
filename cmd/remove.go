package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete"},
	Short:   "Remove a monitored repository",
	Long: `
This command removes a repository from the list of monitored repositories.
It deletes the local checkout and cleans up any references to this repository.`,
	Args: cobra.ExactArgs(1),
	Run:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {
	repo := args[0]
	repoDir := viper.GetString("repo_dir")
	path := filepath.Join(repoDir, repo)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "ERROR: repository not found:", repo)
		os.Exit(1)
	}

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: failed to remove repository:", repo, err)
		os.Exit(1)
	}

	fmt.Println("Removed repository", repo, "- it will no longer be monitored.")
}
