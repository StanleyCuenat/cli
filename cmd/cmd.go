package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// cmdCheckRequiredFlags evaluates the specified flags as parsed in the cobra.Command flagset to check that
// their value is unset (i.e. null/empty/zero, depending on the type), and returns a multierror listing all
// flags missing a required value.
func cmdCheckRequiredFlags(cmd *cobra.Command, flags []string) error {
	var err *multierror.Error

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		for _, fn := range flags {
			if flag.Name == fn {
				var hasValue bool

				switch flag.Value.Type() {
				case "string":
					if flag.Value.String() != "" {
						hasValue = true
					}

				case "int", "int8", "int16", "int32", "int64", "float32", "float64":
					if flag.Value.String() != "0" {
						hasValue = true
					}
				}

				if !hasValue {
					err = multierror.Append(err, fmt.Errorf("no value specified for flag %q", fn))
				}
			}
		}
	})

	return err.ErrorOrNil()
}

func cmdExitOnUsageError(cmd *cobra.Command, reason string) {
	cmd.PrintErrln(fmt.Sprintf("error: %s", reason))
	cmd.Usage() // nolint:errcheck
	os.Exit(1)
}
