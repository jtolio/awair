package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	dataFlags       = flag.NewFlagSet("", flag.ExitOnError)
	dataFlagsFormat = dataFlags.String("format", "human",
		"output format (human|json|csv|tsv)")

	cmdData = &ffcli.Command{
		Name:       "data",
		ShortHelp:  "get data",
		ShortUsage: fmt.Sprintf("%s [opts] data [opts] subcommand [opts]", os.Args[0]),
		Subcommands: []*ffcli.Command{
			cmdDataLatest,
		},
		FlagSet: dataFlags,
		Exec:    help,
	}
	cmdDataLatest = &ffcli.Command{
		Name:       "latest",
		ShortHelp:  "get the latest data",
		ShortUsage: fmt.Sprintf("%s [opts] data [opts] latest", os.Args[0]),
		Exec:       Latest,
	}
)

func Latest(ctx context.Context, args []string) error {
	if len(args) != 0 {
		return flag.ErrHelp
	}

	device, err := getDevice(ctx)
	if err != nil {
		return err
	}

	obs, err := device.Latest(ctx)
	if err != nil {
		return err
	}

	sort.Slice(obs.Sensors, func(i, j int) bool {
		return obs.Sensors[i].Component < obs.Sensors[j].Component
	})

	switch *dataFlagsFormat {
	case "human":
		fmt.Println(obs.Timestamp)
		fmt.Printf("score:\t%v\n", obs.Score)
		for _, r := range obs.Sensors {
			fmt.Printf("%s:\t%v\n", r.Component, r.Value)
		}
	case "json":
		data := map[string]interface{}{
			"timestamp": obs.Timestamp,
			"score":     obs.Score,
		}
		for _, r := range obs.Sensors {
			data[r.Component] = r.Value
		}
		serialized, err := json.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Println(string(serialized))
	case "csv", "tsv":
		sep := ","
		if *dataFlagsFormat == "tsv" {
			sep = "\t"
		}
		fmt.Print("timestamp" + sep + "score")
		for _, r := range obs.Sensors {
			fmt.Print(sep + r.Component)
		}
		fmt.Println()
		fmt.Printf("%v%s%v", obs.Timestamp, sep, obs.Score)
		for _, r := range obs.Sensors {
			fmt.Printf("%s%v", sep, r.Value)
		}
		fmt.Println()
	default:
		return fmt.Errorf("invalid format: %q", *dataFlagsFormat)
	}

	return nil
}
