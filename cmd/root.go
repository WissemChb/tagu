/*
Copyright © 2022 Wissem BEN CHAABANE <benchaaben.wissem@gmail.com>

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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tagu",
	Short: "Dump the resources tags from public cloud provider",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:tvv

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: rootCmdRunE,
}

func rootCmdRunE(c *cobra.Command, args []string) (err error) {
	configPath, err := c.Flags().GetString("config")
	if err != nil {
		return err
	}
	err = rootLoadConfig(configPath)
	if err != nil {
		return err
	}
	c.Printf("Load configuration file %s", viper.ConfigFileUsed())
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initRootFlags(cmd *cobra.Command) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cmd.Flags().StringP("config", "c", "", "config file (default is $HOME/.tagu.yaml)")
	// cmd.Flags().StringP("input-file", "i", "", "the input file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func init() {
	initRootFlags(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
func rootLoadConfig(flag string) (err error) {
	if flag != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flag)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tagu" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("input-tags")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err = viper.ReadInConfig() // Read config file
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}
