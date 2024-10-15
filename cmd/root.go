package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "v0.0.1"

var webhook *string
var content *string
var RootCmd = &cobra.Command{
	Use: "drone-discord-notify",
	Long: `
 ____                         ____  _                       _ 
|  _ \ _ __ ___  _ __   ___  |  _ \(_)___  ___ ___  _ __ __| |
| | | | '__/ _ \| '_ \ / _ \ | | | | / __|/ __/ _ \| '__/ _  |
| |_| | | | (_) | | | |  __/ | |_| | \__ \ (_| (_) | | | (_| |
|____/|_|  \___/|_| |_|\___| |____/|_|___/\___\___/|_|  \__,_|
                                                              
 _   _       _   _  __       
| \ | | ___ | |_(_)/ _|_   _ 
|  \| |/ _ \| __| | |_| | | |
| |\  | (_) | |_| |  _| |_| |
|_| \_|\___/ \__|_|_|  \__, |
                       |___/ 

Author: Shun
Repository: https://github.com/nttu-ysc/drone-discord-notify
`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		if *webhook == "" {
			cmd.Help()
			return
		}
		SendDiscordMessage()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	webhook = RootCmd.PersistentFlags().StringP("webhook", "w", os.Getenv("PLUGIN_DISCORD_WEBHOOK_URL"), "discord webhook URL")
	content = RootCmd.PersistentFlags().StringP("content", "c", "", "discord webhook content")
}

func SendDiscordMessage() {
	// TODO: Implement SendDiscordMessage
	fmt.Println("TODO: Implement SendDiscordMessage")
}
