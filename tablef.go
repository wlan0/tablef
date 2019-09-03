// Copyright 2019 Sidhartha Mani <sidharthamn@gmail.com>
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

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "tablef",
	Short: "create tabular output",
	Long: `                                                                                                                                     
echo "data1 data2 123.456 789" | tablef "%10s\t%10s\t%.2f\t%x"                                                                                                                              
`,
	Run: func(c *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Printf("format string expected\n")
			os.Exit(1)
		}
		if err := tablef(os.Stdin, os.Stdout, args[0]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
}

func main() {
	cmd.Execute()
}

var blue = color.New(color.FgCyan)
var green = color.New(color.FgYellow, color.Faint)

func tablef(input io.Reader, output *os.File, format string) error {
	c := blue
	buf := bufio.NewReader(input)

	for {
		newBuf := &bytes.Buffer{}
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("tablef: %v", err)
		}
		fields := strings.Fields(line)
		if err := printf(newBuf, append([]string{format}, fields...)...); err != nil {
			return fmt.Errorf("tablef: %v", err)
		}

		if c == blue {
			c = green
		} else {
			c = blue
		}
		c.Fprintln(output, newBuf.String())
	}
	return nil
}
