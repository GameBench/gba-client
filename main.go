package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/GameBench/gba-client-go"
	"github.com/spf13/cobra"
)

var (
	server string
	client *gba.GbaClient
)

var rootCmd = &cobra.Command{
	Use:   "gba-cli",
	Short: "gba-cli",
	Long:  "gba-cli",
	PersistentPreRun: func(cmd *cobra.Command, args[]string) {
		config := &gba.Config{BaseUrl: server}
		client = gba.New(config)
	},
}

var listDevicesCmd = &cobra.Command{
	Use: "list",
	Short: "List devices",
	Long: "List devices",
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := client.ListDevices()
		if err != nil {
			panic(err)
		}

		for _, device := range devices {
			fmt.Printf("%s\t%s\n", device.Id, device.Name)
		}
	},
}

var listDeviceAppsCmd = &cobra.Command{
	Use: "list-apps [DEVICE ID]",
	Short: "List apps",
	Long: "List apps",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apps, err := client.GetDeviceApps(args[0])
		if err != nil {
			panic(err)
		}

		for _, app := range apps {
			fmt.Printf("%s\n", app.Identifier)
		}
	},
}

var getDeviceCmd = &cobra.Command{
	Use: "describe [DEVICE ID]",
	Short: "Describe device",
	Long: "Describe device",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		device, err := client.GetDevice(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\t%s\n", device.Id, device.Name)
	},
}

var startSessionCmd = &cobra.Command{
	Use: "start [DEVICE ID] [APP ID]",
	Short: "Start recording a session",
	Long: "Start recording a session",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		screenshots, _ := cmd.Flags().GetBool("screenshots")
		autoSync, _ := cmd.Flags().GetBool("auto-sync")
		options := &gba.StartSessionOptions{Screenshots: screenshots, AutoSync: autoSync}
		session, err := client.StartSession(args[0], args[1], options)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", session.Id)
	},
}

var stopSessionCmd = &cobra.Command{
	Use: "stop [SESSION ID]",
	Short: "Stop a session",
	Long: "stop",
	Args: cobra.OnlyValidArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")
		if all {
			sessions, err := client.ListSessions()
			if err != nil {
				panic(err)
			}

			for _, session := range sessions {
				err = client.StopSession(session.Id)
				if err != nil {
					panic(err)
				}
			}

			return nil
		}

		if len(args) < 1 {
			return errors.New("requires a SESSION ID argument")
		}

		err = client.StopSession(args[0])
		if err != nil {
			panic(err)
		}

		return nil
	},
}

var syncSessionsCmd = &cobra.Command{
	Use: "sync",
	Short: "Sync all sessions",
	Long: "Sync all sessions",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.Sync()
		if err != nil {
			panic(err)
		}
	},
}

var listSessionsCmd = &cobra.Command{
	Use: "list",
	Short: "List sessions",
	Long: "List sessions",
	Run: func(cmd *cobra.Command, args []string) {
		sessions, err := client.ListSessions()
		if err != nil {
			panic(err)
		}

		for _, session := range sessions {
			fmt.Printf("%s\n", session.Id)
		}
	},
}

var getPropertiesCmd = &cobra.Command{
	Use: "list",
	Short: "List config properties",
	Long: "List config properties",
	Run: func(cmd *cobra.Command, args []string) {
		properties, err := client.GetProperties()
		if err != nil {
			panic(err)
		}

		for key, value := range properties {
			fmt.Println(key, value)
		}

	},
}

var setPropertyCmd = &cobra.Command{
	Use: "set",
	Short: "Set a config property",
	Long: "Set a config property",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.SetProperty(args[0], args[1])
		if err != nil {
			panic(err)
		}
	},
}

var sessionCmd = &cobra.Command{
	Use: "session",
	Short: "session",
	Long: "session",
}

var deviceCmd = &cobra.Command{
	Use: "device",
	Short: "device",
	Long: "device",
}

var propertyCmd = &cobra.Command{
	Use: "property",
	Short: "property",
	Long: "property",
}

func main() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(startSessionCmd)
	startSessionCmd.Flags().Bool("auto-sync", false, "Automatically sync session after it's stopped")
	startSessionCmd.Flags().Bool("screenshots", false, "Take screenshots during session")
	sessionCmd.AddCommand(stopSessionCmd)
	stopSessionCmd.Flags().Bool("all", false, "Stop all sessions")
	sessionCmd.AddCommand(syncSessionsCmd)
	sessionCmd.AddCommand(listSessionsCmd)

	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(listDevicesCmd)
	deviceCmd.AddCommand(listDeviceAppsCmd)
	deviceCmd.AddCommand(getDeviceCmd)

	rootCmd.AddCommand(propertyCmd)
	propertyCmd.AddCommand(getPropertiesCmd)
	propertyCmd.AddCommand(setPropertyCmd)

	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
