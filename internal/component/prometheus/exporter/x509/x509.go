// package x509 provides a prometheus.exporter x509 certificate component
// Maintainers for the Grafana Alloy wrapper:
// - @cstdev
package x509

import (
	"time"

	x509Exporter "github.com/enix/x509-certificate-exporter/v3/pkg/exporter"
	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/component/prometheus/exporter"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/internal/static/integrations"
	"github.com/grafana/alloy/internal/static/integrations/x509_exporter"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.x509",
		Stability: featuregate.StabilityExperimental,
		Community: true,
		Args:      Arguments{},
		Exports:   exporter.Exports{},

		Build: exporter.New(createExporter, "x509"),
	})
}

func createExporter(opts component.Options, args component.Arguments, defaultInstanceKey string) (integrations.Integration, string, error) {
	a := args.(Arguments)
	return integrations.NewIntegrationWithInstanceKey(opts.Logger, a.Convert(), defaultInstanceKey)
}

// Arguments holds values which are used to configured the prometheus.exporter.x509 component.
type Arguments struct {
	Files                 []string                   `alloy:"files,attr,optional"`
	Directories           []string                   `alloy:"directories,attr,optional"`
	YAMLs                 []string                   `alloy:"yamls,attr,optional"`
	YAMLPaths             []x509Exporter.YAMLCertRef `alloy:"yaml_paths,attr,optional"`
	TrimPathComponents    int                        `alloy:"trim_path_components,attr,optional"`
	MaxCacheDuration      time.Duration              `alloy:"max_cache_duration,attr,optional"`
	ExposeRelativeMetrics bool                       `alloy:"expose_relative_metrics,attr,optional"`
	ExposeErrorMetrics    bool                       `alloy:"expose_error_metrics,attr,optional"`
	ExposeLabels          []string                   `alloy:"expose_labels,attr,optional"`
	ConfigMapKeys         []string                   `alloy:"config_map_keyss,attr,optional"`
	Kubernetes            KubeBlock                  `alloy:"kubernetes,block,optional`
}

type KubeBlock struct {
	KubeConfigPath             string                        `alloy:"kube_config_path,attr,optional"`
	KubeSecretTypes            []x509Exporter.KubeSecretType `alloy:"kube_secret_types,attr,optional"`
	KubeIncludeNamespaces      []string                      `alloy:"kube_include_namespaces,attr,optional"`
	KubeExcludeNamespaces      []string                      `alloy:"kube_exclude_namespaces,attr,optional"`
	KubeIncludeNamespaceLabels []string                      `alloy:"kube_include_namespace_labels,attr,optional"`
	KubeExcludeNamespaceLabels []string                      `alloy:"kube_exclude_namespace_labels,attr,optional"`
	KubeIncludeLabels          []string                      `alloy:"kube_include_labels,attr,optional"`
	KubeExcludeLabels          []string                      `alloy:"kube_exclude_labels,attr,optional"`
}

// Exports holds the values exported by the prometheus.exporter.x509 component.
type Exports struct{}

// SetToDefault implements syntax.Defaulter
func (args *Arguments) SetToDefault() {
	*args = Arguments{}
}

func (a *Arguments) Convert() *x509_exporter.Config {
	kubeEnabled := false
	if a.Kubernetes.KubeConfigPath != "" {
		kubeEnabled = true
	}

	return &x509_exporter.Config{
		Directories:                a.Directories,
		Files:                      a.Files,
		YAMLs:                      a.YAMLs,
		TrimPathComponents:         a.TrimPathComponents,
		MaxCacheDuration:           a.MaxCacheDuration,
		ExposeRelativeMetrics:      a.ExposeRelativeMetrics,
		ExposeErrorMetrics:         a.ExposeErrorMetrics,
		ExposeLabels:               a.ExposeLabels,
		ConfigMapKeys:              a.ConfigMapKeys,
		KubeEnabled:                kubeEnabled,
		KubeConfigPath:             a.Kubernetes.KubeConfigPath,
		KubeIncludeNamespaces:      a.Kubernetes.KubeIncludeNamespaces,
		KubeExcludeNamespaces:      a.Kubernetes.KubeExcludeNamespaces,
		KubeIncludeNamespaceLabels: a.Kubernetes.KubeIncludeNamespaceLabels,
		KubeExcludeNamespaceLabels: a.Kubernetes.KubeExcludeNamespaceLabels,
		KubeIncludeLabels:          a.Kubernetes.KubeIncludeLabels,
		KubeExcludeLabels:          a.Kubernetes.KubeExcludeLabels,
	}
}
