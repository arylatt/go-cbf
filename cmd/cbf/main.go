package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/arylatt/go-cbf/api"
	"github.com/spf13/cobra"
)

var cbfCmd = &cobra.Command{
	Use:          "cbf event category",
	Short:        "Get the JSON data for the provided category and Cambridge Beer Festival event",
	Example:      "cbf cbf2023 beer",
	Args:         cobra.ExactArgs(2),
	RunE:         cbfRunE,
	SilenceUsage: true,
}

func cbfRunE(cmd *cobra.Command, args []string) error {
	wc, err := api.NewWebClient()
	if err != nil {
		return fmt.Errorf("failed to create web client: %w", err)
	}

	resp, err := wc.Get(args[0], args[1])
	if err != nil {
		return fmt.Errorf("could not find data for event/category: %w", err)
	}

	if resp == nil {
		return fmt.Errorf("could not find data for event/category: response is empty")
	}

	data, _ := json.MarshalIndent(resp, "", "  ")

	fmt.Println(string(data))
	return nil
}

func main() {
	if cbfCmd.Execute() != nil {
		os.Exit(1)
	}
}
