/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

type Data struct {
	R030         float64 `json:'r030'`
	Txt          string  `json:'txt'`
	Rate         float64 `json:'rate'`
	Cc           string  `json:'cc'`
	Exchangedate string  `json:'exchangedate'`
}

var (
	// Teletoken bot
	Teletoken = os.Getenv("TELE_TOKEN")
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started ", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  Teletoken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		kbot.Handle("/start", func(ctx telebot.Context) error {
			menu := &telebot.ReplyMarkup{
				ReplyKeyboard: [][]telebot.ReplyButton{
					{{Text: "Hello"}},
					{{Text: "€ EUR"}, {Text: "$ USD"}},
				},
			}
			return ctx.Send("Welcome to bot!", menu)
		})

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			switch m.Text() {
			case "Hello":
				return m.Send(fmt.Sprintf("Hello I'm Kbot %s! I can find out today's exchange rate. Please, select a currency", appVersion))
			case "€ EUR":
				return m.Send("Euro exchange rate: " + getRate("EUR"))
			case "$ USD":
				return m.Send("Euro exchange rate: " + getRate("USD"))
			default:
				return m.Send("Unknown command. Please try again.")
			}
		})

		kbot.Start()
	},
}

func getRate(valcode string) string {
	url := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%v&json", valcode)
	resp, getErr := http.Get(url)

	if getErr != nil {
		log.Fatal(getErr)
	}

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	defer resp.Body.Close()

	body := []Data{}
	err := json.Unmarshal((respBody), &body)

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("body: ", body[0].Rate)
	return fmt.Sprintf("%v", body[0].Rate)
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
