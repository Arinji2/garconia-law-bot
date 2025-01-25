package commands_utils

func FormatDescription(description string) string {
	if len(description) > 90 {
		return description[:90] + "..."
	}
	return description
}
