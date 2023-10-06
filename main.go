package main

import (
	"flag"
	"fmt"
	cloudcontrol "github.com/jesper-nord/go-pcc/client"
	"github.com/jesper-nord/go-pcc/types"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
)

var (
	configFlag   = flag.String("config", "./go-pcc.yaml", "Path of YAML configuration file")
	debugFlag    = flag.Bool("debug", false, "Show debug output")
	deviceFlag   = flag.String("device", "", "Device to issue command to")
	historyFlag  = flag.String("history", "", "Display history: day,week,month,year")
	listFlag     = flag.Bool("list", false, "List available devices")
	modeFlag     = flag.String("mode", "", "Set mode: auto,heat,cool,dry,fan")
	ecoModeFlag  = flag.String("ecomode", "", "Set eco mode: auto,powerful,quiet")
	offFlag      = flag.Bool("off", false, "Turn device off")
	onFlag       = flag.Bool("on", false, "Turn device on")
	suppressFlag = flag.Bool("suppress", false, "Suppress log messages")
	statusFlag   = flag.Bool("status", false, "Display current status of device")
	tempFlag     = flag.Float64("temp", 0, "Set the temperature (in Celsius)")
	fanSpeedFlag = flag.String("speed", "", "Set fan speed: auto,1,2,3,4,5")
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
	log.SetHeader(`${time_rfc3339} ${level} -`)

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	log.SetLevel(log.INFO)
	if *suppressFlag {
		log.SetLevel(log.OFF)
	}

	if *debugFlag {
		log.SetLevel(log.DEBUG)
		log.SetHeader(`${time_rfc3339} ${level} ${short_file} -`)
		log.Debug("log set to debug level")
	}

	readConfig()
	user := viper.GetString("username")
	pass := viper.GetString("password")
	token := viper.GetString("token")

	client := cloudcontrol.NewClient()

	if token == "" {
		createAndSaveSessionToken(user, pass, &client)
	} else {
		if body, err := client.ValidateSession(token); err != nil {
			log.Infof("invalid session token: %s", string(body))
			createAndSaveSessionToken(user, pass, &client)
		} else {
			log.Debug("session token is valid")
		}
	}

	if *listFlag {
		log.Info("listing available devices")
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
		log.Info("fetching device status")
		status, err := client.GetDeviceStatus()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Status: %s\n", types.Operate[status.Parameters.Operate])
		fmt.Printf("Mode: %s\n", types.ModesReverse[status.Parameters.OperationMode])
		fmt.Printf("Temperature: %0.1f\n", status.Parameters.TemperatureSet)
		fmt.Printf("Outside temperature: %0.1f\n", status.Parameters.OutsideTemperature)
		fmt.Printf("Fan speed: %s\n", types.FanSpeedReverse[status.Parameters.FanSpeed])
		fmt.Printf("Eco mode: %s\n", types.EcoModeReverse[status.Parameters.EcoMode])
	}

	if *historyFlag != "" {
		log.Infof("fetching historical data for %s\n", *historyFlag)
		history, err := client.GetDeviceHistory(types.HistoryDataMode[*historyFlag])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("#,AverageSettingTemp,AverageOutsideTemp,Consumption")
		for _, v := range history.HistoryEntries {
			fmt.Printf("%v,%v,%v,%v\n", v.DataNumber+1, v.AverageSettingTemp, v.AverageOutsideTemp, v.Consumption)
		}
	}

	if *onFlag {
		log.Info("turning device on")
		_, err := client.TurnOn()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("device turned on")
	}

	if *offFlag {
		log.Info("turning device off")
		_, err := client.TurnOff()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("device turned off")
	}

	if *tempFlag != 0 {
		log.Infof("setting temperature to %v degrees", *tempFlag)
		_, err := client.SetTemperature(*tempFlag)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("temperature set to %v degrees", *tempFlag)
	}

	if *fanSpeedFlag != "" {
		log.Infof("setting fan speed to %s", *fanSpeedFlag)
		_, err := client.SetFanSpeed(types.FanSpeed[*fanSpeedFlag])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("fan speed set to %s", *fanSpeedFlag)
	}

	if *modeFlag != "" {
		log.Infof("setting mode to %s", *modeFlag)
		_, err := client.SetMode(types.Modes[*modeFlag])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("mode set to %s", *modeFlag)
	}

	if *ecoModeFlag != "" {
		log.Infof("setting eco mode to %s", *ecoModeFlag)
		_, err := client.SetEcoMode(types.EcoMode[*ecoModeFlag])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("eco mode set to %s", *ecoModeFlag)
	}
}

func createAndSaveSessionToken(user string, pass string, client *cloudcontrol.Client) {
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
