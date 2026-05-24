package utils

import "fmt"

func FormatBytes(bytes int64) string {
	if bytes < 0 {
		bytes = 0
	}

	units := []string{"B", "KB", "MB", "GB", "TB"}
	value := float64(bytes)
	unit := units[0]
	for i := 1; i < len(units) && value >= 1024; i++ {
		value /= 1024
		unit = units[i]
	}

	if unit == "B" {
		return fmt.Sprintf("%d B", bytes)
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

func FormatSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 0 {
		bytesPerSecond = 0
	}
	return FormatBytes(int64(bytesPerSecond)) + "/s"
}
