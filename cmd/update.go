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
	client       string
	validClients = []string{"gog", "steam"}
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update cache of games",
	Long: `A local game cache is kept for quick access of games currently known. 
Running update refreshes the cache, ensuring it is fully in sync with the cloud version`,
	Run: func(cmd *cobra.Command, args []string) {
		rootErrorHandler(update(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVar(&client, "client", "", fmt.Sprintf("which cloud client to use, options: %s", strings.Join(validClients, " ")))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func update(cmd *cobra.Command, args []string) error {
	err := validateClient()
	if err != nil {
		return fmt.Errorf("validateClient: %w", err)
	}
	clients := []string{}
	if client == "" {
		clients = validClients
	} else {
		clients = append(clients, client)
	}
	for _, c := range clients {
		log.Info().Msgf("updating %s", c)
		ctx := context.Background()
		switch c {
		case "steam":
			s, err := steam.New(ctx)
			if err != nil {
				return fmt.Errorf("steam.New: %w", err)
			}
			err = s.Update(ctx)
			if err != nil {
				return fmt.Errorf("steam.Update: %w", err)
			}
		case "gog":
			s, err := gog.New(ctx)
			if err != nil {
				return fmt.Errorf("gog.New: %w", err)
			}
			err = s.Update(ctx)
			if err != nil {
				return fmt.Errorf("gog.Update: %w", err)
			}
		default:
			return fmt.Errorf("%s is unsupported", c)
		}
	}
	return nil
}

func validateClient() error {
	if client == "" {
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
