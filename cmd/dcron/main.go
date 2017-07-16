package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"dcron/config"
)

var (
	// hostFlag is the flag to provide dcrond address.
	hostFlag = flag.String("H", "127.0.0.1:9090", "dcrond address")

	// listFlag is the flag provided to list cron config.
	listFlag = flag.Bool("l", false, "list cron config")

	// editFlag is the flag provided to edit cron config.
	editFlag = flag.Bool("e", false, "edit cron config")

	// baseAddr is the base address of dcrond.
	baseAddr string
)

func main() {
	flag.Parse()

	if !*listFlag && !*editFlag {
		flag.Usage()
		os.Exit(1)
	} else if *listFlag && *editFlag {
		flag.Usage()
		os.Exit(1)
	}

	baseAddr = fmt.Sprintf("http://%s", *hostFlag)
	if *listFlag {
		list()
	} else if *editFlag {
		// edit()
	}
}

func fetch() (*config.CronConfig, error) {
	listUrl := baseAddr + "/list"
	req, err := http.NewRequest("GET", listUrl, nil)
	if err != nil {
		return nil, errors.New("Failed to create list request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to connect to dcrond")
	}
	defer resp.Body.Close()

	var cronConf config.CronConfig
	if err := json.NewDecoder(resp.Body).Decode(&cronConf); err != nil {
		return nil, errors.New("Failed to decode dcrond response")
	}
	return &cronConf, nil
}

func list() {
	conf, err := fetch()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(conf.Config)
}
