package blocky

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ErrSystemdNotAvailable is returned when systemctl is not found on the system.
var ErrSystemdNotAvailable = errors.New("systemd is not available: systemctl not found in PATH")

// SystemdAvailable checks whether systemctl exists on the system.
func SystemdAvailable() bool {
	_, err := exec.LookPath("systemctl")
	return err == nil
}

// ServiceInfo holds parsed systemd service status.
type ServiceInfo struct {
	Active     string `json:"active"`
	SubState   string `json:"sub_state"`
	PID        string `json:"pid,omitempty"`
	Memory     string `json:"memory,omitempty"`
	Uptime     string `json:"uptime,omitempty"`
	FullStatus string `json:"full_status"`
}

// Status queries systemctl for the service status.
func Status(serviceName string) (*ServiceInfo, error) {
	if !SystemdAvailable() {
		return nil, ErrSystemdNotAvailable
	}

	out, _ := exec.Command("systemctl", "status", serviceName).CombinedOutput()
	info := &ServiceInfo{FullStatus: string(out)}

	propsOut, err := exec.Command("systemctl", "show", serviceName,
		"--property=ActiveState,SubState,MainPID,MemoryCurrent").Output()
	if err != nil {
		return info, nil
	}

	props := parseProperties(string(propsOut))
	info.Active = props["ActiveState"]
	info.SubState = props["SubState"]
	if pid := props["MainPID"]; pid != "0" {
		info.PID = pid
	}
	if mem := props["MemoryCurrent"]; mem != "" && mem != "[not set]" {
		info.Memory = formatBytes(mem)
	}

	return info, nil
}

// Restart restarts the systemd service.
func Restart(serviceName string) error {
	if !SystemdAvailable() {
		return ErrSystemdNotAvailable
	}
	out, err := exec.Command("systemctl", "restart", serviceName).CombinedOutput()
	if err != nil {
		return fmt.Errorf("restart failed: %s: %w", string(out), err)
	}
	return nil
}

func parseProperties(output string) map[string]string {
	props := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			props[key] = value
		}
	}
	return props
}

func formatBytes(s string) string {
	var bytes uint64
	if _, err := fmt.Sscanf(s, "%d", &bytes); err != nil {
		return s
	}
	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
