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
	"github.com/spf13/viper"
	"github.com/zzOzz/izlyctl/crous"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"fmt"
)

// cardsCmd represents the cards command
var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Debugf("user: %s", viper.GetString("api.auth.login"))
		logrus.Debugf("auth: %s", crous.GetAuth())
		client := crous.New()
		client.Authenticate()
		card:=client.GetCard(args[0])
		b, err := json.Marshal(card)
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println(string(b))
		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(cardsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cardsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cardsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
