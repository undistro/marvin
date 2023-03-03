// Copyright 2023 Undistro Authors
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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/version"
)

var (
	versionOutput string
	// versionCmd represents the version command
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of Marvin",
		RunE: func(c *cobra.Command, args []string) error {
			v := version.Get()
			var s string
			switch versionOutput {
			case "json":
				b, err := json.MarshalIndent(&v, "", "    ")
				if err != nil {
					return err
				}
				s = string(b)
			case "yaml":
				b, err := yaml.Marshal(&v)
				if err != nil {
					return err
				}
				s = string(b)
			default:
				s = v.String()
			}
			fmt.Println(s)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionOutput, "output", "o", "", `Output format. One of: ("json", "yaml")`)
}
