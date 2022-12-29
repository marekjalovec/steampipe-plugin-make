package make

const (
	EndpointOrganization = "organizations"
	EndpointTeam         = "teams"
)

func ColumnsOrganization() []string {
	// skipped: "teams"
	return []string{"id", "name", "countryId", "timezoneId", "license", "zone", "serviceName", "isPaused", "externalId"}
}

func ColumnsTeam() []string {
	// skipped: "activeScenarios", "activeApps", "operations", "transfer"
	return []string{"id", "name", "organizationId"}
}
