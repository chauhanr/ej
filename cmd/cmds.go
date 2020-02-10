package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func version(cmd *cobra.Command, args []string) {
	fmt.Println(banner)
	fmt.Println("EJ - Explore Jira version 1.0")
}
