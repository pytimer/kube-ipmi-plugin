package constants


var IPMIConfigurationKeys = []string {
	"IP Address Source",
	"IP Address",
	"Subnet Mask",
	"MAC Address",
	"Default Gateway IP",
	"Default Gateway MAC",
}

const (
	IPMIToolPathFlagName = "ipmitool-path"
	KubeConfigFlagName = "kubeconfig"

	DefaultIPMIToolPath = "ipmitool"
	DefaultKubeConfigFile = "~/.kube/config"
)