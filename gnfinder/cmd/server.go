// Copyright © 2018 Dmitry Mozzherin <dmozzherin@gmail.com>
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
	"fmt"
	"log"
	"os"

	"github.com/gnames/gnfinder/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Provides RESTful API to gnfinder functionality. Default port 8778.",
	Long: `Starts an HTTP server using supplied port (8888 is default).

The server provides RESTful API to gnfinder functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Printf("port: %s\n\n", port)
		port = fmt.Sprintf(":%s", port)
		server.Run(port)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("port", "p", "8778", "server's port")
	log.SetFlags(0)
}
