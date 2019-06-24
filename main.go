package main

import (
	"fmt"

	"github.com/pytimer/kube-ipmi-plugin/pkg/ipmi"
)

func main() {
	fmt.Println(ipmi.PrintLANConfiguration())
}
