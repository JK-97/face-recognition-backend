package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.jiangxingai.com/luyor/tf-fence-backend/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of tf-fence-backend",
	Long:  `All software has versions. This is tf-fence-backend`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
