package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/host"
)

/* Actually almost all parameters are mandatory.
1. In case of Docker, we exactly know what to send, so just pass as a command line arguments.
No need to detect anything.
2. In case of packages, we again know exactly what to send as the package is dedicated for particular OS.

Instance ID is optional, however advised to be passed as cmd line arg as well.
1. If the parameter is omitted, gopsutil.GetHost() does
	1. Get UUID from /sys/class/dmi/id/product_uuid (needs root access)
	2. Get UUID from machine-id
	3. Get UUID from sys/kernel/random/boot_id (this ID changes on every boot, so not ideal)
1. For Docker container 1, 2, 3 are not available. We can detect if /telemetry_uuid file exists.
If not, create it with random uuid and pass its content as cmd line arg.
2. For baremetal installation 1, 2 should be available, but still advised to generate UUID
externally and pass as cmd line param to avoid fallback to 3.
*/

var (
	instanceId = kingpin.Flag(
		"instanceId",
		"Instance ID",
	).Short('i').String()
	productFamily = kingpin.Flag(
		"productFamily",
		"Product family",
	).Short('f').Required().String()
	osName = kingpin.Flag(
		"osName",
		"Operating system name",
	).Short('o').Required().String()
	hwArchitecture = kingpin.Flag(
		"hwArchitecture",
		"Hardware architecture",
	).Short('h').Required().String()
	productVersion = kingpin.Flag(
		"productVersion",
		"Product version",
	).Short('v').Required().String()
	telemetryAPI = kingpin.Flag(
		"telemetryApi",
		"Telemetry API endpoint",
	).Short('d').Default("http://localhost:8081/v1/telemetry/GenericReport").String()
)

type telemetryMetric struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type telemetryReport struct {
	Id            string            `json:"id"`
	Time          string            `json:"time"`
	InstanceId    string            `json:"instanceId"`
	ProductFamily string            `json:"productFamily"`
	Metrics       []telemetryMetric `json:"metrics"`
}

type telemetryMessage struct {
	Reports []telemetryReport `json:"reports"`
}

func main() {
	kingpin.Parse()

	// handle optional params
	if *instanceId == "" {
		// Instance ID was not provided, figure out something
		id, err := host.HostID()

		if err != nil {
			// In case of err, id is an empty string, so probably need to figure out something more here
			// like random GUID
			id = uuid.New().String()
		}
		*instanceId = id
	}

	// collect
	metrics := []telemetryMetric{
		{
			Key:   "version",
			Value: *productVersion,
		},
		{
			Key:   "osName",
			Value: *osName,
		},
		{
			Key:   "hwArch",
			Value: *hwArchitecture,
		},
	}

	reportId := uuid.New()
	instId := uuid.MustParse(*instanceId)
	report := telemetryReport{
		Id:            b64.StdEncoding.EncodeToString(reportId[:]),
		Time:          time.Now().UTC().Format(time.RFC3339Nano),
		InstanceId:    b64.StdEncoding.EncodeToString(instId[:]),
		ProductFamily: *productFamily,
		Metrics:       metrics,
	}

	var message telemetryMessage
	message.Reports = append(message.Reports, report)

	// json
	JSON, err := json.Marshal(message)
	if err != nil {
		log.Fatal("impossible to create json: ", err)
		return
	}

	// this is just for debug
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, JSON, "", "\t")
	log.Println(string(prettyJSON.Bytes()))

	// post
	req, err := http.NewRequest("POST", *telemetryAPI, bytes.NewReader(JSON))
	if err != nil {
		log.Fatal("impossible to build request: ", err)
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-Status", "0")

	client := http.Client{Timeout: 30 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("impossible to send request: ", err)
		return
	}
	log.Println("status Code:", res.StatusCode)

	// do we care about response body?
	res.Body.Close()
}
