package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"fmt"
	"log"
	"strconv"
	"time"
)

// CollectorWorker write me.
type CollectorWorker struct {
	RunConfig *configuration.Config
}

// RunCollectors runs all collectors
func (cw *CollectorWorker) RunCollectors() {
	for _, device := range cw.RunConfig.Devices {
		if device.Disabled {
			log.Printf("Config disabled: %v", device.Name)
			continue
		}

		cw.Run(device)
	}

	fmt.Printf("Completed collections: %d\n", len(cw.RunConfig.Devices))
}

// Run the collector for this device.
func (cw *CollectorWorker) Run(device configuration.DeviceConfig) {
	fmt.Printf("Collect Start: %s\n", device.Name)

	collector, _ := MakeCollector(device)

	batchSlice, _ := collector.BuildBatcher()

	// Set up SSHConfig
	s := &utils.SSHConfig{
		User: device.User,
		Pass: device.Pass,
		Host: device.Host + ":" + strconv.Itoa(device.Port),
	}

	// Special case... only some collectors need to make some modifications.
	if collectorSpecial, ok := collector.(CollectorSpecial); ok {
		collectorSpecial.ModifySSHConfig(s)
	}

	connection, err := utils.SSHClient(*s)

	// call GatherExpect to collect the configs
	// TODO: Verify pointer/reference/dereference is necessary.
	result, err := utils.GatherExpect(batchSlice, time.Second*10, connection)
	if err != nil {
		panic(err)
	}

	// Grab just the last result.
	lastResult := result[len(result)-1].Output
	parsed, _ := collector.ParseResult(lastResult)

	log.Printf("Writing: %s\nLength: %d\n", device.Name, len(parsed))

	utils.WriteFile(parsed, device.Name+".txt")
}
