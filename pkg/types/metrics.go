package types

import "time"

// SystemMetrics holds all system metrics
type SystemMetrics struct {
	Disk    *DiskMetrics   `json:"disk"`
	Memory  *MemoryMetrics `json:"memory"`
	CPU     *CPUMetrics    `json:"cpu"`
	GPUs    []*GPUMetrics  `json:"gpus"`
	Updated time.Time      `json:"updated"`
}

// DiskMetrics holds disk usage information
type DiskMetrics struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Error       string  `json:"error,omitempty"`
}

// MemoryMetrics holds memory usage information
type MemoryMetrics struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Error       string  `json:"error,omitempty"`
}

// CPUMetrics holds CPU usage and temperature information
type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
	Temperature  float64 `json:"temperature"`
	Error        string  `json:"error,omitempty"`
}

// GPUMetrics holds GPU usage and temperature information
type GPUMetrics struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"` // "nvidia", "amd", "integrated"
	UsagePercent float64 `json:"usage_percent"`
	Temperature  float64 `json:"temperature"`
	Error        string  `json:"error,omitempty"`
}
