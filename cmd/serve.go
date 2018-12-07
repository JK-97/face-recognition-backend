// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net/http"

	"github.com/spf13/cobra"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/config"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/remote"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/route"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/timer"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve http backend for face-recognition",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		model.InitDB()

		go func() {
			port, _ := cmd.Flags().GetString("port")
			log.Fatal(http.ListenAndServe(":" + port, route.Routes()))
		}()

		go func() {
			timer.Init()
		}()

        go func() {
            remote.AddDevices()
        }()

		forever := make(chan struct{})
		<-forever
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().String("port", "80", "Port to run Application server on")

	cfg := config.Config()
	serveCmd.PersistentFlags().String("db-addr", "mongodb://192.168.3.33", "Mongo db server address")
	cfg.BindPFlag("db-addr", serveCmd.PersistentFlags().Lookup("db-addr"))

	serveCmd.PersistentFlags().String("camera-addr", "http://192.168.0.196:8088", "Camera server address")
	cfg.BindPFlag("camera-addr", serveCmd.PersistentFlags().Lookup("camera-addr"))

	serveCmd.PersistentFlags().String("face-ai-addr", "http://192.168.0.196:8008", "Face ai server address")
	cfg.BindPFlag("face-ai-addr", serveCmd.PersistentFlags().Lookup("face-ai-addr"))

    serveCmd.PersistentFlags().String("login_type", "local", "Login type: local or gateway")
	cfg.BindPFlag("login_type", serveCmd.PersistentFlags().Lookup("login_type"))

    serveCmd.PersistentFlags().String("apigateway-addr", "http://192.168.0.196:8008", "APIGateway server address")
	cfg.BindPFlag("apigateway-addr", serveCmd.PersistentFlags().Lookup("apigateway-addr"))
}
