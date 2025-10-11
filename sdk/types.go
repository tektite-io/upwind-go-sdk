// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

// Common types

// Severity levels for findings, detections, and events
const (
	SeverityLow      = "LOW"
	SeverityMedium   = "MEDIUM"
	SeverityHigh     = "HIGH"
	SeverityCritical = "CRITICAL"
)

// Cloud providers
const (
	CloudProviderAWS   = "AWS"
	CloudProviderGCP   = "GCP"
	CloudProviderAzure = "AZURE"
	CloudProviderBYOC  = "BYOC"
)

// Status constants
const (
	StatusOpen     = "OPEN"
	StatusPending  = "PENDING"
	StatusArchived = "ARCHIVED"
	StatusPass     = "PASS"
	StatusFail     = "FAIL"
	StatusEnabled  = "ENABLED"
	StatusDisabled = "DISABLED"
)

// Tag represents a key-value tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Resource represents a cloud resource
type Resource struct {
	ID               string            `json:"id"`
	ExternalID       string            `json:"external_id,omitempty"`
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	Path             string            `json:"path,omitempty"`
	CloudProvider    string            `json:"cloud_provider"`
	CloudAccountID   string            `json:"cloud_account_id"`
	CloudAccountName string            `json:"cloud_account_name,omitempty"`
	CloudAccountTags []Tag             `json:"cloud_account_tags,omitempty"`
	Region           string            `json:"region,omitempty"`
	ClusterID        string            `json:"cluster_id,omitempty"`
	Namespace        string            `json:"namespace,omitempty"`
	InternetExposure *InternetExposure `json:"internet_exposure,omitempty"`
	RiskCategories   []string          `json:"risk_categories,omitempty"`
}

// InternetExposure represents internet exposure information
type InternetExposure struct {
	Ingress *InternetExposureDetails `json:"ingress,omitempty"`
}

// InternetExposureDetails contains details about internet exposure
type InternetExposureDetails struct {
	ActiveCommunication bool `json:"active_communication"`
}

// Image represents container image information
type Image struct {
	Name       string `json:"name"`
	Digest     string `json:"digest"`
	URI        string `json:"uri"`
	Registry   string `json:"registry,omitempty"`
	Repository string `json:"repository,omitempty"`
	OSVersion  string `json:"os_version,omitempty"`
	OSName     string `json:"os_name,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

// Package represents a software package
type Package struct {
	Name      string `json:"name"`
	Framework string `json:"framework,omitempty"`
	Type      string `json:"type,omitempty"`
	Version   string `json:"version"`
	InUse     bool   `json:"in_use"`
}

// Vulnerability Finding Types

// VulnerabilityFinding represents a vulnerability finding
type VulnerabilityFinding struct {
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	Source        string         `json:"source"`
	FirstSeenTime string         `json:"first_seen_time"`
	LastScanTime  string         `json:"last_scan_time"`
	Vulnerability *Vulnerability `json:"vulnerability,omitempty"`
	Image         *Image         `json:"image,omitempty"`
	Package       *Package       `json:"package,omitempty"`
	Resource      *Resource      `json:"resource,omitempty"`
	Remediation   []Remediation  `json:"remediation,omitempty"`
}

// Vulnerability represents vulnerability details
type Vulnerability struct {
	Name              string         `json:"name,omitempty"`
	Description       string         `json:"description,omitempty"`
	Exploitable       bool           `json:"exploitable"`
	NVDCVEID          string         `json:"nvd_cve_id,omitempty"`
	NVDDescription    string         `json:"nvd_description,omitempty"`
	NVDPublishTime    string         `json:"nvd_publish_time,omitempty"`
	CVEFirstSeenTime  string         `json:"cve_first_seen_time,omitempty"`
	NVDCVSSV2Severity string         `json:"nvd_cvss_v2_severity,omitempty"`
	NVDCVSSV2Score    string         `json:"nvd_cvss_v2_score,omitempty"`
	NVDCVSSV3Severity string         `json:"nvd_cvss_v3_severity,omitempty"`
	NVDCVSSV3Score    string         `json:"nvd_cvss_v3_score,omitempty"`
	NVDCVSSV4Severity string         `json:"nvd_cvss_v4_severity,omitempty"`
	NVDCVSSV4Score    string         `json:"nvd_cvss_v4_score,omitempty"`
	ImpactMetrics     *ImpactMetrics `json:"impact_metrics,omitempty"`
}

// ImpactMetrics represents the impact metrics for a vulnerability
type ImpactMetrics struct {
	AffectedResourceCount int `json:"affected_resource_count"`
	AffectedImageCount    int `json:"affected_image_count"`
}

// Remediation represents remediation information
type Remediation struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Configuration Finding Types

// ConfigurationFinding represents a configuration finding
type ConfigurationFinding struct {
	ID            string                  `json:"id"`
	Status        string                  `json:"status"`
	Severity      string                  `json:"severity"`
	Title         string                  `json:"title"`
	Description   string                  `json:"description,omitempty"`
	FirstSeenTime string                  `json:"first_seen_time"`
	LastSeenTime  string                  `json:"last_seen_time"`
	LastSyncTime  string                  `json:"last_sync_time,omitempty"`
	Framework     *ConfigurationFramework `json:"framework,omitempty"`
	Check         *ConfigurationCheck     `json:"check,omitempty"`
	Resource      *Resource               `json:"resource,omitempty"`
}

// ConfigurationFramework represents a compliance framework
type ConfigurationFramework struct {
	ID               string                        `json:"id"`
	Status           string                        `json:"status,omitempty"`
	Version          string                        `json:"version,omitempty"`
	Revision         string                        `json:"revision,omitempty"`
	Title            string                        `json:"title"`
	Description      string                        `json:"description,omitempty"`
	CloudProvider    string                        `json:"cloud_provider,omitempty"`
	CreateTime       string                        `json:"create_time,omitempty"`
	UpdateTime       string                        `json:"update_time,omitempty"`
	LastScanTime     string                        `json:"last_scan_time,omitempty"`
	Type             string                        `json:"type,omitempty"`
	ComplianceStatus *ConfigurationFrameworkStatus `json:"compliance_status,omitempty"`
	RolloutState     string                        `json:"rollout_state,omitempty"`
}

// ConfigurationFrameworkStatus represents framework compliance status
type ConfigurationFrameworkStatus struct {
	Score int `json:"score"`
}

// ConfigurationCheck represents a configuration check
type ConfigurationCheck struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Remediation string `json:"remediation,omitempty"`
}

// ConfigurationRule represents a configuration rule
type ConfigurationRule struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Framework     string `json:"framework,omitempty"`
	FindingsCount int    `json:"findings_count"`
	CreateTime    string `json:"create_time,omitempty"`
	UpdateTime    string `json:"update_time,omitempty"`
}

// Threat Detection Types

// ThreatDetection represents a threat detection
type ThreatDetection struct {
	ID              string                   `json:"id"`
	Type            string                   `json:"type"`
	Category        string                   `json:"category"`
	Severity        string                   `json:"severity"`
	Status          string                   `json:"status"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description,omitempty"`
	FirstSeenTime   string                   `json:"first_seen_time"`
	LastSeenTime    string                   `json:"last_seen_time"`
	OccurrenceCount int                      `json:"occurrence_count"`
	Resource        *Resource                `json:"resource,omitempty"`
	MitreAttacks    []MitreAttackDetails     `json:"mitre_attacks,omitempty"`
	Triggers        []ThreatDetectionTrigger `json:"triggers,omitempty"`
}

// ThreatDetectionTrigger represents a policy trigger
type ThreatDetectionTrigger struct {
	PolicyID   string                 `json:"policy_id"`
	PolicyName string                 `json:"policy_name"`
	Events     []ThreatDetectionEvent `json:"events,omitempty"`
}

// ThreatDetectionEvent represents an event in a detection
type ThreatDetectionEvent struct {
	ID          string                 `json:"id"`
	EventType   string                 `json:"event_type"`
	Description string                 `json:"description,omitempty"`
	EventTime   string                 `json:"event_time"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// MitreAttackDetails represents MITRE ATT&CK framework information
type MitreAttackDetails struct {
	TacticID      string `json:"tactic_id"`
	TacticName    string `json:"tactic_name"`
	TechniqueID   string `json:"technique_id"`
	TechniqueName string `json:"technique_name"`
}

// ThreatEvent represents a threat event
type ThreatEvent struct {
	ID            string    `json:"id"`
	Type          string    `json:"type"`
	Severity      string    `json:"severity"`
	Category      string    `json:"category"`
	Status        string    `json:"status"`
	Title         string    `json:"title"`
	FirstSeenTime string    `json:"first_seen_time"`
	LastSeenTime  string    `json:"last_seen_time"`
	Resource      *Resource `json:"resource,omitempty"`
}

// ThreatPolicy represents a threat policy
type ThreatPolicy struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Scope       string `json:"scope"`
	OpenIssues  int    `json:"open_issues"`
	ManagedBy   string `json:"managed_by"`
	Enabled     bool   `json:"enabled"`
}

// Workflow Types

// Workflow represents a workflow
type Workflow struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Type              string          `json:"type"`
	Status            string          `json:"status"`
	LastExecutionTime string          `json:"last_execution_time,omitempty"`
	Config            *WorkflowConfig `json:"config,omitempty"`
}

// WorkflowConfig represents workflow configuration
type WorkflowConfig struct {
	Selectors []WorkflowSelector `json:"selectors,omitempty"`
	Actions   []WorkflowAction   `json:"actions,omitempty"`
	Trigger   *WorkflowTrigger   `json:"trigger,omitempty"`
}

// WorkflowTrigger represents workflow trigger configuration
type WorkflowTrigger struct {
	Type       string   `json:"type"`
	Severities []string `json:"severities,omitempty"`
	Categories []string `json:"categories,omitempty"`
}

// WorkflowSelector represents a workflow selector (interface for different types)
type WorkflowSelector map[string]interface{}

// WorkflowAction represents a workflow action (interface for different types)
type WorkflowAction map[string]interface{}

// Integration Webhook Types

// IntegrationWebhook represents an integration webhook
type IntegrationWebhook struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Vendor     string                 `json:"vendor"`
	Status     string                 `json:"status"`
	Config     map[string]interface{} `json:"config,omitempty"`
	CreateTime string                 `json:"create_time,omitempty"`
	UpdateTime string                 `json:"update_time,omitempty"`
}

// API Security Types

// ApiEndpoint represents an API endpoint
type ApiEndpoint struct {
	ID            string                   `json:"id"`
	Method        string                   `json:"method"`
	URI           string                   `json:"uri"`
	ResourceID    string                   `json:"resource_id"`
	FirstSeenTime string                   `json:"first_seen_time"`
	LastSeenTime  string                   `json:"last_seen_time"`
	Domains       []string                 `json:"domains,omitempty"`
	StatusCodes   []string                 `json:"status_codes,omitempty"`
	RiskOverview  *ApiEndpointRiskOverview `json:"risk_overview,omitempty"`
}

// ApiEndpointRiskOverview represents risk overview for an API endpoint
type ApiEndpointRiskOverview struct {
	Authentication        *ApiEndpointAuthentication        `json:"authentication,omitempty"`
	InternetExposure      *ApiEndpointInternetExposure      `json:"internet_exposure,omitempty"`
	SensitiveDataFindings []ApiEndpointSensitiveDataFinding `json:"sensitive_data_findings,omitempty"`
}

// ApiEndpointAuthentication represents authentication state
type ApiEndpointAuthentication struct {
	State string `json:"state"`
}

// ApiEndpointInternetExposure represents internet exposure for API endpoint
type ApiEndpointInternetExposure struct {
	Ingress *ApiEndpointExposureDetails `json:"ingress,omitempty"`
}

// ApiEndpointExposureDetails represents exposure details
type ApiEndpointExposureDetails struct {
	LastSeenTime string `json:"last_seen_time"`
}

// ApiEndpointSensitiveDataFinding represents sensitive data finding
type ApiEndpointSensitiveDataFinding struct {
	Type         string `json:"type"`
	Category     string `json:"category"`
	LastSeenTime string `json:"last_seen_time"`
}

// SBOM Package Types

// SbomPackage represents an SBOM package
type SbomPackage struct {
	Name                   string                  `json:"name"`
	Version                string                  `json:"version"`
	PackageManager         string                  `json:"package_manager,omitempty"`
	Framework              string                  `json:"framework,omitempty"`
	Licenses               []string                `json:"licenses,omitempty"`
	VulnerabilitiesSummary *VulnerabilitiesSummary `json:"vulnerabilities_summary,omitempty"`
	ResourcesSummary       *ResourcesSummary       `json:"resources_summary,omitempty"`
	ImagesSummary          *ImagesSummary          `json:"images_summary,omitempty"`
}

// VulnerabilitiesSummary represents a summary of vulnerabilities
type VulnerabilitiesSummary struct {
	CriticalCount     int `json:"critical_count"`
	HighCount         int `json:"high_count"`
	MediumCount       int `json:"medium_count"`
	LowCount          int `json:"low_count"`
	UnclassifiedCount int `json:"unclassified_count"`
	TotalCount        int `json:"total_count"`
}

// ResourcesSummary represents a summary of resources
type ResourcesSummary struct {
	InUseCount int `json:"in_use_count"`
	TotalCount int `json:"total_count"`
}

// ImagesSummary represents a summary of images
type ImagesSummary struct {
	AffectedCount int `json:"affected_count"`
}

// Cloud Account Types

// CloudAccount represents a cloud account
type CloudAccount struct {
	ID        string                 `json:"id"`
	AccountID string                 `json:"account_id"`
	Name      string                 `json:"name"`
	Provider  string                 `json:"provider"`
	Config    map[string]interface{} `json:"config,omitempty"`
}

// Query parameter types

// VulnerabilityFindingsQuery represents query parameters for vulnerability findings
type VulnerabilityFindingsQuery struct {
	PageToken                  string
	PerPage                    int
	CloudAccountID             string
	ClusterID                  string
	Namespace                  string
	IngressActiveCommunication *bool
	InternetExposure           *bool
	InUse                      *bool
	Exploitable                *bool
	FixAvailable               *bool
	Severity                   string
	ImageName                  string
	Framework                  string
}

// ConfigurationFindingsQuery represents query parameters for configuration findings
type ConfigurationFindingsQuery struct {
	MinLastSeenTime         string
	MaxLastSeenTime         string
	Status                  string
	Severity                string
	ResourceName            string
	CheckTitle              string
	CheckID                 string
	FrameworkID             string
	FrameworkTitle          string
	CloudAccountTags        []string
	IncludeCloudAccountTags bool
}

// ThreatDetectionsQuery represents query parameters for threat detections
type ThreatDetectionsQuery struct {
	Severity         string
	Type             string
	Category         string
	MinFirstSeenTime string
	MaxFirstSeenTime string
	MinLastSeenTime  string
	MaxLastSeenTime  string
}

// ThreatEventsQuery represents query parameters for threat events
type ThreatEventsQuery struct {
	CloudAccountID   string
	Severity         string
	Category         string
	MinFirstSeenTime string
	MaxFirstSeenTime string
	MinLastSeenTime  string
	MaxLastSeenTime  string
	Page             int
	PerPage          int
}

// ApiEndpointsQuery represents query parameters for API endpoints
type ApiEndpointsQuery struct {
	PerPage                 int
	PageToken               string
	Method                  string
	AuthenticationState     string
	HasInternetIngress      *bool
	HasVulnerability        *bool
	HasSensitiveData        *bool
	CloudAccountID          string
	CloudProvider           string
	ResourceType            string
	CloudOrganizationID     string
	CloudOrganizationUnitID string
	Domain                  string
	ClusterID               string
	Namespace               string
}
