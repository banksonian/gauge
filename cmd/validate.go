// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/getgauge/gauge/config"
	"github.com/getgauge/gauge/logger"
	"github.com/getgauge/gauge/util"
	"github.com/getgauge/gauge/validation"
	"github.com/spf13/cobra"
)

const (
	hideSuggestionDefault = false
	hideSuggestionName    = "hide-suggestion"
)

var (
	validateCmd = &cobra.Command{
		Use:     "validate [flags] [args]",
		Short:   "Check for validation and parse errors",
		Long:    `Check for validation and parse errors.`,
		Example: "  gauge validate specs/",
		Run: func(cmd *cobra.Command, args []string) {
			loadEnvAndReinitLogger(cmd)
			validation.HideSuggestion = hideSuggestion
			if err := config.SetProjectRoot(args); err != nil {
				exit(err, cmd.UsageString())
			}
			if len(args) == 0 {
				args = append(args, util.GetSpecDirs()...)
			}
			res := validation.ValidateSpecs(args, false)
			if currentContext != nil && currentContext.Value(CommandContext("Command")) == "execution" {
				currentContext = context.WithValue(currentContext, CommandContext("ValidationResult"), res)
			} else {
				if len(res.Errs) > 0 {
					os.Exit(1)
				}
				if res.SpecCollection.Size() < 1 {
					logger.Infof(true, "No specifications found in %s.", strings.Join(args, ", "))
					res.Runner.Kill()
					if res.ParseOk {
						os.Exit(0)
					}
					os.Exit(1)
				}
				res.Runner.Kill()
				if res.ErrMap.HasErrors() {
					os.Exit(1)
				}
				logger.Infof(true, "No errors found.")
			}
		},
		DisableAutoGenTag: true,
	}
	hideSuggestion bool
)

func init() {
	GaugeCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVarP(&hideSuggestion, "hide-suggestion", "", false, "Prints a step implementation stub for every unimplemented step")

}
