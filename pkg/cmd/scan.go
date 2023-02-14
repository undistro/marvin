package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/utils/pointer"
)

type ScanOptions struct {
	*genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	ChecksPath *string
}

func NewScanOptions() *ScanOptions {
	return &ScanOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(false),
		IOStreams:   genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		ChecksPath:  pointer.String("checks/"),
	}
}

func (o *ScanOptions) AddFlags(flags *pflag.FlagSet) {
	o.ConfigFlags.AddFlags(flags)
	if o.ChecksPath != nil {
		flags.StringVar(o.ChecksPath, "checks", *o.ChecksPath, "Path to the check files directory")
	}
}

func (o *ScanOptions) ToDynamicClient() (*dynamic.DynamicClient, error) {
	config, err := o.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

func (o *ScanOptions) Run() error {
	//dynamicClient, err := o.ToDynamicClient()
	//if err != nil {
	//	return err
	//}
	discoveryClient, err := o.ToDiscoveryClient()
	if err != nil {
		return err
	}
	kubeVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return err
	}
	_, apiResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return err
	}

	//TODO test
	fmt.Println(*o.ChecksPath)
	fmt.Println(kubeVersion)
	fmt.Println(len(apiResources))

	return nil
}
