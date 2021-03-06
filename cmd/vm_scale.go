package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

// scaleCmd represents the scale command
var vmScaleCmd = &cobra.Command{
	Use:   "scale <vm name> [vm name] ...",
	Short: "Scale virtual machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		so, err := cmd.Flags().GetString("service-offering")
		if err != nil {
			return err
		}

		serviceoffering, err := getServiceOfferingByName(so)
		if err != nil {
			return err
		}

		tasks := make([]task, 0, len(args))
		for _, v := range args {
			vm, err := getVirtualMachineByNameOrID(v)
			if err != nil {
				return err
			}

			if vm.State != (string)(egoscale.VirtualMachineStopped) {
				fmt.Fprintf(os.Stderr, "this operation is not permitted if your VM is not stopped\n")
			}

			tasks = append(tasks, task{
				&egoscale.ScaleVirtualMachine{ID: vm.ID, ServiceOfferingID: serviceoffering.ID},
				fmt.Sprintf("Scaling %q ", vm.Name),
			})
		}

		taskResponses := asyncTasks(tasks)
		errors := filterErrors(taskResponses)
		if len(errors) > 0 {
			return errors[0]
		}

		return nil
	},
}

func init() {
	vmCmd.AddCommand(vmScaleCmd)
	vmScaleCmd.Flags().StringP("service-offering", "o", "", "<name | id> (micro|tiny|small|medium|large|extra-large|huge|mega|titan")
	if err := vmScaleCmd.MarkFlagRequired("service-offering"); err != nil {
		log.Fatal(err)
	}
}
