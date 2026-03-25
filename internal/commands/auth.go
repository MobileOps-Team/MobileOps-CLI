package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MobileOps-Team/mobileops-cli/internal/client"
	"github.com/MobileOps-Team/mobileops-cli/internal/config"
	"github.com/MobileOps-Team/mobileops-cli/internal/envelope"
	"github.com/MobileOps-Team/mobileops-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the MobileOps API",
	RunE:  runAuthLogin,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	RunE:  runAuthLogout,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	RunE:  runAuthStatus,
}

func init() {
	authLoginCmd.Flags().String("host", "", "API host URL (overrides environment default)")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authStatusCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	env := config.ResolveEnvironment(envFlag)
	host, _ := cmd.Flags().GetString("host")
	if host == "" {
		host = config.HostFor(env)
	}

	fmt.Printf("Environment: %s (%s)\n\n", env, host)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Access Key: ")
	accessKey, _ := reader.ReadString('\n')
	accessKey = strings.TrimSpace(accessKey)

	fmt.Print("Secret Key: ")
	secretKey, _ := reader.ReadString('\n')
	secretKey = strings.TrimSpace(secretKey)

	if accessKey == "" || secretKey == "" {
		fmt.Fprintln(os.Stderr, "Error: Access key and secret key are required.")
		os.Exit(1)
	}

	fmt.Print("Verifying credentials... ")

	c := client.NewWithKeys(host, accessKey, secretKey)
	_, err := c.Get("assets", map[string]string{"limit": "1"})
	if err != nil {
		fmt.Println("FAILED")
		if _, ok := err.(*client.AuthError); ok {
			fmt.Fprintln(os.Stderr, "Error: Invalid credentials. Please check your access key and secret key.")
		} else {
			fmt.Fprintf(os.Stderr, "Error: Could not connect to %s (%s)\n", host, err.Error())
		}
		os.Exit(1)
	}

	if err := config.Save(host, accessKey, secretKey, env); err != nil {
		fmt.Println("FAILED")
		fmt.Fprintf(os.Stderr, "Error: Could not save credentials (%s)\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("OK")
	fmt.Println()
	fmt.Println("Authenticated successfully!")
	fmt.Printf("Credentials saved to %s\n", config.Path(env))
	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	env := config.ResolveEnvironment(envFlag)

	if config.Exists(env) {
		if err := config.Delete(env); err != nil {
			return err
		}
		fmt.Printf("Credentials removed for %s.\n", env)
	} else {
		fmt.Printf("No credentials found for %s.\n", env)
	}
	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	env := config.ResolveEnvironment(envFlag)

	if !config.Exists(env) {
		if jsonOutput {
			env := envelope.WrapError(
				fmt.Sprintf("Not authenticated for %s", env),
				[]string{fmt.Sprintf("mobileops auth login --env %s", env)},
			)
			formatter.Output(env, true)
		} else {
			fmt.Printf("Not authenticated for %s. Run `mobileops auth login --env %s` to get started.\n", env, env)
		}
		return nil
	}

	creds, err := config.Load(env)
	if err != nil {
		return err
	}

	keyPreview := creds.AccessKey
	if len(keyPreview) > 8 {
		keyPreview = keyPreview[:8] + "..."
	}

	// Test connectivity
	connected := false
	c, err := client.New("", "", "", env)
	if err == nil {
		_, err = c.Get("assets", map[string]string{"limit": "1"})
		connected = err == nil
	}

	if jsonOutput {
		data := map[string]interface{}{
			"authenticated": true,
			"environment":   env,
			"host":          creds.Host,
			"access_key":    keyPreview,
			"api_connected": connected,
		}

		summary := fmt.Sprintf("Authenticated and connected (%s)", env)
		if !connected {
			summary = fmt.Sprintf("Authenticated but API unreachable (%s)", env)
		}

		env := envelope.Wrap(data, summary, []string{
			"mobileops vessels list",
			"mobileops jobs list",
		})
		formatter.Output(env, true)
	} else {
		status := "Connected"
		if !connected {
			status = "API unreachable"
		}
		fmt.Printf("Environment: %s\n", env)
		fmt.Printf("Status:      %s\n", status)
		fmt.Printf("Host:        %s\n", creds.Host)
		fmt.Printf("Access Key:  %s\n", keyPreview)
		fmt.Printf("Config:      %s\n", config.Path(env))
	}
	return nil
}
