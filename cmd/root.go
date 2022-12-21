/*
Copyright © 2022 Kaleb Hawkins <KalebHawkins@outlook.com>

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
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/KalebHawkins/ggpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var apiKey string

var (
	prompt         string
	size           string
	outfile        string
	displayVersion bool
	version        string
	commit         string
)

var (
	errNoApiKey = errors.New("there was no API_KEY specified. You can do so by exporting `AIIMG_API_KEY` or adding `AIIMG_API_KEY` to your ~/.aiimg.yaml configuration file")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aiimg",
	Short: "Generate images with text and AI :)",
	Long: `AIImg can be used to generate images from text using OpenAI's 
images API endpoint.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if displayVersion {
			fmt.Printf("Version: %s\nCommit: %s", version, commit)
			os.Exit(0)
		}

		apiKey = viper.GetString("AIIMG_API_KEY")
		if apiKey == "" {
			return errNoApiKey
		}

		c := ggpt3.NewClient(apiKey)

		imgReq := ggpt3.ImageRequest{
			Prompt:         prompt,
			N:              1,
			Size:           size,
			ResponseFormat: "b64_json",
		}

		imgResp, err := c.RequestImages(context.Background(), &imgReq)
		if err != nil {
			return err
		}

		encodedImg := imgResp.Data[0].B64
		imgBuf, err := base64.StdEncoding.DecodeString(encodedImg)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		err = os.WriteFile(outfile, imgBuf, 0644)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aiimg.yaml)")
	rootCmd.Flags().StringVarP(&prompt, "prompt", "p", "A chicken with it's head cut off", "describe what to generate")
	rootCmd.Flags().StringVarP(&size, "size", "s", "512x512", "size of the image to generate")
	rootCmd.Flags().StringVarP(&outfile, "outfile", "o", "img.png", "the file to output")
	rootCmd.Flags().BoolVarP(&displayVersion, "version", "v", false, "version information")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".aiimg")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
