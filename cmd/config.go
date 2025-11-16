package cmd

import (
	"fmt"
	"net/url"

	"github.com/cbridges1/portainer-cli/pkg/storage"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Configure Portainer CLI settings like URL and authentication`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `Set configuration values such as Portainer URL`,
}

var configSetUrlCmd = &cobra.Command{
	Use:   "url <portainer-url>",
	Short: "Set the Portainer URL",
	Long:  `Set the Portainer instance URL (e.g., https://portainer.example.com)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		portainerURL := args[0]

		// Validate URL
		if _, err := url.ParseRequestURI(portainerURL); err != nil {
			return fmt.Errorf("invalid URL format: %w", err)
		}

		config, err := storage.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := config.SetURL(portainerURL); err != nil {
			return fmt.Errorf("failed to save URL: %w", err)
		}

		fmt.Printf("Portainer URL set to: %s\n", portainerURL)
		return nil
	},
}

var configSetEndpointCmd = &cobra.Command{
	Use:   "endpoint <endpoint-id>",
	Short: "Set the Portainer endpoint ID",
	Long:  `Set the default Portainer endpoint ID (e.g., 1 for local Docker)`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		endpointIDStr := args[0]

		// Parse endpoint ID
		endpointID := 0
		if _, err := fmt.Sscanf(endpointIDStr, "%d", &endpointID); err != nil {
			return fmt.Errorf("invalid endpoint ID format: %w", err)
		}

		if endpointID <= 0 {
			return fmt.Errorf("endpoint ID must be a positive integer")
		}

		config, err := storage.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := config.SetEndpoint(endpointID); err != nil {
			return fmt.Errorf("failed to save endpoint ID: %w", err)
		}

		fmt.Printf("Portainer endpoint ID set to: %d\n", endpointID)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get configuration values",
	Long:  `Display current configuration values`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := storage.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Printf("Portainer URL: %s\n", config.URL)
		fmt.Printf("Endpoint ID: %d\n", config.EndpointID)

		if config.Credentials != nil {
			fmt.Printf("Stored credentials: %s (%s)\n", config.Credentials.Username, config.Credentials.Type)
		} else {
			fmt.Println("Stored credentials: None")
		}

		if config.Token != nil {
			if config.IsTokenValid() {
				fmt.Println("Token status: Valid")
			} else {
				fmt.Println("Token status: Expired")
			}
			fmt.Printf("Token expires: %s\n", config.Token.ExpiresAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("Token status: None")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configSetCmd.AddCommand(configSetUrlCmd)
	configSetCmd.AddCommand(configSetEndpointCmd)
}
