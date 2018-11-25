package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	ssgl "github.com/ranForkce/shadowsocks-go/cmd/shadowsocks-local/local"
	ss "github.com/ranForkce/shadowsocks-go/shadowsocks"
)

var debug ss.DebugLog

func main() {
	log.SetOutput(os.Stdout)

	var configFile, cmdServer, cmdURI string
	var cmdConfig ss.Config
	var printVer bool

	flag.BoolVar(&printVer, "version", false, "print version")
	flag.StringVar(&configFile, "c", "config.json", "specify config file")
	flag.StringVar(&cmdServer, "s", "", "server address")
	flag.StringVar(&cmdConfig.LocalAddress, "b", "", "local address, listen only to this address if specified")
	flag.StringVar(&cmdConfig.Password, "k", "", "password")
	flag.IntVar(&cmdConfig.ServerPort, "p", 0, "server port")
	flag.IntVar(&cmdConfig.Timeout, "t", 300, "timeout in seconds")
	flag.IntVar(&cmdConfig.LocalPort, "l", 0, "local socks5 proxy port")
	flag.StringVar(&cmdConfig.Method, "m", "", "encryption method, default: aes-256-cfb")
	flag.BoolVar((*bool)(&debug), "d", false, "print debug message")
	flag.StringVar(&cmdURI, "u", "", "shadowsocks URI")

	flag.Parse()

	if s, e := ssgl.ParseURI(cmdURI, &cmdConfig); e != nil {
		log.Printf("invalid URI: %s\n", e.Error())
		flag.Usage()
		os.Exit(1)
	} else if s != "" {
		cmdServer = s
	}

	if printVer {
		ss.PrintVersion()
		os.Exit(0)
	}

	cmdConfig.Server = cmdServer
	ss.SetDebug(debug)

	exists, err := ss.IsFileExists(configFile)
	// If no config file in current directory, try search it in the binary directory
	// Note there's no portable way to detect the binary directory.
	binDir := path.Dir(os.Args[0])
	if (!exists || err != nil) && binDir != "" && binDir != "." {
		oldConfig := configFile
		configFile = path.Join(binDir, "config.json")
		log.Printf("%s not found, try config file %s\n", oldConfig, configFile)
	}

	config, err := ss.ParseConfig(configFile)
	if err != nil {
		config = &cmdConfig
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", configFile, err)
			os.Exit(1)
		}
	} else {
		ss.UpdateConfig(config, &cmdConfig)
	}
	ssgl.RunConf(config)
}
