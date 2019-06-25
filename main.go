package main

import (
	"flag"
	"fmt"

	"github.com/pytimer/kube-ipmi-plugin/pkg/constants"
	"github.com/pytimer/kube-ipmi-plugin/pkg/ipmi"

	"github.com/spf13/pflag"
	"k8s.io/klog"
)

type PluginOptions struct {
	IPMIToolPath string
	KubeConfig   string
}

func main() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Set("logtostderr", "true")
	// We do not want these flags to show up in --help
	// These MarkHidden calls must be after the lines above
	pflag.CommandLine.MarkHidden("version")
	pflag.CommandLine.MarkHidden("log_flush_frequency")
	pflag.CommandLine.MarkHidden("alsologtostderr")
	pflag.CommandLine.MarkHidden("log_backtrace_at")
	pflag.CommandLine.MarkHidden("log_dir")
	pflag.CommandLine.MarkHidden("logtostderr")
	pflag.CommandLine.MarkHidden("stderrthreshold")
	pflag.CommandLine.MarkHidden("vmodule")
	pflag.CommandLine.MarkHidden("log_file")
	pflag.CommandLine.MarkHidden("log_file_max_size")
	pflag.CommandLine.MarkHidden("skip_headers")
	pflag.CommandLine.MarkHidden("skip_log_headers")

	po := &PluginOptions{}
	pflag.StringVar(&po.IPMIToolPath, constants.IPMIToolPathFlagName, "", "Path to the ipmitool")
	pflag.StringVar(&po.KubeConfig, constants.KubeConfigFlagName, constants.DefaultKubeConfigFile, "The kubeconfig use connect to the Kubernetes cluster.")
	pflag.Parse()

	toolPath, err := ipmi.CheckIPMIToolPath(po.IPMIToolPath)
	if err != nil {
		klog.Fatalf("Failed to check ipmitool path, %v", err)
	}

	fmt.Println(ipmi.PrintLANConfiguration(toolPath))
}
