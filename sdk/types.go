package sdk

import "time"

type VulnerabilityFinding struct {
	ID            string        `json:"id"`
	Vulnerability Vulnerability `json:"vulnerability"`
	Image         ImageInfo     `json:"image"`
	Resource      ResourceInfo  `json:"resource"`
	Remediation   []Remediation `json:"remediation"`
	Source        string        `json:"source"`
	FirstSeenTime time.Time     `json:"first_seen_time"`
	LastScanTime  time.Time     `json:"last_scan_time"`
	Package       PackageInfo   `json:"package"`
}

type Vulnerability struct {
	Exploitable       bool      `json:"exploitable"`
	NVDDescription    string    `json:"nvd_description"`
	NVDCVEID          string    `json:"nvd_cve_id"`
	NVDPublishTime    time.Time `json:"nvd_publish_time"`
	NVDCVSSV3Severity string    `json:"nvd_cvss_v3_severity"`
	NVDCVSSV3Score    string    `json:"nvd_cvss_v3_score"`
}

type ImageInfo struct {
	Name      string `json:"name"`
	Digest    string `json:"digest"`
	URI       string `json:"uri"`
	Tag       string `json:"tag"`
	OSVersion string `json:"os_version"`
	OSName    string `json:"os_name"`
}

type ResourceInfo struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	Region           string           `json:"region"`
	Namespace        string           `json:"namespace"`
	CloudProvider    string           `json:"cloud_provider"`
	ClusterID        string           `json:"cluster_id"`
	CloudAccountID   string           `json:"cloud_account_id"`
	InternetExposure InternetExposure `json:"internet_exposure"`
}

type InternetExposure struct {
	Ingress IngressExposure `json:"ingress"`
}

type IngressExposure struct {
	ActiveCommunication bool `json:"active_communication"`
	Exposed             bool `json:"exposed"`
}

type Remediation struct {
	Type string          `json:"type"`
	Data RemediationData `json:"data"`
}

type RemediationData struct {
	FixedInVersion string `json:"fixed_in_version"`
}

type PackageInfo struct {
	Name      string `json:"name"`
	Framework string `json:"framework"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	InUse     bool   `json:"in_use"`
}

type Finding struct {
	ID        string `json:"id"`
	Severity  string `json:"severity"`
	ImageName string `json:"image_name"`
	Component string `json:"component"`
	CVE       string `json:"cve"`
}

type FindingsQuery struct {
	OrgID                      string
	ImageName                  *string
	InUse                      *bool
	IngressActiveCommunication *bool
	Severities                 []string
	PerPage                    *int
}
