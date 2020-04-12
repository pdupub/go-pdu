// Copyright 2019 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package console

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Console is the struct of cmd line console
type Console struct {
	targetURL string
	showCmd   *cobra.Command
	history   []string
}

// NewConsole used to build a new console
func NewConsole() (*Console, error) {
	cls := &Console{}
	cls.BldShowCmd()
	return cls, nil
}

// Close the console
func (c *Console) Close() {

}

// Run the console
func (c *Console) Run() {
	c.welcome()
	for {
		fmt.Print("> ")
		content := c.scanLine()
		if content == "quit" || content == "q" {
			c.Close()
			break
		} else if content == "show" {
			c.ExeShowCmd()
		} else {
			fmt.Println(content)
		}
	}
}

// initInfoDisplay
func (c Console) welcome() {
	fmt.Print(`
##############################################################
#                                                            #
#               Welcome to PDU console ;-)                   #
#                                                            #
##############################################################

`)
}

// SetTargetURL set the url of remote url
func (c *Console) SetTargetURL(url string) {
	c.targetURL = url
}

// BldShowCmd initial the show command
func (c *Console) BldShowCmd() {
	c.showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show command line",
		RunE: func(_ *cobra.Command, args []string) error {
			fmt.Println("Show information ......")
			return nil
		},
	}
}

// ExeShowCmd execute the show command
func (c *Console) ExeShowCmd() error {
	if err := c.showCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func (c Console) scanLine() string {
	var input string
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	input = string(data)
	return input
}
