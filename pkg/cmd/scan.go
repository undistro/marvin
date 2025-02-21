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
	"github.com/go-logr/logr"
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
	NoFail                *bool
	SkipAnnotation        *string
	DisableAnnotationSkip *bool
	DisableZoraBanner     *bool
	CostLimit             *uint64

	ctx          context.Context
	log          logr.Logger
	printer      printers.Printer
	client       *dynamic.DynamicClient
	kubeVersion  *version.Info
	apiResources []*metav1.APIResourceList
	resources    map[string][]unstructured.Unstructured
	gvrs         map[string]string
}

// NewScanOptions returns a ScanOptions with the default values
func NewScanOptions() *ScanOptions {
	return &ScanOptions{
		ConfigFlags:           genericclioptions.NewConfigFlags(false),
		IOStreams:             genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		ChecksPath:            pointer.String(""),
		DisableBuiltIn:        pointer.Bool(false),
		OutputFormat:          pointer.String("table"),
		NoFail:                pointer.Bool(false),
		DisableAnnotationSkip: pointer.Bool(false),
		DisableZoraBanner:     pointer.Bool(false),
		SkipAnnotation:        pointer.String("marvin.undistro.io/skip"),
		CostLimit:             pointer.Uint64(1000000),
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
	if o.NoFail != nil {
		flags.BoolVar(o.NoFail, "no-fail", *o.NoFail, "Return an exit code of zero even if there are errors in the report")
	}
	if o.SkipAnnotation != nil {
		flags.StringVar(o.SkipAnnotation, "skip-annotation", *o.SkipAnnotation, "Annotation name for skipping checks")
	}
	if o.DisableAnnotationSkip != nil {
		flags.BoolVar(o.DisableAnnotationSkip, "disable-annotation-skip", *o.DisableAnnotationSkip, "Disable resource skipping by annotation")
	}
	if o.DisableZoraBanner != nil {
		flags.BoolVar(o.DisableZoraBanner, "disable-zora-banner", *o.DisableZoraBanner, "Disable Zora banner on output")
	}
	if o.CostLimit != nil {
		flags.Uint64Var(o.CostLimit, "cost-limit", *o.CostLimit, "CEL cost limit. Set 0 to disable it.")
	}
}

// Init initializes the kubernetes clients, get server version and API resources
func (o *ScanOptions) Init(ctx context.Context) error {
	if err := o.Validate(); err != nil {
		return err
	}
	o.ctx = ctx
	o.log = logr.FromContextOrDiscard(o.ctx)

	var printer printers.Printer
	switch *o.OutputFormat {
	case "json":
		printer = &printers.JSONPrinter{}
	case "yaml":
		printer = &printers.YAMLPrinter{}
	case "table":
		printer = &printers.TablePrinter{DisableZoraBanner: *o.DisableZoraBanner}
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
	o.gvrs = make(map[string]string)
	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ScanOptions) Validate() error {
	if *o.DisableBuiltIn == true && *o.ChecksPath == "" {
		return errors.New(`please set '--checks/-f' or keep 'disable-builtin' 'false'`)
	}
	return nil
}

// ToDynamicClient returns a DynamicClient using a computed RESTConfig.
func (o *ScanOptions) ToDynamicClient() (*dynamic.DynamicClient, error) {
	config, err := o.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(config)
}

// Run executes the scan command
func (o *ScanOptions) Run() (bool, error) {
	allChecks, err := o.getChecks()
	if err != nil {
		return false, err
	}
	report := types.NewReport(o.kubeVersion)
	for _, check := range allChecks {
		cr := o.runCheck(check)
		report.Add(cr)
	}

	report.GVRs = o.gvrs
	hasError := report.HasError()
	if hasError {
		o.log.Info("scan finished with errors")
	}
	return hasError, o.printer.PrintObj(*report, o.Out)
}

// getChecks returns a list of checks.Check based on the flags, including built-in checks or/and from a path.
func (o *ScanOptions) getChecks() ([]types.Check, error) {
	o.log.V(3).Info("loading checks", "builtin", !*o.DisableBuiltIn, "custom", *o.ChecksPath != "")
	var allChecks []types.Check
	if !*o.DisableBuiltIn {
		allChecks = loader.Builtins
		o.log.V(3).Info("builtin checks loaded", "total", len(loader.Builtins))
	}
	if *o.ChecksPath != "" {
		localChecks, err := loader.LoadChecks(*o.ChecksPath)
		if err != nil {
			o.log.Error(err, "failed to load checks", "path", *o.ChecksPath)
			return nil, fmt.Errorf("load checks error: %s", err.Error())
		}
		o.log.V(2).Info("custom checks loaded", "total", len(localChecks))
		allChecks = append(allChecks, localChecks...)
	}
	return allChecks, nil
}

// runCheck performs a check
func (o *ScanOptions) runCheck(check types.Check) *types.CheckResult {
	log := o.log.WithValues("check", check.ID)
	cr := types.NewCheckResult(check)
	defer cr.UpdateStatus()
	v, err := validator.Compile(check, o.apiResources, o.kubeVersion, *o.CostLimit)
	if err != nil {
		log.Error(err, "failed to compile check "+check.ID)
		cr.AddError(fmt.Errorf("%s compile error: %s", check.Path, err.Error()))
		return cr
	}
	log.V(3).Info("check compiled successfully")
	resources, errs := o.loadResources(check)
	cr.AddErrors(errs...)
	for gvr, objs := range resources {
		for _, obj := range objs {
			log := log.WithValues("obj", fmt.Sprintf("%s/%s", types.GVK(obj), types.NamespacedName(obj)))
			o.addGVR(obj, gvr)
			if o.isSkipped(check.ID, obj.GetAnnotations()) {
				log.V(4).Info("skipped")
				cr.AddSkipped(obj)
				continue
			}
			passed, _, err := v.Validate(obj, check.Params)
			if err != nil {
				log.Error(err, "failed to validate check "+check.ID)
				cr.AddError(fmt.Errorf("%s validate error: %s", check.Path, err.Error()))
				continue
			}
			if passed {
				log.V(4).Info("passed")
				cr.AddPassed(obj)
			} else {
				log.V(4).Info("failed")
				cr.AddFailed(obj)
			}
		}
	}
	return cr
}

// loadResources returns a map of resource slice by GVR to be validated by the given check
func (o *ScanOptions) loadResources(check types.Check) (map[string][]unstructured.Unstructured, []error) {
	resources := map[string][]unstructured.Unstructured{}
	var errs []error
	for _, r := range check.Match.Resources {
		gvr := r.ToGVR()
		gvrs := fmt.Sprintf("%s/%s", gvr.GroupVersion().String(), gvr.Resource)
		log := o.log.WithValues("check", check.ID)
		objs, cached := o.resources[gvrs]
		if cached {
			log.V(3).Info(gvrs+" resources cached", "total", len(objs))
			resources[gvrs] = objs
		} else {
			log.V(3).Info(fmt.Sprintf("listing %s from kubernetes", gvrs))
			ul, err := o.client.Resource(gvr).Namespace(*o.Namespace).List(o.ctx, metav1.ListOptions{})
			if err != nil {
				log.Error(err, "failed to list "+gvrs)
				errs = append(errs, fmt.Errorf("list %s error: %s", gvr.Resource, err.Error()))
				continue
			}
			log.V(1).Info(gvrs+" loaded from kubernetes", "total", len(ul.Items))
			o.resources[gvrs] = ul.Items
			resources[gvrs] = ul.Items
		}
	}
	return resources, errs
}

// isSkipped returns true if the checkID is annotated to be skipped
func (o *ScanOptions) isSkipped(checkID string, annotations map[string]string) bool {
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

// addGVR updates the map of GVR by GVK
func (o *ScanOptions) addGVR(obj unstructured.Unstructured, gvr string) {
	gvk := types.GVK(obj)
	if _, ok := o.gvrs[gvk]; !ok {
		o.gvrs[gvk] = gvr
	}
}
