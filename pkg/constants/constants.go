package constants

import "time"

var IPMIConfigurationKeys = []string{
	"IP Address Source",
	"IP Address",
	"Subnet Mask",
	"MAC Address",
	"Default Gateway IP",
	"Default Gateway MAC",
}

const (
	IPMIToolPathFlagName = "ipmitool-path"
	KubeConfigFlagName   = "kubeconfig"
	PeriodFlagName       = "period"
	NodeNameFlagName     = "nodename"

	DefaultIPMIToolPath   = "ipmitool"
	DefaultKubeConfigFile = "~/.kube/config"
	DefaultPeriod         = "1h"

	// APICallRetryInterval defines how long plugin should wait before retrying a failed API operation
	APICallRetryInterval = 500 * time.Millisecond
	// PatchNodeTimeout specifies how long plugin should wait for applying the label and taint on the master before timing out
	PatchNodeTimeout = 2 * time.Minute

	NodeNameEnvName = "NODENAME"

	IPMIAnnotationKey = "ipmi.alpha.kubernetes.io/net"
)
