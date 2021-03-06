// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"net"
	"os"
	"os/signal"
	"syscall"

	sf "github.com/fpay/snowflake-go"
	pb "github.com/fpay/snowflake-go/pb"

	"github.com/fpay/foundation-go/cache"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "snowflake",
	Short: "A simple Twitter snowflake generator.",
	Long:  `Snowflake is a network service for generating unique ID numbers at high scale with some simple guarantees.`,
	Run: func(cmd *cobra.Command, args []string) {
		start(cmd, args)
	},
}

// 启动snowflake服务
func start(cmd *cobra.Command, args []string) {
	opts := new(ConfigOptions)
	opts.Load()

	redis, err := cache.NewRedisCache(opts.Redis)
	handleInitError("redis", err)

	sfs, err := sf.NewSnowflakeServer(redis)
	handleInitError("snowflake server", err)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.Server.Port))
	handleInitError("grpc server", err)

	gs := grpc.NewServer()
	pb.RegisterSnowflakeServiceServer(gs, sfs)
	go gs.Serve(listen)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	gs.GracefulStop()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snowflake.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".my" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".snowflake")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type ServerOptions struct {
	Port int `mapstructure:"port"`
}

type ConfigOptions struct {
	Server ServerOptions      `mapstructure:"server"`
	Redis  cache.RedisOptions `mapstructure:"redis"`
}

func (co *ConfigOptions) Load() {
	err := viper.Unmarshal(co)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}
}

func handleInitError(module string, err error) {
	if err == nil {
		return
	}
	log.Fatalf("init %s failed, error: %s", module, err)
}
