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

	"github.com/gnames/gnfinder/grpc"
	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Provides gRPC API to gnfinder functionality. Default port 8778.",
	Long: `Starts a gRPC server using supplied port (8778 is default).

Provides gRPC API to gnfinder functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Printf("port: %d\n\n", port)
		grpc.Run(port)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	grpcCmd.Flags().IntP("port", "p", 8778, "grpc's port")
	log.SetFlags(0)
}
