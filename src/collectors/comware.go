package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type deviceComware struct {
	configuration.DeviceConfig
}

func makeComware(d configuration.DeviceConfig) Collector {
	return &deviceComware{d}
}

// BuildBatcher for Comware
func (p *deviceComware) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"<.*.>", "screen-length disable"},
		{"<.*.>", "display current-configuration"},
		{"return"},
	})
}

// ParseResult for Comware
func (p *deviceComware) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`#[\s\S]*?return`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
