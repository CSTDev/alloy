package x509exporter

import (
	"time"

	"github.com/go-kit/log"

	"github.com/enix/x509-certificate-exporter/v3/pkg/exporter"
	"github.com/grafana/alloy/internal/static/integrations"
	integrations_v2 "github.com/grafana/alloy/internal/static/integrations/v2"
	"github.com/grafana/alloy/internal/static/integrations/v2/metricsutils"
)

type Config struct {
	Directories                []string                  `yaml:"directories,omitempty"`
	Files                      []string                  `yaml:"files,omitempty"`
	YAMLs                      []string                  `yaml:"yamls,omitempty"`
	YAMLPaths                  []exporter.YAMLCertRef    `yaml:"yaml_paths,omitempty"`
	TrimPathComponents         int                       `yaml:"trim_path_components,omitempty"`
	MaxCacheDuration           time.Duration             `yaml:"max_cache_duration,omitempty"`
	ExposeRelativeMetrics      bool                      `yaml:"expose_relative_metrics,omitempty"`
	ExposeErrorMetrics         bool                      `yaml:"expose_error_metrics,omitempty"`
	ExposeLabels               []string                  `yaml:"expose_labels,omitempty"`
	ConfigMapKeys              []string                  `yaml:"config_map_keyss,omitempty"`
	KubeEnabled                bool                      `yaml:"kube_enabled,omitempty"`
	KubeConfigPath             string                    `yaml:"kube_config_path,omitempty"`
	KubeSecretTypes            []exporter.KubeSecretType `yaml:"kube_secret_types,omitempty"`
	KubeIncludeNamespaces      []string                  `yaml:"kube_include_namespaces,omitempty"`
	KubeExcludeNamespaces      []string                  `yaml:"kube_exclude_namespaces,omitempty"`
	KubeIncludeNamespaceLabels []string                  `yaml:"kube_include_namespace_labels,omitempty"`
	KubeExcludeNamespaceLabels []string                  `yaml:"kube_exclude_namespace_labels,omitempty"`
	KubeIncludeLabels          []string                  `yaml:"kube_include_labels,omitempty"`
	KubeExcludeLabels          []string                  `yaml:"kube_exclude_labels,omitempty"`
}

var DefaultConfig = Config{
	ExposeRelativeMetrics: false,
	ExposeErrorMetrics:    false,
}

// Name returns the name of the integration that this config represents.
func (c *Config) Name() string {
	return "x509_exporter"
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig
	type plain Config
	return unmarshal((*plain)(c))
}

// InstanceKey returns the hostname:port of the agent process.
func (c *Config) InstanceKey(agentKey string) (string, error) {
	return agentKey, nil
}

// NewIntegration converts this config into an instance of an integration.
func (c *Config) NewIntegration(logger log.Logger) (integrations.Integration, error) {
	return New(logger, c)
}

func init() {
	integrations.RegisterIntegration(&Config{})
	integrations_v2.RegisterLegacy(&Config{}, integrations_v2.TypeSingleton, metricsutils.Shim)
}

// New creates a new x509_exporter integration.
func New(log log.Logger, c *Config) (integrations.Integration, error) {
	// This is where you'd create the actual exporter instance
	options := exporter.Options{
		Directories:                c.Directories,
		Files:                      c.Files,
		YAMLs:                      c.YAMLs,
		TrimPathComponents:         c.TrimPathComponents,
		MaxCacheDuration:           c.MaxCacheDuration,
		ExposeRelativeMetrics:      c.ExposeRelativeMetrics,
		ExposeErrorMetrics:         c.ExposeErrorMetrics,
		ExposeLabels:               c.ExposeLabels,
		ConfigMapKeys:              c.ConfigMapKeys,
		KubeEnabled:                c.KubeEnabled,
		KubeConfigPath:             c.KubeConfigPath,
		KubeIncludeNamespaces:      c.KubeIncludeNamespaces,
		KubeExcludeNamespaces:      c.KubeExcludeNamespaces,
		KubeIncludeNamespaceLabels: c.KubeIncludeNamespaceLabels,
		KubeExcludeNamespaceLabels: c.KubeExcludeNamespaceLabels,
		KubeIncludeLabels:          c.KubeIncludeLabels,
		KubeExcludeLabels:          c.KubeExcludeLabels,
	}

	newExporter, _ := exporter.New(options)
	collector := &exporter.Collector{Exporter: &newExporter}

	return integrations.NewCollectorIntegration(
		c.Name(),
		integrations.WithCollectors(collector),
	), nil
}
