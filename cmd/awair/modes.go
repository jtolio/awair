package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jtolio/awair"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type GetterFn func(context.Context) (string, error)
type SetterFn func(context.Context, string) error

type modeControl struct {
	Name           string
	Methods        func(d *awair.Device) (GetterFn, SetterFn)
	PossibleStates []string
}

var (
	modes = []modeControl{
		{
			Name: "display",
			Methods: func(d *awair.Device) (GetterFn, SetterFn) {
				return d.GetDisplayMode, d.SetDisplayMode
			},
			PossibleStates: []string{
				"default", "score", "clock",
				"temp_humid_celsius", "temp_humid_fahrenheit",
			},
		},
		{
			Name: "knocking",
			Methods: func(d *awair.Device) (GetterFn, SetterFn) {
				return d.GetKnockingMode, d.SetKnockingMode
			},
			PossibleStates: []string{"on", "off", "sleep"},
		},
		{
			Name: "led",
			Methods: func(d *awair.Device) (GetterFn, SetterFn) {
				return d.GetLEDMode, d.SetLEDMode
			},
			PossibleStates: []string{"on", "dim", "sleep"},
		},
	}
)

func init() {
	for _, m := range modes {
		m := m
		getter := &ffcli.Command{
			Name:       "get",
			ShortHelp:  fmt.Sprintf("get the %s mode", m.Name),
			ShortUsage: fmt.Sprintf("%s [opts] %s get", os.Args[0], m.Name),
			Exec: func(ctx context.Context, args []string) error {
				if len(args) != 0 {
					return flag.ErrHelp
				}
				device, err := getDevice(ctx)
				if err != nil {
					return err
				}
				getter, _ := m.Methods(device)
				mode, err := getter(ctx)
				if err != nil {
					return err
				}
				_, err = fmt.Println(mode)
				return err
			},
		}
		options := strings.Join(m.PossibleStates, ", ")
		setter := &ffcli.Command{
			Name:       "set",
			ShortHelp:  fmt.Sprintf("set the %s mode (%s)", m.Name, options),
			ShortUsage: fmt.Sprintf("%s [opts] %s set <mode>", os.Args[0], m.Name),
			Exec: func(ctx context.Context, args []string) error {
				if len(args) != 1 {
					fmt.Println("mode should be one of " + options)
					return flag.ErrHelp
				}

				device, err := getDevice(ctx)
				if err != nil {
					return err
				}
				_, setter := m.Methods(device)
				return setter(ctx, args[0])
			},
		}
		cmdRoot.Subcommands = append(cmdRoot.Subcommands, &ffcli.Command{
			Name:        m.Name,
			ShortHelp:   fmt.Sprintf("control the %s mode", m.Name),
			ShortUsage:  fmt.Sprintf("%s [opts] %s subcommand [opts]", os.Args[0], m.Name),
			Subcommands: []*ffcli.Command{getter, setter},
			Exec:        help,
		})
	}
	sort.Slice(cmdRoot.Subcommands, func(i, j int) bool { return cmdRoot.Subcommands[i].Name < cmdRoot.Subcommands[j].Name })
}
