package cmd

import (
	"fmt"
	"syscall"

	portainer "github.com/cbridges1/portainer-cli/pkg/client"
	"github.com/cbridges1/portainer-cli/pkg/storage"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Portainer",
	Long:  `Login to Portainer using username/password or access code`,
}

var loginPasswordCmd = &cobra.Command{
	Use:   "password",
	Short: "Login with username and password",
	Long:  `Authenticate with Portainer using username and password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := storage.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if config.URL == "" {
			return fmt.Errorf("Portainer URL not configured. Use 'portainer-cli config set url <url>' first")
		}

		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			fmt.Print("Username: ")
			fmt.Scanln(&username)
		}

		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		password := string(passwordBytes)
		fmt.Println() // New line after password input

		// Create auth client
		authClient := portainer.NewAuthClient(config.URL)

		// Authenticate
		authResp, err := authClient.AuthenticateUser(username, password)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Parse token expiry
		expiresAt, err := portainer.ParseJWTExpiry(authResp.JWT)
		if err != nil {
			return fmt.Errorf("failed to parse JWT expiry: %w", err)
		}

		// Store credentials and token
		if err := config.SetCredentials(username, password, "password"); err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}

		if err := config.SetToken(authResp.JWT, expiresAt); err != nil {
			return fmt.Errorf("failed to store token: %w", err)
		}

		fmt.Println("Login successful!")
		fmt.Printf("Token expires at: %s\n", expiresAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

var loginCodeCmd = &cobra.Command{
	Use:   "code",
	Short: "Login with access code",
	Long:  `Authenticate with Portainer using an access code`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := storage.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if config.URL == "" {
			return fmt.Errorf("Portainer URL not configured. Use 'portainer-cli config set url <url>' first")
		}

		accessCode, _ := cmd.Flags().GetString("code")
		if accessCode == "" {
			fmt.Print("Access Code: ")
			fmt.Scanln(&accessCode)
		}

		// Create auth client
		authClient := portainer.NewAuthClient(config.URL)

		// Authenticate
		authResp, err := authClient.AuthenticateWithCode(accessCode)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		// Parse token expiry
		expiresAt, err := portainer.ParseJWTExpiry(authResp.JWT)
		if err != nil {
			return fmt.Errorf("failed to parse JWT expiry: %w", err)
		}

		// Store credentials and token (use access code as "password")
		if err := config.SetCredentials("", accessCode, "code"); err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}

		if err := config.SetToken(authResp.JWT, expiresAt); err != nil {
			return fmt.Errorf("failed to store token: %w", err)
		}

		fmt.Println("Login successful!")
		fmt.Printf("Token expires at: %s\n", expiresAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear stored credentials",
	Long:  `Remove stored credentials and JWT token from local storage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := storage.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := config.ClearCredentials(); err != nil {
			return fmt.Errorf("failed to clear credentials: %w", err)
		}

		fmt.Println("Logged out successfully. Credentials and tokens have been cleared.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)

	loginCmd.AddCommand(loginPasswordCmd)
	loginCmd.AddCommand(loginCodeCmd)

	loginPasswordCmd.Flags().StringP("username", "u", "", "Username for authentication")
	loginCodeCmd.Flags().StringP("code", "c", "", "Access code for authentication")
}
