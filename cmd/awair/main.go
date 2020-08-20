package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jtolio/awair"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	sysFlags      = flag.NewFlagSet("", flag.ExitOnError)
	sysFlagConfig = sysFlags.String("config", filepath.Join(homeDir(), ".awair"),
		"path to config file")
	sysFlagToken       = sysFlags.String("token", "", "awair bearer token")
	sysFlagsFahrenheit = sysFlags.Bool("fahrenheit", false, "prefer fahrenheit")

	cmdRoot = &ffcli.Command{
		ShortHelp:  "control your awair",
		ShortUsage: fmt.Sprintf("%s [opts] subcommand [opts]", os.Args[0]),
		Subcommands: []*ffcli.Command{
			cmdData,
		},
		FlagSet: sysFlags,
		Options: []ff.Option{
			ff.WithAllowMissingConfigFile(true),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithConfigFileFlag("config"),
		},
		Exec: help,
	}
)

func homeDir() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	if u.HomeDir == "" {
		panic("no homedir found")
	}
	return u.HomeDir
}

func main() {
	err := cmdRoot.ParseAndRun(context.Background(), os.Args[1:])
	if err == nil {
		return
	}
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %+v\n", err)
	os.Exit(1)
}

func help(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}

func getDevice(ctx context.Context) (*awair.Device, error) {
	client := awair.NewClientFromBearerToken(*sysFlagToken)
	client.Options.PreferFahrenheit = *sysFlagsFahrenheit
	devices, err := client.GetDevices(ctx)
	if err != nil {
		return nil, err
	}
	if len(devices) != 1 {
		return nil, fmt.Errorf("not exactly one device found")
	}
	return devices[0], nil
}
