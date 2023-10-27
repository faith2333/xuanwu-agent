package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"xuanwu-agent/internal/tsagents"
)

var tsAgentsServer *tsagents.Server
var tsAgentServerParams = &tsagents.FlagParams{}

// TSAgent cmd
var tsAgentCmd = &cobra.Command{
	Use:   "tsagent",
	Short: "tsagent server agent for xuanwu",
	Long: `
usage: xuanwu-agent tsagent --listen 0.0.0.0:8000 
`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(tsAgentsServer.Listen(tsAgentServerParams))
	},
}

func init() {
	tsAgentsServer = tsagents.NewServer()
	tsAgentCmd.Flags().StringVarP(&tsAgentServerParams.Address, "address", "a", "127.0.0.1:8080", "The Address and Port which TSAgent server listen at")
}
