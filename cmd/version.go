package cmd

var version string

func Version() string {
	if version == "" {
		return "development"
	}
	return version
}
