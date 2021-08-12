package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/xackery/gamecli/gog"
	"github.com/xackery/gamecli/steam"
)

var (
	name     string
	isDirect bool
	appID    string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a game",
	Long:  `Run a game on a preferred platform, either appid or client is required`,
	Run: func(cmd *cobra.Command, args []string) {
		rootErrorHandler(run(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVar(&appID, "appid", "", "application id of a game to run")
	runCmd.PersistentFlags().StringVar(&name, "name", "", "name of game to run")
	runCmd.PersistentFlags().BoolVar(&isDirect, "direct", true, "run the game without the client")
	runCmd.PersistentFlags().StringVar(&client, "client", "", fmt.Sprintf("which cloud client to use, options: %s", strings.Join(validClients, " ")))
}

func run(cmd *cobra.Command, args []string) error {
	err := validateAppID()
	if err != nil {
		return fmt.Errorf("validateAppID: %w", err)
	}
	clients := []string{}
	if client == "" {
		clients = validClients
	} else {
		clients = append(clients, client)
	}
	for _, c := range clients {
		log.Info().Msgf("running %s", c)
		ctx := context.Background()
		switch c {
		case "steam":
			s, err := steam.New(ctx)
			if err != nil {
				return fmt.Errorf("steam.New: %w", err)
			}
			err = s.Run(ctx, appID, name, isDirect)
			if err != nil {
				return fmt.Errorf("steam.Run: %w", err)
			}
			return nil
		case "gog":
			s, err := gog.New(ctx)
			if err != nil {
				return fmt.Errorf("gog.New: %w", err)
			}
			err = s.Run(ctx, appID)
			if err != nil {
				return fmt.Errorf("gog.Run: %w", err)
			}
			return nil
		default:
			return fmt.Errorf("%s is unsupported", c)
		}
	}
	return nil
}

func validateAppID() error {
	if appID == "" {
		return nil
	}
	client = strings.ToLower(client)

	switch client {
	case "gog", "steam":
		return nil
	default:
		return fmt.Errorf("unknown client: %s, (valid options include %s)", client, strings.Join(validClients, " "))
	}
}
