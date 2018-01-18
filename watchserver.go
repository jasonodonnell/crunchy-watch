/*
 Copyright 2017 Crunchy Data Solutions, Inc.
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
	"github.com/crunchydata/crunchy-watch/plugins"
	api "github.com/crunchydata/crunchy-watch/watchapi"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var POLL_INT = int64(3)

func main() {
	api.GetEnv()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		api.Logger.Println(sig)
		api.Logger.Println("caught signal, cleaning up and exiting...")
		os.Exit(0)
	}()

	var VERSION = os.Getenv("CCP_VERSION")

	api.Logger.Println("watchserver " + VERSION + ": starting")

	for true {
		if api.HealthCheck() == false {
			switch api.EnvVars.PROJECT_TYPE {
			case "docker":
				plugins.DockerFailover()
				api.Logger.Println("docker failover exits normally")
				os.Exit(0)
			case "kube":
				plugins.KubeFailover()
			case "openshift":
				plugins.OpenshiftFailover()
			case "cf":
				plugins.CloudFoundryFailover()
			default:
				api.Logger.Println(api.EnvVars.PROJECT_TYPE + " handling not implemented")
			}
		} else {
			api.Logger.Println("master is responding...")
		}
		api.Logger.Printf("sleeping for %d\n", api.EnvVars.SLEEP_TIME)
		time.Sleep(time.Duration(api.EnvVars.SLEEP_TIME) * time.Second)
	}

}