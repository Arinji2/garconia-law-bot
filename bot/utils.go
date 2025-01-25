package bot

func formatDescription(description string) string {
	if len(description) > 90 {
		return description[:90] + "..."
	}
	return description
}
