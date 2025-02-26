package cmd

import (
	_ "embed"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"scheduler0/pkg/secrets"
	"strings"
)

var CredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "create, view or modify scheduler0 credentials",
	Long: `
Your credentials is composed of a secret key used to create public keys for your clients and basic http authentication 
details for node to node communication.

Usage:

	credential init

This will go through the credentials init flow.
Note that starting the server without going through the init flow will not work.
`,
}

// InitCmd initializes scheduler0 configuration
var initCCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize credentials for the scheduler0 server",
	Long: `
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr, "[cmd] ", log.LstdFlags)
		logger.Println("Initializing Scheduler0 Credentials")

		scheduler0Secrets := secrets.NewScheduler0Secrets()
		credentials := scheduler0Secrets.GetSecrets()

		if credentials.SecretKey != "" && credentials.AuthUsername != "" && credentials.AuthPassword != "" {
			recreateKey := promptui.Prompt{
				Label:       "Credentials already exist are you sure you want to re-create it[Y/n]:",
				HideEntered: false,
			}
			recreate, _ := recreateKey.Run()

			if strings.ToLower(recreate) == "n" || strings.ToLower(recreate) == "no" {
				return
			}
		}

		secretKeyPrompt := promptui.Prompt{
			Label:       "Secret Key",
			Mask:        '*',
			HideEntered: true,
		}
		SecretKey, _ := secretKeyPrompt.Run()
		credentials.SecretKey = SecretKey

		authUserNamePrompt := promptui.Prompt{
			Label:       "Auth Username",
			HideEntered: true,
		}
		usernameKey, _ := authUserNamePrompt.Run()
		credentials.AuthUsername = usernameKey

		passwordKeyPrompt := promptui.Prompt{
			Label:       "Auth Password",
			Mask:        '*',
			HideEntered: true,
		}
		passwordKey, _ := passwordKeyPrompt.Run()
		credentials.AuthPassword = passwordKey

		scheduler0Secrets.SaveSecrets(credentials)
		logger.Println("Scheduler0 Initialized")
	},
}

// ShowCmd show scheduler0 password configuration
var showCCmd = &cobra.Command{
	Use:   "show",
	Short: "This will show the configurations that have been set.",
	Long: `
Using this credentials you can tell what credentials have been set.

Usage:

	scheduler0 credentials show

Use the --show-password flag if you want the password to be visible.
`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr, "[cmd] ", log.LstdFlags)
		scheduler0Secrets := secrets.NewScheduler0Secrets()
		credentials := scheduler0Secrets.GetSecrets()
		logger.Println("Credentials:")
		logger.Println("SecretKey:", credentials.SecretKey)
		logger.Println("AuthUsername:", credentials.AuthUsername)
		logger.Println("AuthPassword:", credentials.AuthPassword)
	},
}

func init() {
	CredentialCmd.AddCommand(initCCmd)
	CredentialCmd.AddCommand(showCCmd)
}
