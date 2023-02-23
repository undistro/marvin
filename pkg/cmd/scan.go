package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/utils/pointer"

	"github.com/spf13/pflag"

	"github.com/undistro/marvin/pkg/checks"
	"github.com/undistro/marvin/pkg/loader"
	"github.com/undistro/marvin/pkg/printers"
	"github.com/undistro/marvin/pkg/report"
	"github.com/undistro/marvin/pkg/validator"
)

type ScanOptions struct {
	*genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	ChecksPath     *string
	DisableBuiltIn *bool
	OutputFormat   *string

	printer      printers.Printer
	client       *dynamic.DynamicClient
	kubeVersion  *version.Info
	apiResources []*metav1.APIResourceList
}

func NewScanOptions() *ScanOptions {
	return &ScanOptions{
		ConfigFlags:    genericclioptions.NewConfigFlags(false),
		IOStreams:      genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		ChecksPath:     pointer.String(""),
		DisableBuiltIn: pointer.Bool(false),
		OutputFormat:   pointer.String("table"),
	}
}

// AddFlags binds scan configuration flags to a given flagset
func (o *ScanOptions) AddFlags(flags *pflag.FlagSet) {
	o.ConfigFlags.AddFlags(flags)
	if o.ChecksPath != nil {
		flags.StringVarP(o.ChecksPath, "checks", "f", *o.ChecksPath, "Path to the check files directory")
	}
	if o.DisableBuiltIn != nil {
		flags.BoolVar(o.DisableBuiltIn, "disable-builtin", *o.DisableBuiltIn, "Disable builtin checks")
	}
	if o.OutputFormat != nil {
		flags.StringVarP(o.OutputFormat, "output", "o", *o.OutputFormat, `Output format. One of: ('table', json', 'yaml').`)
	}
}

func (o *ScanOptions) ToDynamicClient() (*dynamic.DynamicClient, error) {
	config, err := o.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

// Validate ensures that all required arguments and flag values are provided
func (o *ScanOptions) Validate() error {
	if *o.DisableBuiltIn == true && *o.ChecksPath == "" {
		return errors.New(`please set '--checks/-f' or keep 'disable-builtin' 'false'`)
	}
	return nil
}

// Init initializes the kubernetes clients, get server version and API resources
func (o *ScanOptions) Init() error {
	dynamicClient, err := o.ToDynamicClient()
	if err != nil {
		return fmt.Errorf("dynamic client error: %s", err.Error())
	}
	discoveryClient, err := o.ToDiscoveryClient()
	if err != nil {
		return fmt.Errorf("kubernetes client error: %s", err.Error())
	}
	kubeVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return fmt.Errorf("server version error: %s", err.Error())
	}
	_, apiResources, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return fmt.Errorf("server groups error: %s", err.Error())
	}

	var printer printers.Printer
	switch *o.OutputFormat {
	case "json":
		printer = &printers.JSONPrinter{}
	case "yaml":
		printer = &printers.YAMLPrinter{}
	case "table":
		printer = &printers.TablePrinter{}
	default:
		return fmt.Errorf("invalid output format '%s'", *o.OutputFormat)
	}

	o.client = dynamicClient
	o.kubeVersion = kubeVersion
	o.apiResources = apiResources
	o.printer = printer
	return nil
}

func (o *ScanOptions) Run() error {
	allChecks, err := o.getChecks()
	if err != nil {
		return err
	}
	rep := report.New(o.kubeVersion)
	cache := make(map[string][]unstructured.Unstructured)
	for _, check := range allChecks {
		cr := report.NewCheckResult(check)
		rep.Add(cr)
		v, err := validator.Compile(check, o.apiResources, o.kubeVersion)
		if err != nil {
			cr.AddError(fmt.Errorf("compile error: %s", err.Error()))
			continue
		}
		var resources []unstructured.Unstructured
		for _, r := range check.Match.Resources {
			gvr := r.ToGVR()
			objs, cached := cache[gvr.String()]
			if cached {
				resources = append(resources, objs...)
			} else {
				ul, err := o.client.Resource(gvr).Namespace(*o.Namespace).List(context.Background(), metav1.ListOptions{})
				if err != nil {
					cr.AddError(fmt.Errorf("list %s error: %s", gvr.Resource, err.Error()))
					continue
				}
				cache[gvr.String()] = ul.Items
				resources = append(resources, ul.Items...)
			}
		}
		for _, obj := range resources {
			passed, _, err := v.Validate(obj, check.Params)
			if err != nil {
				cr.AddError(fmt.Errorf("%s validate error: %s", check.Path, err.Error()))
				continue
			}
			if passed {
				cr.AddPassed(obj)
			} else {
				cr.AddFailed(obj)
			}
		}
	}

	return o.printer.PrintObj(*rep, o.Out)
}

// getChecks returns a list of checks.Check based on the flags, including built-in checks or/and from a path.
func (o *ScanOptions) getChecks() ([]checks.Check, error) {
	var allChecks []checks.Check
	if !*o.DisableBuiltIn {
		allChecks = loader.Builtins
	}
	if *o.ChecksPath != "" {
		localChecks, err := loader.LoadChecks(*o.ChecksPath)
		if err != nil {
			return nil, fmt.Errorf("load checks error: %s", err.Error())
		}
		allChecks = append(allChecks, localChecks...)
	}
	return allChecks, nil
}
