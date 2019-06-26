package ipmi

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pytimer/kube-ipmi-plugin/pkg/constants"
	"github.com/pytimer/kube-ipmi-plugin/pkg/kube"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"k8s.io/utils/path"
)

type Plugin struct {
	toolPath  string
	period    time.Duration
	nodeName  string
	k8sClient kubernetes.Interface
	stopC     chan struct{}
	data      map[string]string
}

func NewPlugin(toolPath string, loopPeriod int, nodeName string, c kubernetes.Interface, ch chan struct{}) *Plugin {
	return &Plugin{
		toolPath:  toolPath,
		period:    time.Duration(loopPeriod) * time.Second,
		nodeName:  nodeName,
		k8sClient: c,
		stopC:     ch,
	}
}

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

// PrintLANConfiguration using `ipmitool -I open lan print` to get the necessary information,
// such as ipmi ip address, gateway, etc.
func (p *Plugin) PrintLANConfiguration() (map[string]string, error) {
	out, err := exec.Command(p.toolPath, "-I", "open", "lan", "print").CombinedOutput()
	if err != nil {
		return nil, err
	}

	info := rawDecode(string(out))
	used := filterUsedConfiguration(info, constants.IPMIConfigurationKeys)

	return used, nil
}

// CheckIPMIToolPath check ipmitool exists, if the specified path not exists, use default ipmitool.
// If the ipmitool not installed, return error.
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

func (p *Plugin) generateLanConfigurationAnnotations() (string, error) {
	annotations := make(map[string]string)
	for k, v := range p.data {
		newKey := strings.Replace(strings.ToLower(k), " ", "_", -1)
		annotations[newKey] = v
	}

	b, err := json.Marshal(annotations)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (p *Plugin) patchLanConfigurationToAnnotations(n *corev1.Node) {
	v, err := p.generateLanConfigurationAnnotations()
	if err != nil {
		klog.Warningf("Failed to generate ipmi information annotation, %v", err)
		return
	}
	n.Annotations[constants.IPMIAnnotationKey] = v
}

// Run business logic
func (p *Plugin) Run() error {
	wait.Until(func() {
		toolPath, err := CheckIPMIToolPath(p.toolPath)
		if err != nil {
			klog.Warningf("Failed to check ipmitool path, %v", err)
			return
		}
		p.toolPath = toolPath

		configs, err := p.PrintLANConfiguration()
		if err != nil {
			klog.Warningf("Failed to get the ipmi information, %v", err)
			return
		}
		p.data = configs

		if err := kube.PatchNode(p.k8sClient, p.nodeName, p.patchLanConfigurationToAnnotations); err != nil {
			klog.Warningf("Failed to patch ipmi to the node %s, %v", p.nodeName, err)
			return
		}
	}, p.period, p.stopC)

	return nil
}
