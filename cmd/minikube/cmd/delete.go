/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package cmd

import (
	"os"

	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdUtil "k8s.io/minikube/cmd/util"
	"k8s.io/minikube/pkg/minikube/cluster"
	pkg_config "k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/console"
	"k8s.io/minikube/pkg/minikube/constants"
	"k8s.io/minikube/pkg/minikube/machine"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a local kubernetes cluster",
	Long: `Deletes a local kubernetes cluster. This command deletes the VM, and removes all
associated files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			console.ErrStyle("usage", "usage: minikube delete")
			os.Exit(1)
		}
		profile := viper.GetString(pkg_config.MachineProfile)
		console.OutStyle("deleting-vm", "Deleting %q Kubernetes cluster ...", profile)
		api, err := machine.NewAPIClient()
		if err != nil {
			console.Fatal("Error getting client: %v", err)
			os.Exit(1)
		}
		defer api.Close()

		if err = cluster.DeleteHost(api); err != nil {
			switch err := errors.Cause(err).(type) {
			case mcnerror.ErrHostDoesNotExist:
				console.OutStyle("meh", "%q VM does not exist", profile)
			default:
				console.Fatal("Failed to delete VM: %v", err)
				os.Exit(1)
			}
		} else {
			console.OutStyle("crushed", "VM deleted.")
		}

		if err := cmdUtil.KillMountProcess(); err != nil {
			console.Fatal("Failed to kill mount process: %v", err)
		}

		if err := os.Remove(constants.GetProfileFile(viper.GetString(pkg_config.MachineProfile))); err != nil {
			if os.IsNotExist(err) {
				console.OutStyle("meh", "%q profile does not exist", profile)
				os.Exit(0)
			}
			console.Fatal("Failed to remove profile: %v", err)
			os.Exit(1)
		}
		console.Success("Removed %q profile!", profile)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
