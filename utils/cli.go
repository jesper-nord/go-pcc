package main

import (
	"flag"
	"fmt"
	"github.com/hacktobeer/go-panasonic/cloudcontrol"
	pt "github.com/hacktobeer/go-panasonic/types"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
)

var (
	commit  = "development"
	date    = "development"
	version = "development"

	configFlag  = flag.String("config", "gopanasonic.yaml", "Path of YAML configuration file")
	debugFlag   = flag.Bool("debug", false, "Show debug output")
	deviceFlag  = flag.String("device", "", "Device to issue command to")
	historyFlag = flag.String("history", "", "Display history: day,week,month,year")
	listFlag    = flag.Bool("list", false, "List available devices")
	modeFlag    = flag.String("mode", "", "Set mode: auto,heat,cool,dry,fan")
	offFlag     = flag.Bool("off", false, "Turn device off")
	onFlag      = flag.Bool("on", false, "Turn device on")
	quietFlag   = flag.Bool("quiet", false, "Don't output any log messages")
	statusFlag  = flag.Bool("status", false, "Display current status of device")
	tempFlag    = flag.Float64("temp", 0, "Set the temperature (in Celsius)")
	versionFlag = flag.Bool("version", false, "Show build version information")
)

func readConfig() {
	viper.SetConfigFile(*configFlag)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("version: %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("date: %s\n", date)
		os.Exit(0)
	}

	log.SetLevel(log.INFO)

	if *quietFlag {
		log.SetLevel(log.ERROR)
	}

	if *debugFlag {
		log.SetLevel(log.DEBUG)
		log.Debug("log set to debug level")
	}

	readConfig()
	user := viper.GetString("username")
	pass := viper.GetString("password")
	token := viper.GetString("token")

	client := cloudcontrol.NewClient()

	if token == "" {
		createAndSaveSession(user, pass, &client)
	} else {
		if body, err := client.ValidateSession(token); err != nil {
			log.Debugf("invalid session token: %s", string(body))
			createAndSaveSession(user, pass, &client)
		} else {
			log.Debug("session token is valid")
		}
	}

	if *listFlag {
		log.Debug("listing available devices")
		devices, err := client.ListDevices()
		if err != nil {
			log.Fatal(err)
		}

		if len(devices) == 0 {
			log.Fatal("found no devices for configured account")
		}

		log.Infof("%d device(s) found:\n", len(devices))
		for _, device := range devices {
			fmt.Println(device)
		}
		os.Exit(0)
	}

	// read device from configuration file
	configDevice := viper.GetString("device")
	if configDevice != "" {
		log.Debugf("using device %s from config file", configDevice)
		client.SetDevice(configDevice)
	}

	// read device from flag (higher priority)
	if *deviceFlag != "" {
		log.Debugf("using device %s from flag", *deviceFlag)
		client.SetDevice(*deviceFlag)
	}

	if client.DeviceGUID == "" {
		log.Fatal("no device configured, use -device flag or set device in config file")
	}

	if *statusFlag {
		log.Debug("fetching device status")
		status, err := client.GetDeviceStatus()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Device GUID: %s\n", status.DeviceGUID)
		fmt.Println("Capabilities:")
		fmt.Printf("Auto mode: %t\n", status.AutoMode)
		fmt.Printf("Heat mode: %t\n", status.HeatMode)
		fmt.Printf("Dry mode: %t\n", status.DryMode)
		fmt.Printf("Cool mode: %t\n", status.CoolMode)
		fmt.Printf("Fan mode: %t\n", status.FanMode)
		fmt.Printf("Fan Speed mode: %d\n", status.FanSpeedMode)
		fmt.Printf("Quiet mode: %t\n", status.QuietMode)
		fmt.Printf("Eco function: %d\n", status.EcoFunction)
		fmt.Printf("EcoNavi function: %t\n", status.EcoNavi)
		fmt.Printf("iAutoX: %t\n", status.IautoX)
		fmt.Printf("NanoeX: %t\n", status.Nanoe)
		fmt.Println("Current status:")
		fmt.Printf("Status: %s\n", pt.Operate[status.Parameters.Operate])
		fmt.Printf("Online: %t\n", status.Parameters.Online)
		fmt.Printf("Temperature: %0.1f\n", status.Parameters.TemperatureSet)
		fmt.Printf("Mode: %s\n", pt.ModesReverse[status.Parameters.OperationMode])
	}

	if *historyFlag != "" {
		log.Debugf("fetching historical data for %s\n", *historyFlag)
		history, err := client.GetDeviceHistory(pt.HistoryDataMode[*historyFlag])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("#,AverageSettingTemp,AverageInsideTemp,AverageOutsideTemp")
		for _, v := range history.HistoryEntries {
			fmt.Printf("%v,%v,%v,%v\n", v.DataNumber+1, v.AverageSettingTemp, v.AverageInsideTemp, v.AverageOutsideTemp)
		}
	}

	if *onFlag {
		log.Debug("turning device on")
		_, err := client.TurnOn()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("device turned on")
	}

	if *offFlag {
		log.Debug("turning device off")
		_, err := client.TurnOff()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("device turned off")
	}

	if *tempFlag != 0 {
		log.Debugf("setting temperature to %v degrees", *tempFlag)
		_, err := client.SetTemperature(*tempFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("temperature set to %v degrees", *tempFlag)
	}

	if *modeFlag != "" {
		log.Debugf("setting mode to %s", *modeFlag)
		_, err := client.SetMode(pt.Modes[*modeFlag])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("mode set to %s", *modeFlag)
	}
}

func createAndSaveSession(user string, pass string, client *cloudcontrol.Client) {
	if user == "" || pass == "" {
		log.Fatal("missing username and/or password in config file")
	}

	_, err := client.CreateSession(user, pass)
	if err != nil {
		log.Fatal(err)
	}

	viper.Set("token", client.Utoken)
	err = viper.WriteConfig()
	if err != nil {
		log.Fatal("unable to write session token to config file")
	}

	log.Debugf("new session token created and written to config file: %s", client.Utoken)
}
