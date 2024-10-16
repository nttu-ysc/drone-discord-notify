package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	webhook = RootCmd.PersistentFlags().StringP("webhook", "w", os.Getenv("PLUGIN_DISCORD_WEBHOOK_URL"), "discord webhook URL")
	content = RootCmd.PersistentFlags().StringP("content", "c", os.Getenv("PLUGIN_DISCORD_CONTENT"), "discord webhook content")

	var emojiCircle string = ":green_circle:"
	var emojiStatus string = ":white_check_mark:"
	if os.Getenv("DRONE_STAGE_STATUS") != "success" {
		emojiCircle = ":red_circle:"
		emojiStatus = ":x:"
	}

	DefaultEmbed = Embeds{
		Title:       os.Getenv("DRONE_REPO"),
		Description: fmt.Sprintf("%s **Build #%s %s**", emojiCircle, os.Getenv("DRONE_BUILD_NUMBER"), os.Getenv("DRONE_COMMIT_MESSAGE")),
		Url:         os.Getenv("DRONE_BUILD_LINK"),
		Color:       3066993,
		Fields: []EmbedsField{
			{
				Name:   "Branch",
				Value:  os.Getenv("DRONE_COMMIT_BRANCH"),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  fmt.Sprintf("%s %s", emojiStatus, os.Getenv("DRONE_STAGE_STATUS")),
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  os.Getenv("DRONE_COMMIT_AUTHOR"),
				Inline: true,
			},
			{
				Name:   "Event",
				Value:  os.Getenv("DRONE_BUILD_EVENT"),
				Inline: true,
			},
			{
				Name:   "Changes",
				Value:  fmt.Sprintf("[View Changes](%s)", os.Getenv("DRONE_COMMIT_LINK")),
				Inline: false,
			},
			{
				Name:   "Current time",
				Value:  time.Now().Local().Format("2006-01-02 15:04:05 -07:00"),
				Inline: false,
			},
		},
		Footer: Footer{
			Text:    fmt.Sprintf("Build #%s by %s", os.Getenv("DRONE_BUILD_NUMBER"), os.Getenv("DRONE_COMMIT_AUTHOR")),
			IconUrl: os.Getenv("DRONE_COMMIT_AUTHOR_AVATAR"),
		},
	}
}

type DiscordMessage struct {
	Content string   `json:"content"`
	Embeds  []Embeds `json:"embeds"`
}

type Embeds struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Url         string        `json:"url"`
	Color       int           `json:"color"`
	Fields      []EmbedsField `json:"fields"`
	Footer      Footer        `json:"footer"`
}

type EmbedsField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Footer struct {
	Text    string `json:"text"`
	IconUrl string `json:"icon_url"`
}

var DefaultEmbed Embeds

func SendDiscordMessage() {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(DiscordMessage{Content: *content, Embeds: []Embeds{DefaultEmbed}})

	client := new(http.Client)
	client.Timeout = time.Second * 10

	req, err := http.NewRequest(http.MethodPost, *webhook, body)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		log.Printf("Response Status: %d\nResponse Body: %s\n", resp.StatusCode, respBody)
		return
	}

	log.Println("Success!")
}
