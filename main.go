package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Viperoo/golog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var ServerTypes = map[string]string{
	"150sk10": "KS-1",
	"150sk20": "KS-2a",
	"150sk21": "KS-2b",
	"150sk22": "KS-2c",
	"150sk30": "KS-3",
	"150sk31": "KS-3",
	"150sk40": "KS-4",
	"150sk41": "KS-4",
	"150sk42": "KS-4",
	"150sk50": "KS-5",
	"150sk60": "KS-6",
}

var DataCenters = map[string]string{
	"bhs": "Beauharnois, Canada (Americas)",
	"gra": "Gravelines, France",
	"rbx": "Roubaix, France (Western Europe)",
	"sbg": "Strasbourg, France (Central Europe)",
	"par": "Paris, France",
}

type Kimsufi struct {
	Answer struct {
		Availability []struct {
			Class     string `json:"__class"`
			Reference string `json:"reference"`
			Zones     []struct {
				Availability string `json:"availability"`
				Class        string `json:"__class"`
				Zone         string `json:"zone"`
			} `json:"zones"`
		} `json:"availability"`
		Class string `json:"__class"`
	} `json:"answer"`
	Error   interface{} `json:"error"`
	Id      int64       `json:"id"`
	Version string      `json:"version"`
}

type AvalibleServers struct {
	Server, Zone string
}

var logger log.Logger
var logfile = flag.String("l", "kimsufi.log", "Log file")
var debug = flag.Bool("d", false, "Debug mode")
var configfile = flag.String("c", "kimsufi.conf", "Configuration file")

func main() {
	/*
	* Parse flags
	 */
	flag.Parse()
	/*
	* Set logger level
	 */
	setLogger()

	loadConfig(*configfile)

	response, err := http.Get("https://ws.ovh.com/dedicated/r2/ws.dispatcher/getAvailability2")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Warn(err.Error())
		}
		var m Kimsufi
		err = json.Unmarshal(contents, &m)

		if err != nil {
			logger.Warn(err.Error())
			os.Exit(1)
		}
		Avalible := make(map[int]AvalibleServers)

		for _, value := range m.Answer.Availability {

			if serverName, ok := ServerTypes[value.Reference]; ok {
				for _, zone := range value.Zones {
					if zone.Availability == "1H-high" {

						Avalible[len(Avalible)+1] = AvalibleServers{Server: serverName, Zone: zone.Zone}
						//fmt.Printf("Serwer jest dostępny ! %[2]d ------ %[1]d\n", serverName, zone)
					}
				}
			}

		}

		FindServers(Avalible)

	}
}

func FindServers(Servers map[int]AvalibleServers) {
	for _, server := range Servers {
		for _, s := range Config.Notifier.Server {

			if s == server.Server {
				fmt.Print("Ślpij gołębia.")

			}

		}

	}
}

func setLogger() {
	file, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Print("Error log is not wrtiable.")
		os.Exit(1)
	}
	var multi io.Writer
	if *debug == true {
		multi = io.MultiWriter(file, os.Stdout)
	} else {
		multi = io.MultiWriter(file)
	}

	logger, _ = log.NewLogger(multi,
		log.TIME_FORMAT_SEC,
		log.LOG_FORMAT_SIMPLE,
		log.LogLevel_Debug)
}
