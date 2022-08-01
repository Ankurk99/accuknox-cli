// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of KubeArmor

package cmd

import (
	"github.com/accuknox/accuknox-cli/summary"
	"github.com/spf13/cobra"
)

var summaryOptions summary.Options

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Policy summary from discovery engine",
	Long:  `Policy summary from discovery engine`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := summary.StartSummary(summaryOptions); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().StringVarP(&summaryOptions.Labels, "labels", "l", "", "Labels for resources")
	summaryCmd.Flags().StringVarP(&summaryOptions.Namespace, "namespace", "n", "", "Namespace for resources")
	summaryCmd.Flags().BoolVarP(&summaryOptions.RevDNSLookup, "rev-dns-lookup", "r", false, "Reverse DNS Lookup")
}
