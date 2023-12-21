package main

import (
	"fmt"
	"os"
)

func main() {
	rootCmd := NewRootCmd()
	versionCmd := NewVersionCmd()
	yahooAPICmd := NewYahooAPICmd()
	qqAPICmd := NewQQAPICmd()
	sinaAPICmd := NewSinaAPICmd()
	neteaseAPICmd := NewNeteaseAPICmd()
	eastmoneyAPICmd := NewEastmoneyAPICmd()
	eastmoneyLimitupAPICmd := NewEastmoneyLimitupAPICmd()
	eastmoneyLhbAPICmd := NewEastmoneyLhbAPICmd()
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(yahooAPICmd)
	rootCmd.AddCommand(qqAPICmd)
	rootCmd.AddCommand(sinaAPICmd)
	rootCmd.AddCommand(neteaseAPICmd)
	rootCmd.AddCommand(eastmoneyAPICmd)
	rootCmd.AddCommand(eastmoneyLimitupAPICmd)
	rootCmd.AddCommand(eastmoneyLhbAPICmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
