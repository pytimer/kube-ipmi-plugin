package ipmi

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pytimer/kube-ipmi-plugin/pkg/constants"

	"k8s.io/klog"
	"k8s.io/utils/path"
)

func rawDecode(raw string) map[string]string {
	m := make(map[string]string)
	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		items := strings.SplitN(line, ":", 2)
		if len(items) < 2 {
			continue
		}
		title := strings.TrimSpace(items[0])
		value := items[1]
		if title == "" || title == "Auth Type Enable" || title == "Cipher Suite Priv Max" {
			continue
		}
		m[title] = strings.TrimSpace(value)
	}

	return m
}

func filterUsedConfiguration(configurations map[string]string, filters []string) (used map[string]string) {
	for _, k := range filters {
		v, ok := configurations[k]
		if !ok {
			continue
		}
		if used == nil {
			used = make(map[string]string)
		}
		used[k] = v
	}
	return
}

func PrintLANConfiguration(toolPath string) (map[string]string, error) {
	out, err := exec.Command(toolPath, "-I", "open", "lan", "print").CombinedOutput()
	if err != nil {
		return nil, err
	}

	info := rawDecode(string(out))
	used := filterUsedConfiguration(info, constants.IPMIConfigurationKeys)

	return used, nil
}

func CheckIPMIToolPath(toolPath string) (string, error) {
	var err error
	if toolPath == "" {
		klog.Warningf("The ipmitool path is empty, we should check the default ipmitool path.")
		toolPath, err = exec.LookPath(constants.DefaultIPMIToolPath)
		if err != nil {
			return "", err
		}
	}

	exists, err := path.Exists(path.CheckFollowSymlink, toolPath)
	if err != nil {
		return toolPath, err
	}
	if !exists {
		return toolPath, fmt.Errorf("impitool path [%s] not exists", toolPath)
	}

	return toolPath, nil
}
