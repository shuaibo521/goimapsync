package main

// Copyright (c) 2020 - Valentin Kuznetsov <vkuznet AT gmail dot com>
// configuration module for goimapsync
//

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Server structure keeps IMAP server's credentials
type Server struct {
	Name     string `json:"name"`     // name of IMAP server
	Uri      string `json:"uri"`      // IMAP URI
	Username string `json:"username"` // user name
	Password string `json:"password"` // user password
	UseTls   bool   `json:"useTls"`   // use TLS connection
}

// Configuration stores DAS configuration parameters
type Configuration struct {
	Servers     []Server `json:"servers"`     // list of IMAP server credentials
	Maildir     string   `json:"maildir"`     // maildir directory
	CommonInbox bool     `json:"commonInbox"` // use common inbox for all imap servers
	DBUri       string   `json:"dbUri"`       // DB URI
	Verbose     int      `json:"verbose"`     // verbosity level
}

// Config variable represents configuration object
var Config Configuration

// ParseConfig parse given config file
func ParseConfig(configFile string) error {
	var data []byte
	var err error
	if configFile == "-" {
		// read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		var content string
		for scanner.Scan() {
			content = fmt.Sprintf("%s%s", content, scanner.Text())
		}
		data = []byte(content)
		if err := scanner.Err(); err != nil {
			log.Println("unable to read from stdin", err)
			return err
		}
	} else {
		data, err = ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("Unable to read: file %s, error %v\n", configFile, err)
			return err
		}
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Printf("Unable to parse: file %s, error %v\n", configFile, err)
		return err
	}
	if Config.Maildir == "" {
		log.Fatal("Please specify maildir in your configuration")
	}
	if Config.DBUri == "" {
		Config.DBUri = fmt.Sprintf("sqlite3://%s/.goimapsync.db\n", Config.Maildir)
	}
	return nil
}
