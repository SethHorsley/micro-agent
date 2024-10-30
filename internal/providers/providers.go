package providers

func GetAvailableModels() map[string][]string {
	return map[string][]string{
		"open_ai": {
			"gpt-4-turbo-preview",
			"gpt-4",
			"gpt-3.5-turbo",
		},
		"anthropic": {
			"claude-3-opus-20240229",
			"claude-3-sonnet-20240229",
			"claude-3-haiku-20240307",
		},
	}
}
