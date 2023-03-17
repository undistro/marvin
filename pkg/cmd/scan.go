// Copyright 2023 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/utils/pointer"

	"github.com/spf13/pflag"

	"github.com/undistro/marvin/pkg/loader"
	"github.com/undistro/marvin/pkg/printers"
	"github.com/undistro/marvin/pkg/types"
	"github.com/undistro/marvin/pkg/validator"
)

type ScanOptions struct {
	*genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	ChecksPath            *string
	DisableBuiltIn        *bool
	OutputFormat          *string
	NoColor               *bool
	SkipAnnotation        *string
	DisableAnnotationSkip *bool

	printer      printers.Printer
	client       *dynamic.DynamicClient
	kubeVersion  *version.Info
	apiResources []*metav1.APIResourceList
	resources    map[string][]unstructured.Unstructured
}

func NewScanOptions() *ScanOptions {
	return &ScanOptions{
		ConfigFlags:           genericclioptions.NewConfigFlags(false),
		IOStreams:             genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		ChecksPath:            pointer.String(""),
		DisableBuiltIn:        pointer.Bool(false),
		OutputFormat:          pointer.String("table"),
		NoColor:               pointer.Bool(false),
		DisableAnnotationSkip: pointer.Bool(false),
		SkipAnnotation:        pointer.String("marvin.undistro.io/skip"),
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
		flags.StringVarP(o.OutputFormat, "output", "o", *o.OutputFormat, `Output format. One of: ("table", "json", "yaml" or "markdown")`)
	}
	if o.NoColor != nil {
		flags.BoolVar(o.NoColor, "no-color", *o.NoColor, "Disable color output")
	}
	if o.SkipAnnotation != nil {
		flags.StringVar(o.SkipAnnotation, "skip-annotation", *o.SkipAnnotation, "Annotation name for skipping checks")
	}
	if o.DisableAnnotationSkip != nil {
		flags.BoolVar(o.DisableAnnotationSkip, "disable-annotation-skip", *o.DisableAnnotationSkip, "Disable resource skipping by annotation")
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
	color.NoColor = *o.NoColor

	var printer printers.Printer
	switch *o.OutputFormat {
	case "json":
		printer = &printers.JSONPrinter{}
	case "yaml":
		printer = &printers.YAMLPrinter{}
	case "table":
		printer = &printers.TablePrinter{}
	case "markdown":
		color.NoColor = true
		printer = &printers.MarkdownPrinter{}
	default:
		return fmt.Errorf("invalid output format '%s'", *o.OutputFormat)
	}

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

	o.client = dynamicClient
	o.kubeVersion = kubeVersion
	o.apiResources = apiResources
	o.printer = printer
	o.resources = make(map[string][]unstructured.Unstructured)
	return nil
}

func (o *ScanOptions) Run() error {
	allChecks, err := o.getChecks()
	if err != nil {
		return err
	}
	report := types.NewReport(o.kubeVersion)
	for _, check := range allChecks {
		cr := types.NewCheckResult(check)
		report.Add(cr)
		v, err := validator.Compile(check, o.apiResources, o.kubeVersion)
		if err != nil {
			cr.AddError(fmt.Errorf("compile error: %s", err.Error()))
			continue
		}
		resources, errs := o.loadResources(check)
		cr.AddErrors(errs...)
		for _, obj := range resources {
			if o.IsSkipped(check.ID, obj.GetAnnotations()) {
				cr.AddSkipped(obj)
				continue
			}
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
		cr.UpdateStatus()
	}

	return o.printer.PrintObj(*report, o.Out)
}

// loadResources returns the resources to be validated by the given check
func (o *ScanOptions) loadResources(check types.Check) ([]unstructured.Unstructured, []error) {
	var resources []unstructured.Unstructured
	var errs []error
	for _, r := range check.Match.Resources {
		gvr := r.ToGVR()
		objs, cached := o.resources[gvr.String()]
		if cached {
			resources = append(resources, objs...)
		} else {
			ul, err := o.client.Resource(gvr).Namespace(*o.Namespace).List(context.Background(), metav1.ListOptions{})
			if err != nil {
				errs = append(errs, fmt.Errorf("list %s error: %s", gvr.Resource, err.Error()))
				continue
			}
			o.resources[gvr.String()] = ul.Items
			resources = append(resources, ul.Items...)
		}
	}
	return resources, errs
}

// getChecks returns a list of checks.Check based on the flags, including built-in checks or/and from a path.
func (o *ScanOptions) getChecks() ([]types.Check, error) {
	var allChecks []types.Check
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

func (o *ScanOptions) IsSkipped(checkID string, annotations map[string]string) bool {
	if annotations == nil {
		return false
	}
	if *o.DisableAnnotationSkip {
		return false
	}
	v, ok := annotations[*o.SkipAnnotation]
	if !ok {
		return false
	}
	ids := strings.Split(v, ",")
	for _, s := range ids {
		if strings.TrimSpace(s) == checkID {
			return true
		}
	}
	return false
}
