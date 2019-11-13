package main

import (
	"fmt"
	"os"

	"github.com/GameBench/gba-client-go"
	"github.com/spf13/cobra"
)

var (
	server string
	username string
	password string
	client *gba.GbaClient
)

var rootCmd = &cobra.Command{
	Use:   "foo",
	Short: "Foo",
	Long:  "Foo",
	PersistentPreRun: func(cmd *cobra.Command, args[]string) {
		config := &gba.Config{BaseUrl: server, Username: username, Password: password}
		client = gba.New(config)
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Foo")
		devices, err := client.ListDevices()
		if err != nil {
			panic(err)
		}

		fmt.Println(devices)
	},
}

var listDevicesCmd = &cobra.Command{
	Use: "list-devices",
	Short: "list-devices",
	Long: "list-devices",
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := client.ListDevices()
		if err != nil {
			panic(err)
		}

		fmt.Println(devices)
	},
}

var listDeviceAppsCmd = &cobra.Command{
	Use: "list-device-apps [DEVICE ID]",
	Short: "list-device-apps",
	Long: "list-device-apps",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apps, err := client.GetDeviceApps(args[0])
		if err != nil {
			panic(err)
		}

		for _, app := range apps {
			fmt.Println(app)
		}
	},
}

var getDeviceCmd = &cobra.Command{
	Use: "get-device [DEVICE ID]",
	Short: "get-device",
	Long: "get-device",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := client.GetDevice(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Println(devices)
	},
}

var startSessionCmd = &cobra.Command{
	Use: "start-session [DEVICE ID] [APP ID]",
	Short: "start-session",
	Long: "start-session",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := client.StartSession(args[0], args[1], nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(devices)
	},
}

var stopSessionCmd = &cobra.Command{
	Use: "stop-session [SESSION ID]",
	Short: "stop-session",
	Long: "stop-session",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := client.StopSession(args[0])
		if err != nil {
			panic(err)
		}
	},
}

var syncSessionsCmd = &cobra.Command{
	Use: "sync-sessions",
	Short: "sync-sessions",
	Long: "sync-sessions",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.Sync()
		if err != nil {
			panic(err)
		}
	},
}

var getPropertiesCmd = &cobra.Command{
	Use: "get-properties",
	Short: "get-properties",
	Long: "get-properties",
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
	Use: "set-property",
	Short: "set-property",
	Long: "set-property",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.SetProperty(args[0], args[1])
		if err != nil {
			panic(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(listDevicesCmd)
	rootCmd.AddCommand(listDeviceAppsCmd)
	rootCmd.AddCommand(getDeviceCmd)
	rootCmd.AddCommand(startSessionCmd)
	rootCmd.AddCommand(stopSessionCmd)
	rootCmd.AddCommand(syncSessionsCmd)
	rootCmd.AddCommand(getPropertiesCmd)
	rootCmd.AddCommand(setPropertyCmd)

	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
