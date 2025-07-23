/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"iidexic.dotstrike/dscore"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List user data",
	Long: `
	list all user data :)
`,
	Run: func(cmd *cobra.Command, args []string) {
		g, e := dscore.GetGlobals()
		if e != nil {
			panic(fmt.Errorf("globals failed:%v", e))
		}
		if len(args) == 0 {
			cmd.Print("USER DATA ---------\n")
			printsl := g.DescribeAllUserData()
			for _, p := range printsl {
				cmd.Print(p, "\n")
			}
		} else {
			for i, a := range args {
				cmd.Printf("%d: %v", i, g.CfgData(a))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
