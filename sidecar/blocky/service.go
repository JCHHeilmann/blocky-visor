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
	// Get full status output
	out, _ := exec.Command("systemctl", "status", serviceName).CombinedOutput()
	fullStatus := string(out)

	info := &ServiceInfo{FullStatus: fullStatus}

	// Parse active state
	activeOut, err := exec.Command("systemctl", "show", serviceName, "--property=ActiveState").Output()
	if err == nil {
		info.Active = strings.TrimPrefix(strings.TrimSpace(string(activeOut)), "ActiveState=")
	}

	// Parse sub state
	subOut, err := exec.Command("systemctl", "show", serviceName, "--property=SubState").Output()
	if err == nil {
		info.SubState = strings.TrimPrefix(strings.TrimSpace(string(subOut)), "SubState=")
	}

	// Parse main PID
	pidOut, err := exec.Command("systemctl", "show", serviceName, "--property=MainPID").Output()
	if err == nil {
		pid := strings.TrimPrefix(strings.TrimSpace(string(pidOut)), "MainPID=")
		if pid != "0" {
			info.PID = pid
		}
	}

	// Parse memory
	memOut, err := exec.Command("systemctl", "show", serviceName, "--property=MemoryCurrent").Output()
	if err == nil {
		mem := strings.TrimPrefix(strings.TrimSpace(string(memOut)), "MemoryCurrent=")
		if mem != "" && mem != "[not set]" {
			info.Memory = formatBytes(mem)
		}
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

func formatBytes(s string) string {
	// MemoryCurrent is in bytes as a string
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
