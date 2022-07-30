/*
Copyright © 2022 Wissem BEN CHAABANE<benchaaben.wissem@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"

	"tagu/models"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Dump tags from AWS cloud provider",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: awsCmdRunE,
}

func awsCmdRunE(c *cobra.Command, args []string) (err error) {
	filePath, err := c.Flags().GetString("input-file")
	if err != nil {
		return err
	}
	_, err = awsloadConfig(filePath)
	if err != nil {
		return err
	}
	c.Printf("Load configuration file %s", viper.ConfigFileUsed())
	return nil
}

func initAwsFlags(c *cobra.Command) {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	c.Flags().StringP("input-file", "i", "", "the input file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func init() {
	rootCmd.AddCommand(awsCmd)
	initAwsFlags(awsCmd)
}

// awsloadConfig reads in config file and ENV variables if set.
func awsloadConfig(tags string) (spec models.Spec, err error) {
	viper.SetConfigType("yaml")
	if tags != "" {
		viper.SetConfigFile(tags)
	} else {
		viper.AutomaticEnv() // read in environment variables that match
		viper.SetConfigFile(viper.GetString("aws_config"))
	}

	err = viper.ReadInConfig() // Read config file
	if err != nil {
		return spec, fmt.Errorf("%s", err)
	}
	if viper.IsSet("accounts") {
		var gspec models.GeneralSpec
		err = viper.Unmarshal(&gspec) // Load config file in struct object
		if err != nil {
			return spec, err
		}
		spec.UniformConfig(gspec) // Build generic config into detailed one
	} else if viper.IsSet("filter-input") {
		err = viper.Unmarshal(&spec) // Load config file in struct object
		if err != nil {
			return spec, err
		}
	} else {
		return spec, fmt.Errorf("invalid configuration format for file %s", viper.ConfigFileUsed())
	}

	return spec, nil
}
