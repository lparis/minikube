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
	"time"

	"github.com/docker/machine/libmachine/mcnerror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cmdUtil "k8s.io/minikube/cmd/util"
	"k8s.io/minikube/pkg/minikube/cluster"
	pkg_config "k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/console"
	"k8s.io/minikube/pkg/minikube/machine"
	pkgutil "k8s.io/minikube/pkg/util"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops a running local kubernetes cluster",
	Long: `Stops a local kubernetes cluster running in Virtualbox. This command stops the VM
itself, leaving all files intact. The cluster can be started again with the "start" command.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := viper.GetString(pkg_config.MachineProfile)
		console.OutStyle("stopping", "Stopping %q Kubernetes cluster...", profile)

		api, err := machine.NewAPIClient()
		if err != nil {
			console.Fatal("Error getting client: %v", err)
			os.Exit(1)
		}
		defer api.Close()

		nonexistent := false

		stop := func() (err error) {
			err = cluster.StopHost(api)
			switch err := errors.Cause(err).(type) {
			case mcnerror.ErrHostDoesNotExist:
				console.OutStyle("meh", "%q VM does not exist, nothing to stop", profile)
				nonexistent = true
				return nil
			default:
				return err
			}
		}
		if err := pkgutil.RetryAfter(5, stop, 1*time.Second); err != nil {
			console.Fatal("Unable to stop VM: %v", err)
			cmdUtil.MaybeReportErrorAndExit(err)
		}
		if !nonexistent {
			console.OutStyle("stopped", "%q stopped.", profile)
		}

		if err := cmdUtil.KillMountProcess(); err != nil {
			console.Fatal("Unable to kill mount process: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
