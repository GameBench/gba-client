package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/GameBench/gba-client-go"
	"github.com/spf13/cobra"
)

var (
	server string
	client *gba.GbaClient
	commitHash string
	version string
)

var rootCmd = &cobra.Command{
	Use:   "gba-cli",
	Short: "gba-cli",
	Long:  "gba-cli",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config := &gba.Config{BaseUrl: server}
		client = gba.New(config)
	},
}

var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List devices",
	Long:  "List devices",
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
	Use:   "list-apps [DEVICE ID]",
	Short: "List apps",
	Long:  "List apps",
	Args:  cobra.MinimumNArgs(1),
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
	Use:   "describe [DEVICE ID]",
	Short: "Describe device",
	Long:  "Describe device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		device, err := client.GetDevice(args[0])
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\t%s\n", device.Id, device.Name)
	},
}

var shellCmd = &cobra.Command{
	Use: "shell [DEVICE ID] -- [COMMAND]",
	Short: "Execute shell command on Android device",
	Long:  `Execute shell command on Android device

List packages

    gba-client device shell HT83G1C00069 -- pm list packages

Launch clock app

    gba-client device shell HT83G1C00069 -- cmd activity start-activity -n com.google.android.deskclock/com.android.deskclock.DeskClock`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
		output, err := client.ExecuteShellCommandOnDevice(args[0], strings.Join(args[1:], " "))
		if err != nil {
			panic(err)
		}

		fmt.Println(*output)
	},
}

var startSessionCmd = &cobra.Command{
	Use:   "start [DEVICE ID] [APP ID]",
	Short: "Start recording a session",
	Long:  "Start recording a session",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		screenshots, _ := cmd.Flags().GetBool("screenshots")
		autoSync, _ := cmd.Flags().GetBool("auto-sync")
		tags, _ := cmd.Flags().GetStringToString("tags")
		options := &gba.StartSessionOptions{Screenshots: screenshots, AutoSync: autoSync, Tags: tags}
		session, err := client.StartSession(args[0], args[1], options)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", session.Id)
	},
}

var stopSessionCmd = &cobra.Command{
	Use:   "stop [SESSION ID]",
	Short: "Stop a session",
	Long:  "stop",
	Args:  cobra.OnlyValidArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")
		if all {
			sessions, err := client.ListSessions()
			if err != nil {
				panic(err)
			}

			options := &gba.StopSessionOptions{}

			for _, session := range sessions {
				_, err = client.StopSession(session.Id, options)
				if err != nil {
					panic(err)
				}
			}

			return nil
		}

		if len(args) < 1 {
			return errors.New("requires a SESSION ID argument")
		}

		outputJson, err := cmd.Flags().GetBool("output-json")
		if err != nil {
			panic(err)
		}

		options := &gba.StopSessionOptions{}

		if outputJson {
			options.IncludeSessionJsonInResponse = true
		}

		outputDir, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			panic(err)
		}

		if outputDir != "" {
			options.OutputDir = outputDir
		}

		response, err := client.StopSession(args[0], options)
		if err != nil {
			panic(err)
		}

		if outputJson && response != nil {
			fmt.Println(*response)
		}

		return nil
	},
}

var syncSessionsCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all sessions",
	Long:  "Sync all sessions",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.Sync()
		if err != nil {
			panic(err)
		}
	},
}

var listSessionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List sessions",
	Long:  "List sessions",
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
	Use:   "list",
	Short: "List config properties",
	Long:  "List config properties",
	Run: func(cmd *cobra.Command, args []string) {
		properties, err := client.GetProperties()
		if err != nil {
			panic(err)
		}

		encodedProperties, err := json.MarshalIndent(properties, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(encodedProperties))
	},
}

var setPropertyCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config properties",
	Long:  `Set config properties

Retrieve properties

    gba-client property list > properties

Set properties

    cat properties | gba-client property set`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

		var input strings.Builder

		for scanner.Scan() {
			input.WriteString(scanner.Text())
		}

		if scanner.Err() != nil {
			panic(scanner.Err())
		}

		var properties map[string]interface{}
		err := json.Unmarshal([]byte(input.String()), &properties)
		if err != nil {
			panic(err)
		}

		err = client.SetProperties(properties)
		if err != nil {
			panic(err)
		}
	},
}

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "session",
	Long:  "session",
}

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "device",
	Long:  "device",
}

var propertyCmd = &cobra.Command{
	Use:   "property",
	Short: "property",
	Long:  "property",
}

type ClientVersionInfo struct {
	CommitHash string
	Version string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Run: func(cmd *cobra.Command, args []string) {
		clientVersionInfo := &ClientVersionInfo{CommitHash: commitHash, Version:version}
		fmt.Printf("Client: %+v\n", clientVersionInfo)

		serverVersionInfo, err := client.GetServerVersionInfo()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Server: %+v\n", serverVersionInfo)
	},
}

func main() {
	rootCmd.AddCommand(sessionCmd)
	sessionCmd.AddCommand(startSessionCmd)
	startSessionCmd.Flags().Bool("auto-sync", false, "Automatically sync session after it's stopped")
	startSessionCmd.Flags().Bool("screenshots", false, "Take screenshots during session")
	startSessionCmd.Flags().StringToString("tags", make(map[string]string), "")
	sessionCmd.AddCommand(stopSessionCmd)
	stopSessionCmd.Flags().Bool("all", false, "Stop all sessions")
	stopSessionCmd.Flags().Bool("output-json", false, "Output JSON")
	stopSessionCmd.Flags().String("output-dir", "", "Write JSON file to this directory")
	sessionCmd.AddCommand(syncSessionsCmd)
	sessionCmd.AddCommand(listSessionsCmd)

	rootCmd.AddCommand(deviceCmd)
	deviceCmd.AddCommand(listDevicesCmd)
	deviceCmd.AddCommand(listDeviceAppsCmd)
	deviceCmd.AddCommand(getDeviceCmd)
	deviceCmd.AddCommand(shellCmd)

	rootCmd.AddCommand(propertyCmd)
	propertyCmd.AddCommand(getPropertiesCmd)
	propertyCmd.AddCommand(setPropertyCmd)

	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "", "")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
