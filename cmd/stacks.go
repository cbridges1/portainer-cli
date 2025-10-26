package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	portainer "github.com/jalenbridges/portainer-cli/pkg/client"
	"github.com/sanity-io/litter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "Manage Portainer stacks",
	Long:  `Commands for managing Docker Compose stacks in Portainer`,
}

var stacksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stacks",
	Long:  `List all Docker Compose stacks in the Portainer environment`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getPortainerClient()
		if err != nil {
			return err
		}

		stacks, err := client.ListStacks()
		if err != nil {
			return fmt.Errorf("failed to list stacks: %w", err)
		}

		if len(stacks) == 0 {
			fmt.Println("No stacks found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED\tCREATED BY")

		for _, stack := range stacks {
			created := time.Unix(stack.CreationDate, 0).Format("2006-01-02 15:04:05")
			status := getStackStatus(stack.Status)
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				stack.ID, stack.Name, status, created, stack.CreatedBy)
		}

		return w.Flush()
	},
}

func getPortainerClient() (*portainer.Client, error) {
	// Check if legacy flags are provided (API key mode)
	url := viper.GetString("url")
	token := viper.GetString("token")
	endpoint := viper.GetInt("endpoint")

	if url != "" && token != "" && endpoint != 0 {
		// Legacy API key mode
		config := portainer.Config{
			URL:        url,
			Token:      token,
			EndpointID: endpoint,
			AuthType:   "api-key",
		}
		return portainer.NewClient(config), nil
	}

	// Use JWT authentication with stored credentials
	return portainer.NewAuthenticatedClient()
}

func getStackStatus(status int) string {
	switch status {
	case 1:
		return "ACTIVE"
	case 2:
		return "INACTIVE"
	default:
		return "UNKNOWN"
	}
}

var stacksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new stack",
	Long:  `Create a new Docker Compose stack in Portainer from a file, repository, or inline content`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getPortainerClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("stack name is required (use --name flag)")
		}

		file, _ := cmd.Flags().GetString("file")
		content, _ := cmd.Flags().GetString("content")
		repository, _ := cmd.Flags().GetString("repository")
		composePath, _ := cmd.Flags().GetString("compose-path")
		swarmID, _ := cmd.Flags().GetString("swarm-id")

		var stack *portainer.Stack

		if repository != "" {
			if composePath == "" {
				composePath = "docker-compose.yml"
			}

			req := portainer.CreateStackRequest{
				Name:                        name,
				SwarmID:                     swarmID,
				RepositoryURL:               repository,
				ComposeFilePathInRepository: composePath,
			}

			stack, err = client.CreateStackFromRepository(req)
		} else if file != "" {
			fileContent, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", file, err)
			}

			req := portainer.CreateStackRequest{
				Name:             name,
				SwarmID:          swarmID,
				StackFileContent: string(fileContent),
			}

			litter.Dump(req)

			stack, err = client.CreateStackFromString(req)
		} else if content != "" {
			req := portainer.CreateStackRequest{
				Name:             name,
				SwarmID:          swarmID,
				StackFileContent: content,
			}

			stack, err = client.CreateStackFromString(req)
		} else {
			return fmt.Errorf("one of --file, --content, or --repository must be specified")
		}

		if err != nil {
			return fmt.Errorf("failed to create stack: %w", err)
		}

		fmt.Printf("Stack '%s' created successfully with ID: %d\n", stack, stack)
		return nil
	},
}

var stacksUpdateCmd = &cobra.Command{
	Use:   "update <stack-id>",
	Short: "Update an existing stack",
	Long:  `Update an existing Docker Compose stack in Portainer with new content`,
	//Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getPortainerClient()
		if err != nil {
			return err
		}

		stackID, err := cmd.Flags().GetInt("id")
		name, _ := cmd.Flags().GetString("name")

		if err != nil && name == "" {
			return fmt.Errorf("invalid stack ID: %s", args[0])
		}

		if name != "" {
			stack, err := client.FindStackByName(name)
			if err != nil {
				return fmt.Errorf("failed to update stack: %w", err)
			}
			if stack != nil {
				stackID = stack.ID
			}
		}

		file, _ := cmd.Flags().GetString("file")
		content, _ := cmd.Flags().GetString("content")
		prune, _ := cmd.Flags().GetBool("prune")
		pullImage, _ := cmd.Flags().GetBool("pull-image")

		var stackContent string

		if file != "" {
			fileContent, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", file, err)
			}
			stackContent = string(fileContent)
		} else if content != "" {
			stackContent = content
		} else {
			return fmt.Errorf("one of --file or --content must be specified")
		}

		req := portainer.UpdateStackRequest{
			StackFileContent: stackContent,
			Prune:            prune,
			PullImage:        pullImage,
		}

		stack, err := client.UpdateStack(stackID, req)
		if err != nil {
			return fmt.Errorf("failed to update stack: %w", err)
		}

		fmt.Printf("Stack '%s' (ID: %d) updated successfully\n", stack.Name, stack.ID)

		return nil
	},
}

var stacksDeleteCmd = &cobra.Command{
	Use:   "delete <stack-id>",
	Short: "Delete a stack",
	Long:  `Delete a Docker Compose stack from Portainer`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getPortainerClient()
		if err != nil {
			return err
		}

		stackID, err := cmd.Flags().GetInt("id")
		name, _ := cmd.Flags().GetString("name")

		if err != nil && name == "" {
			return fmt.Errorf("invalid stack ID: %s", args[0])
		}

		if name != "" {
			stack, err := client.FindStackByName(name)
			if err != nil {
				return fmt.Errorf("failed to update stack: %w", err)
			}
			if stack != nil {
				stackID = stack.ID
			}
		}

		force, _ := cmd.Flags().GetBool("force")

		if !force {
			stack, err := client.GetStack(stackID)
			if err != nil {
				return fmt.Errorf("failed to get stack info: %w", err)
			}

			fmt.Printf("Are you sure you want to delete stack '%s' (ID: %d)? [y/N]: ", stack.Name, stack.ID)
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
				fmt.Println("Deletion cancelled")
				return nil
			}
		}

		err = client.DeleteStack(stackID)
		if err != nil {
			return fmt.Errorf("failed to delete stack: %w", err)
		}

		fmt.Printf("Stack with ID %d deleted successfully\n", stackID)
		return nil
	},
}

var stacksInspectCmd = &cobra.Command{
	Use:   "inspect <stack-id>",
	Short: "Inspect a stack",
	Long:  `Display detailed information about a Docker Compose stack`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getPortainerClient()
		if err != nil {
			return err
		}

		stackID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid stack ID: %s", args[0])
		}

		stack, err := client.GetStack(stackID)
		if err != nil {
			return fmt.Errorf("failed to get stack: %w", err)
		}

		showFile, _ := cmd.Flags().GetBool("show-file")

		fmt.Printf("ID: %d\n", stack.ID)
		fmt.Printf("Name: %s\n", stack.Name)
		fmt.Printf("Status: %s\n", getStackStatus(stack.Status))
		fmt.Printf("Endpoint ID: %d\n", stack.EndpointID)
		fmt.Printf("Entry Point: %s\n", stack.EntryPoint)

		if stack.SwarmID != "" {
			fmt.Printf("Swarm ID: %s\n", stack.SwarmID)
		}

		if stack.Namespace != "" {
			fmt.Printf("Namespace: %s\n", stack.Namespace)
		}

		created := time.Unix(stack.CreationDate, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("Created: %s\n", created)
		fmt.Printf("Created By: %s\n", stack.CreatedBy)

		if stack.UpdateDate > 0 {
			updated := time.Unix(stack.UpdateDate, 0).Format("2006-01-02 15:04:05")
			fmt.Printf("Updated: %s\n", updated)
			fmt.Printf("Updated By: %s\n", stack.UpdatedBy)
		}

		if len(stack.Env) > 0 {
			fmt.Printf("\nEnvironment Variables:\n")
			for _, env := range stack.Env {
				fmt.Printf("  %s=%s\n", env.Name, env.Value)
			}
		}

		if stack.GitConfig != nil {
			fmt.Printf("\nGit Configuration:\n")
			fmt.Printf("  URL: %s\n", stack.GitConfig.URL)
			fmt.Printf("  Reference: %s\n", stack.GitConfig.ReferenceName)
			fmt.Printf("  Config File Path: %s\n", stack.GitConfig.ConfigFilePath)
		}

		if showFile {
			fmt.Printf("\nStack File Content:\n")
			fmt.Println("---")
			fileContent, err := client.GetStackFile(stackID)
			if err != nil {
				fmt.Printf("Error getting stack file: %v\n", err)
			} else {
				fmt.Println(fileContent)
			}
			fmt.Println("---")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stacksCmd)
	stacksCmd.AddCommand(stacksListCmd)
	stacksCmd.AddCommand(stacksCreateCmd)
	stacksCmd.AddCommand(stacksUpdateCmd)
	stacksCmd.AddCommand(stacksDeleteCmd)
	stacksCmd.AddCommand(stacksInspectCmd)

	stacksCreateCmd.Flags().StringP("name", "n", "", "Name of the stack (required)")
	stacksCreateCmd.Flags().StringP("file", "f", "", "Path to docker-compose.yml file")
	stacksCreateCmd.Flags().StringP("content", "c", "", "Inline docker-compose content")
	stacksCreateCmd.Flags().StringP("repository", "r", "", "Git repository URL")
	stacksCreateCmd.Flags().String("compose-path", "docker-compose.yml", "Path to compose file in repository")
	stacksCreateCmd.Flags().String("swarm-id", "", "Docker Swarm ID (required for Swarm mode)")

	stacksUpdateCmd.Flags().IntP("id", "i", 0, "Id of the stack")
	stacksUpdateCmd.Flags().StringP("name", "n", "", "Name of the stack")
	stacksUpdateCmd.Flags().StringP("file", "f", "", "Path to docker-compose.yml file")
	stacksUpdateCmd.Flags().StringP("content", "c", "", "Inline docker-compose content")
	stacksUpdateCmd.Flags().Bool("prune", false, "Remove services that are no longer defined")
	stacksUpdateCmd.Flags().Bool("pull-image", false, "Pull latest images before updating")

	stacksDeleteCmd.Flags().IntP("id", "i", 0, "Id of the stack")
	stacksDeleteCmd.Flags().StringP("name", "n", "", "Name of the stack")
	stacksDeleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")

	stacksInspectCmd.Flags().IntP("id", "i", 0, "Id of the stack")
	stacksInspectCmd.Flags().StringP("name", "n", "", "Name of the stack")
	stacksInspectCmd.Flags().Bool("show-file", false, "Show the stack file content")
}
