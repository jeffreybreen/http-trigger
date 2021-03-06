/*
Copyright (c) 2016-2017 Bitnami

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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kubeless/http-trigger/pkg/controller"
	httptriggerutils "github.com/kubeless/http-trigger/pkg/utils"
	"github.com/kubeless/http-trigger/pkg/version"
	kubelessutils "github.com/kubeless/kubeless/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "http-trigger-controller",
	Short: "Kubeless HTTP trigger controller",
	Long:  "Kubeless HTTP trigger controller",
	Run: func(cmd *cobra.Command, args []string) {

		kubelessClient, err := kubelessutils.GetFunctionClientInCluster()
		if err != nil {
			logrus.Fatalf("Cannot get kubeless CR API client: %v", err)
		}

		httpTriggerClient, err := httptriggerutils.GetTriggerClientInCluster()
		if err != nil {
			logrus.Fatalf("Cannot get HTTP trigger CR API client: %v", err)
		}

		httpTriggerCfg := controller.HTTPTriggerConfig{
			KubeCli:        httptriggerutils.GetClient(),
			TriggerClient:  httpTriggerClient,
			KubelessClient: kubelessClient,
		}

		httpTriggerController := controller.NewHTTPTriggerController(httpTriggerCfg)

		stopCh := make(chan struct{})
		defer close(stopCh)

		go httpTriggerController.Run(stopCh)

		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM)
		signal.Notify(sigterm, syscall.SIGINT)
		<-sigterm
	},
}

func main() {
	logrus.Infof("Running Kubeless HTTP trigger controller version: %v", version.Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
