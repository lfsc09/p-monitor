package gpu

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"p-monitor/internal/logs"
	"p-monitor/pkg/types"
)

// CollectAllGPUMetrics collects metrics from all available GPUs
func CollectAllGPUMetrics() []*types.GPUMetrics {
	var gpus []*types.GPUMetrics

	// Collect NVIDIA GPUs
	nvidiaGPUs := collectNVIDIAGPUs()
	gpus = append(gpus, nvidiaGPUs...)

	// Collect AMD GPUs
	amdGPUs := collectAMDGPUs()
	gpus = append(gpus, amdGPUs...)

	// Collect integrated GPU (if available)
	integratedGPU := collectIntegratedGPU()
	if integratedGPU != nil {
		gpus = append(gpus, integratedGPU)
	}

	return gpus
}

// collectNVIDIAGPUs collects metrics from NVIDIA GPUs using nvidia-smi
func collectNVIDIAGPUs() []*types.GPUMetrics {
	var gpus []*types.GPUMetrics

	// Check if nvidia-smi is available
	if !commandExists("nvidia-smi") {
		logs.Debug("nvidia-smi not found, skipping NVIDIA GPU monitoring")
		return gpus
	}

	// Run nvidia-smi to get GPU information
	cmd := exec.Command("nvidia-smi", "--query-gpu=index,name,utilization.gpu,temperature.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		logs.Error("Failed to run nvidia-smi: %v", err)
		return gpus
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		gpu := parseNVIDIALine(line)
		if gpu != nil {
			gpus = append(gpus, gpu)
		}
	}

	return gpus
}

// collectAMDGPUs collects metrics from AMD GPUs using radeontop
func collectAMDGPUs() []*types.GPUMetrics {
	var gpus []*types.GPUMetrics

	// Check if radeontop is available
	if !commandExists("radeontop") {
		logs.Debug("radeontop not found, skipping AMD GPU monitoring")
		return gpus
	}

	// Run radeontop to get GPU information
	cmd := exec.Command("radeontop", "-l", "1", "-d", "-")
	output, err := cmd.Output()
	if err != nil {
		logs.Error("Failed to run radeontop: %v", err)
		return gpus
	}

	gpu := parseRadeontopOutput(string(output))
	if gpu != nil {
		gpus = append(gpus, gpu)
	}

	return gpus
}

// collectIntegratedGPU collects metrics from integrated GPU
func collectIntegratedGPU() *types.GPUMetrics {
	// Try to get integrated GPU info from various sources
	// This is a simplified implementation - in practice, you might need to
	// check different sources depending on the system

	// For now, we'll try to get basic info from lspci or similar
	if commandExists("lspci") {
		cmd := exec.Command("lspci", "-v")
		output, err := cmd.Output()
		if err != nil {
			logs.Debug("Failed to run lspci: %v", err)
			return nil
		}

		// Look for integrated graphics
		if strings.Contains(string(output), "VGA") &&
			(strings.Contains(string(output), "Intel") || strings.Contains(string(output), "AMD")) {
			return &types.GPUMetrics{
				Name:         "Integrated GPU",
				Type:         "integrated",
				UsagePercent: 0, // Integrated GPU usage is harder to get
				Temperature:  0, // Temperature might not be available
				Error:        "Integrated GPU monitoring not fully implemented",
			}
		}
	}

	return nil
}

// parseNVIDIALine parses a line from nvidia-smi output
func parseNVIDIALine(line string) *types.GPUMetrics {
	parts := strings.Split(line, ", ")
	if len(parts) < 4 {
		logs.Error("Invalid nvidia-smi output line: %s", line)
		return nil
	}

	index := strings.TrimSpace(parts[0])
	name := strings.TrimSpace(parts[1])

	usageStr := strings.TrimSpace(parts[2])
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		logs.Error("Failed to parse GPU usage: %s", usageStr)
		usage = 0
	}

	tempStr := strings.TrimSpace(parts[3])
	temperature, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		logs.Error("Failed to parse GPU temperature: %s", tempStr)
		temperature = 0
	}

	return &types.GPUMetrics{
		Name:         fmt.Sprintf("NVIDIA GPU %s (%s)", index, name),
		Type:         "nvidia",
		UsagePercent: usage,
		Temperature:  temperature,
	}
}

// parseRadeontopOutput parses radeontop output
func parseRadeontopOutput(output string) *types.GPUMetrics {
	// radeontop output format is complex, this is a simplified parser
	// In practice, you might need to parse the output more carefully

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "gpu") && strings.Contains(line, "%") {
			// Extract usage percentage
			re := regexp.MustCompile(`gpu\s+(\d+\.?\d*)%`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				usage, err := strconv.ParseFloat(matches[1], 64)
				if err != nil {
					continue
				}

				return &types.GPUMetrics{
					Name:         "AMD GPU",
					Type:         "amd",
					UsagePercent: usage,
					Temperature:  0, // radeontop might not provide temperature
				}
			}
		}
	}

	return nil
}

// commandExists checks if a command exists in the system
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
