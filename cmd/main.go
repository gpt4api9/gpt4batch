/*
Copyright 2023 The gpt4batch Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"github.com/spf13/cobra"
	"gitlab.com/gpt4batch/cmd/authsvc"
	"gitlab.com/gpt4batch/cmd/batchsvc"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Seed the random number generator.
	// This is only done once, on program startup.
	rand.Seed(time.Now().Unix())

	// Create a context that is cancelled when the command is interrupted.
	// This context is used throughout the program.
	ctx := context.Background()
	rootCmd := NewCommand()
	rootCmd.AddCommand(authsvc.NewAuthenticationCommand(ctx))
	rootCmd.AddCommand(batchsvc.NewBatchCommand(ctx))
	rootCmd.SilenceUsage = true
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// NewCommand creates a new root command for gpt4batch.
func NewCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "batchsvc",
		Args:  cobra.NoArgs,
		Short: "Command line tool to gpt4api batch serve.<https://gpt4api.shop>",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.PrintErrf("See '%s -h' for help\n", cmd.CommandPath())
		},
	}
	return rootCmd
}
