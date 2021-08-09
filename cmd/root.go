package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command = newRootCmd()

func newRootCmd() *cobra.Command {
	r := &cobra.Command{
		Use:   "git-monitor",
		Short: "git-monitor keeps track of changes in monitored git repositories",
		Long: `
git-monitor manages a list of git repositories and can check them for new
commits or changes on branches. It works against local and remote repositories.

git-monitor is a CLI application which can tell you if a repository has changes
since the last time you checked it. Think of it as "apt-get update" which tells
you that repository X on branch Y has new commits.`,
		Run: runRoot,
	}

	r.AddCommand(addCmd)
	r.AddCommand(checkCmd)
	r.AddCommand(listCmd)
	r.AddCommand(removeCmd)

	return r
}

// Execute of the root command is our main entry point.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-monitor.yaml)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().String("repo_dir", "", "directory where to store local repositories")
	viper.BindPFlag("repo_dir", rootCmd.PersistentFlags().Lookup("repo_dir"))
	viper.SetDefault("repo_dir", "~/.git-monitor")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".git-monitor" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".git-monitor")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	repoDir := viper.GetString("repo_dir")
	if len(repoDir) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: repo_dir not set")
		os.Exit(1)
	}

	repoDir, err := homedir.Expand(repoDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: homedir expand failed for repo_dir %s: %v\n", repoDir, err)
		os.Exit(1)
	}
	repoDir, err = filepath.Abs(repoDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: absolute path failed for repo_dir %s: %v\n", repoDir, err)
		os.Exit(1)
	}
	repoDir, err = filepath.EvalSymlinks(repoDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: eval-symlink failed for repo_dir %s: %v\n", repoDir, err)
		os.Exit(1)
	}
	viper.Set("repo_dir", repoDir)

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		err := os.MkdirAll(repoDir, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: could not create missing repo_dir under %s: %v\n", repoDir, err)
			os.Exit(1)
		}

		if viper.GetBool("verbose") {
			fmt.Println("Created empty repo_dir:", repoDir)
		}
	}
	if viper.GetBool("verbose") {
		fmt.Println("Using repo_dir:", repoDir)
	}

	if viper.GetBool("verbose") {
		ignoredBranches := viper.GetStringSlice("ignored_branches")
		if len(ignoredBranches) > 0 {
			fmt.Println("Ignoring all changes on branches:", strings.Join(ignoredBranches, ", "))
		}
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	checkCmd.Run(cmd, args)
}
