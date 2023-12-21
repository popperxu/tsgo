package main

import (
	"fmt"
	"strings"

	//"github.com/saycv/tsgo"
	//"github.com/saycv/tsgo/pkg/configuration"

	"github.com/pkg/errors"
	logsupport "github.com/saycv/tsgo/pkg/log"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command
func NewRootCmd() *cobra.Command {

	var logLevel string

	rootCmd := &cobra.Command{
		Use:   "tsgo [-logLevel]",
		Short: `tsgo is a command-line utility that displays stocks`,
		Args:  cobra.OnlyValidArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			lvl, err := log.ParseLevel(logLevel)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "unable to parse log level '%v'", logLevel)
				return err
			}
			logsupport.Setup()
			log.SetLevel(lvl)
			log.SetOutput(cmd.OutOrStdout())
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd, args, e := cmd.Root().Find(args)
			if cmd == nil || e != nil || len(args) > 0 {
				return errors.Errorf("unknown help topic: %v", strings.Join(args, " "))
			}

			return nil
		},
	}
	flags := rootCmd.Flags()

	flags.StringVarP(&logLevel, "log", "l", "debug", "log level to set [debug|info|warning|error|fatal|panic]")

	return rootCmd
}
