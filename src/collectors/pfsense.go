package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type devPfsense struct {
	configuration.DeviceConfig
}

func makePfsense(d configuration.DeviceConfig) Collector {
	return &devPfsense{d}
}

// BuildBatcher for pfSense
func (p *devPfsense) BuildBatcher() ([]expect.Batcher, error) {
	// The "OK" result must be the first entry for variable.
	// The more the better, since this is constantly checking every case for a match.
	// - So simply having .*root will match multiple times throughout the dump.
	return utils.VariableBatcher([][]string{
		{`</pfsense>`}, // Found the "OK" result!
		{`Enter an option: `, "8\n"},
		{`/root(\x1b\[[0-9;]*m)?:`, "cat /conf/config.xml\n"},
		{`\$ `, "cat /conf/config.xml\n"},
	})
}

// ParseResult for pfSense
func (p *devPfsense) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
