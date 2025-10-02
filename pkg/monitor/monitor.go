package monitor

import (
	"context"
	"time"

	"p-monitor/internal/logs"
	"p-monitor/pkg/config"
	"p-monitor/pkg/gpu"
	"p-monitor/pkg/types"
)

// Use types from the types package
type SystemMetrics = types.SystemMetrics
type DiskMetrics = types.DiskMetrics
type MemoryMetrics = types.MemoryMetrics
type CPUMetrics = types.CPUMetrics
type GPUMetrics = types.GPUMetrics

// Monitor handles system monitoring
type Monitor struct {
	config  *config.Config
	metrics chan *SystemMetrics
	ctx     context.Context
	cancel  context.CancelFunc
	latest  *SystemMetrics
}

// New creates a new monitor instance
func New(cfg *config.Config) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Monitor{
		config:  cfg,
		metrics: make(chan *SystemMetrics, 1),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start starts the monitoring loop
func (m *Monitor) Start() {
	logs.Info("Starting system monitor")

	ticker := time.NewTicker(time.Duration(m.config.GetUpdateIntervalSeconds()) * time.Second)
	defer ticker.Stop()

	// Initial collection
	m.collectMetrics()

	for {
		select {
		case <-ticker.C:
			m.collectMetrics()
		case <-m.ctx.Done():
			logs.Info("Stopping system monitor")
			return
		}
	}
}

// Stop stops the monitoring loop
func (m *Monitor) Stop() {
	m.cancel()
}

// GetMetrics returns the latest metrics
func (m *Monitor) GetMetrics() *SystemMetrics {
	return m.latest
}

// collectMetrics collects all system metrics
func (m *Monitor) collectMetrics() {
	metrics := &SystemMetrics{
		Updated: time.Now(),
	}

	// Collect disk metrics
	metrics.Disk = m.collectDiskMetrics()

	// Collect memory metrics
	metrics.Memory = m.collectMemoryMetrics()

	// Collect CPU metrics
	metrics.CPU = m.collectCPUMetrics()

	// Collect GPU metrics
	metrics.GPUs = m.collectGPUMetrics()

	// Update latest metrics
	m.latest = metrics

	// Send metrics (non-blocking)
	select {
	case m.metrics <- metrics:
	default:
		// Channel is full, skip this update
	}
}

// collectDiskMetrics collects disk usage metrics
func (m *Monitor) collectDiskMetrics() *DiskMetrics {
	disk := &DiskMetrics{}

	// Get root filesystem usage
	usage, err := getDiskUsage("/")
	if err != nil {
		disk.Error = err.Error()
		logs.Error("Failed to get disk usage: %v", err)
		return disk
	}

	disk.Total = usage.Total
	disk.Used = usage.Used
	disk.UsedPercent = usage.UsedPercent
	return disk
}

// collectMemoryMetrics collects memory usage metrics
func (m *Monitor) collectMemoryMetrics() *MemoryMetrics {
	memory := &MemoryMetrics{}

	usage, err := getMemoryUsage()
	if err != nil {
		memory.Error = err.Error()
		logs.Error("Failed to get memory usage: %v", err)
		return memory
	}

	memory.Total = usage.Total
	memory.Used = usage.Used
	memory.UsedPercent = usage.UsedPercent
	return memory
}

// collectCPUMetrics collects CPU usage and temperature metrics
func (m *Monitor) collectCPUMetrics() *CPUMetrics {
	cpu := &CPUMetrics{}

	// Get CPU usage
	usage, err := getCPUUsage()
	if err != nil {
		cpu.Error = err.Error()
		logs.Error("Failed to get CPU usage: %v", err)
		return cpu
	}
	cpu.UsagePercent = usage

	// Get CPU temperature
	temp, err := getCPUTemperature()
	if err != nil {
		logs.Error("Failed to get CPU temperature: %v", err)
		// Don't set error for temperature as it's optional
	} else {
		cpu.Temperature = temp
	}

	return cpu
}

// collectGPUMetrics collects GPU usage and temperature metrics
func (m *Monitor) collectGPUMetrics() []*GPUMetrics {
	return gpu.CollectAllGPUMetrics()
}
