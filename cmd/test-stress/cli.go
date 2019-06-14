package main

import (
	"github.com/gsmcwhirter/go-util/v3/cli"
	"github.com/gsmcwhirter/go-util/v3/errors"
	"github.com/spf13/viper"
)

func setup(start func(config) error) *cli.Command {
	c := cli.NewCLI(AppName, BuildVersion, BuildSHA, BuildDate, cli.CommandOptions{
		ShortHelp:    "discord bot lib stress test",
		Args:         cli.NoArgs,
		SilenceUsage: true,
	})

	var configFile string
	c.Flags().StringVar(&configFile, "config", "", "The config file to use")
	c.Flags().String("pprof_hostport", "", "The host and port for the pprof http server to listen on")

	c.SetRunFunc(func(cmd *cli.Command, args []string) (err error) {
		v := viper.New()

		if configFile != "" {
			v.SetConfigFile(configFile)
		} else {
			v.SetConfigName("stress-config")
			v.AddConfigPath(".") // working directory
		}

		v.SetEnvPrefix("EDB")
		v.AutomaticEnv()

		err = v.BindPFlags(cmd.Flags())
		if err != nil {
			return errors.Wrap(err, "could not bind flags to viper")
		}

		err = v.ReadInConfig()
		if err != nil {
			return errors.Wrap(err, "could not read in config file")
		}

		conf := config{}
		err = v.Unmarshal(&conf)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal config into struct")
		}

		return start(conf)
	})

	return c
}
