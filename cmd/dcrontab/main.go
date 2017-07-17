package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	"dcron/config"
	"dcron/cron"
	log "github.com/Sirupsen/logrus"
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

	// listUrl is the list API url.
	listUrl string

	// updateUrl is the update API url.
	updateUrl string
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
	listUrl = baseAddr + "/list"
	updateUrl = baseAddr + "/update"

	if *listFlag {
		list()
	} else if *editFlag {
		edit()
	}
}

func fetch() (*config.CronConfig, error) {
	req, err := http.NewRequest("GET", listUrl, nil)
	if err != nil {
		log.WithField("err", err).Error("Failed to create fetch request")
		return nil, errors.New("Failed to create list request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.WithField("err", err).Error("Failed to connect to dcrond")
		return nil, errors.New("Failed to connect to dcrond")
	}
	defer resp.Body.Close()

	var cronConf config.CronConfig
	if err := json.NewDecoder(resp.Body).Decode(&cronConf); err != nil {
		log.WithField("err", err).Error("Failed to decode dcrond response")
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
	fmt.Print(conf.Config)
}

func edit() {
	// Fetch config
	conf, err := fetch()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Write config to temp file
	log.Info("Loading config to temp file")
	rand.Seed(time.Now().UTC().UnixNano())
	tmpFilePath := fmt.Sprintf("/tmp/dcron-config-%d.tmp", rand.Int())
	ioutil.WriteFile(tmpFilePath, []byte(conf.Config), 0644)
	log.Info("Loaded config to temp file")

	// Open temp file in vi
	log.Info("Opening editor")
	cmd := exec.Command("vi", tmpFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.WithField("err", err).Error("Failed to open editor")
		fmt.Fprintln(os.Stderr, "Failed to open editor")
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		log.WithField("err", err).Error("Failed to wait for editor")
		fmt.Fprintln(os.Stderr, "Failed to wait for editor")
		os.Exit(1)
	}
	log.Info("Returned from editor")

	// Read updated content from temp file
	content, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		log.WithField("err", err).Error("Failed to read temp file")
		fmt.Fprintln(os.Stderr, "Failed to read temp file")
		os.Exit(1)
	}

	// Validate temp file
	contentStr := string(content)
	_, err = cron.MakeJobsFromString(contentStr)
	if err != nil {
		log.WithField("err", err).Error("Failed to parse config file")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Save new config
	conf.Config = contentStr
	jsonStr, err := json.Marshal(conf)
	if err != nil {
		log.WithField("err", err).Error("Failed to save config")
		fmt.Fprintln(os.Stderr, "Failed to save config")
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", updateUrl, bytes.NewBuffer(jsonStr))
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.WithField("err", err).Error("Failed to save config")
		fmt.Fprintln(os.Stderr, "Failed to save config")
		os.Exit(1)
	}

	// Validate response
	if resp.StatusCode != http.StatusOK {
		log.WithField("resp", resp).Error("Failed to save config")
		fmt.Fprintln(os.Stderr, "Failed to save config due to "+resp.Status+"!")
		os.Exit(1)
	}

	// Remove temp file
	if err = os.Remove(tmpFilePath); err != nil {
		log.WithField("err", err).Warning("Failed to remove temp file, continuing..")
	}
	fmt.Println("Config saved successfully!")
}
