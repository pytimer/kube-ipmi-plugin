package kube

import (
	"encoding/json"

	"github.com/pytimer/kube-ipmi-plugin/pkg/constants"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

// NewClient creates k8s client
func NewClient(kubeConfigFile string) (kubernetes.Interface, error) {
	var clientConfig *rest.Config
	var err error
	if len(kubeConfigFile) > 0 {
		clientConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigFile)
		if err != nil {
			return nil, err
		}
	} else {
		clientConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// PatchNode tries to patch a node using the following client, executing patchFn for the actual mutating logic
func PatchNode(client kubernetes.Interface, nodeName string, patchFn func(*corev1.Node)) error {
	// Loop on every false return. Return with an error if raised. Exit successfully if true is returned.
	return wait.Poll(constants.APICallRetryInterval, constants.PatchNodeTimeout, func() (bool, error) {
		n, err := client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		oldData, err := json.Marshal(n)
		if err != nil {
			return false, err
		}

		patchFn(n)

		newData, err := json.Marshal(n)
		if err != nil {
			return false, err
		}

		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, corev1.Node{})
		if err != nil {
			return false, err
		}

		if _, err := client.CoreV1().Nodes().Patch(n.Name, types.StrategicMergePatchType, patchBytes); err != nil {
			if apierrors.IsConflict(err) {
				klog.Warning("[patchnode] Temporarily unable to update node metadata due to conflict (will retry)")
				return false, nil
			}
			return false, err
		}

		return true, nil
	})
}
