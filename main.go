package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pytimer/kube-ipmi-plugin/pkg/constants"
	"github.com/pytimer/kube-ipmi-plugin/pkg/ipmi"
	"github.com/pytimer/kube-ipmi-plugin/pkg/kube"
	"github.com/pytimer/kube-ipmi-plugin/pkg/util"

	"github.com/spf13/pflag"
	"k8s.io/klog"
)

type PluginOptions struct {
	IPMIToolPath string
	KubeConfig   string
	Period       string
	NodeName     string
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
	pflag.StringVar(&po.Period, constants.PeriodFlagName, constants.DefaultPeriod, "The application worker period.")
	pflag.StringVar(&po.NodeName, constants.NodeNameFlagName, "", "The ipmi information patch to the node name.")
	pflag.Parse()

	if err := po.run(); err != nil {
		klog.Fatal(err)
	}
}

func (o *PluginOptions) run() error {
	period, err := util.GetTimeDurationStringToSeconds(o.Period)
	if err != nil {
		return fmt.Errorf("failed to convert period to time, %v", err)
	}

	if o.NodeName == "" {
		klog.Info("The --nodename is empty, so get the node name by env or hostname.")
		nodeName, err := util.GetNodeName()
		if err != nil {
			return fmt.Errorf("failed to get the node name, %v", err)
		}
		o.NodeName = nodeName
	}

	c, err := kube.NewClient(o.KubeConfig)
	if err != nil {
		return fmt.Errorf("failed to generate kubernetes client, %v", err)
	}

	shutdownC := make(chan struct{})
	go listenToSystemSignal(shutdownC)
	p := ipmi.NewPlugin(o.IPMIToolPath, period, o.NodeName, c, shutdownC)

	return p.Run()
}

// listenToSystemSignal listen system signal and exit.
func listenToSystemSignal(stopC chan<- struct{}) {
	klog.V(5).Info("Listen to system signal.")

	signalChan := make(chan os.Signal, 1)
	ignoreChan := make(chan os.Signal, 1)

	signal.Notify(ignoreChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		klog.V(3).Infof("Shutdown by system signal: %s", sig)
		stopC <- struct{}{}
	}
}
