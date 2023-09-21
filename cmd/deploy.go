package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"xuanwu-agent/internal/deploy"
)

var interDeploy *deploy.Deploy
var flagParams = &deploy.FlagParams{}

// rootCmd represents the base command when called without any subcommands
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy agent for xuanwu",
	Long: `
usage: xuanwu-agent deploy --listen 0.0.0.0:8000 
`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(interDeploy.Listen(flagParams))
	},
}

func init() {
	interDeploy = deploy.NewDeploy()

	deployCmd.Flags().StringVarP(&flagParams.Address, "address", "a", "127.0.0.1:8080", "The Address and Port which deploy listen at")
	deployCmd.Flags().BoolVar(&flagParams.InCluster, "in-cluster", false, "Specify the xuanwu-agent is deploy in a kubernetes cluster")
}
