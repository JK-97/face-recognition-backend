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
	"github.com/spf13/cobra"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/model/device"
	"gitlab.jiangxingai.com/luyor/face-recognition-backend/internal/app/schema"
)

// addDeviceCmd represents the test command
var addDeviceCmd = &cobra.Command{
	Use:   "add_device",
	Short: "display_name device_id rtmpaddr",
	Long: `Add camera devices in database: display_name device_id rtmpaddr`,
	Args: cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
        model.InitDB()

        p := schema.Camera{
            Name: args[0],
            Rtmp: args[2],
            DeviceName: args[1],
        }
        device.AddCamera(&p)
	},
}

func init() {
	rootCmd.AddCommand(addDeviceCmd)
}
