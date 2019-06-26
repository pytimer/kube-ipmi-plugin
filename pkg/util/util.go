package util

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/pytimer/kube-ipmi-plugin/pkg/constants"
)

var (
	ValueFormatRegex = regexp.MustCompile("^\\d+")
	UnitFormatRegex  = regexp.MustCompile("\\w{1}$")
)

var unitMultiplier = map[string]int{
	"s": 1,
	"m": 60,
	"h": 3600,
}

// GetTimeDurationStringToSeconds ...
func GetTimeDurationStringToSeconds(str string) (int, error) {
	multiplier := 1

	matches := ValueFormatRegex.FindAllString(str, 1)

	if len(matches) <= 0 {
		return 0, fmt.Errorf("Time Duration could not be parsed")
	}

	value, err := strconv.Atoi(matches[0])
	if err != nil {
		return 0, err
	}

	unit := UnitFormatRegex.FindAllString(str, 1)[0]

	if val, ok := unitMultiplier[unit]; ok {
		multiplier = val
	}

	return int(value * multiplier), nil
}

func GetNodeName() (string, error) {
	if os.Getenv(constants.NodeNameEnvName) != "" {
		return os.Getenv(constants.NodeNameEnvName), nil
	}
	return os.Hostname()
}
