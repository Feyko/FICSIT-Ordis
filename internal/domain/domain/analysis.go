package domain

type AnalysisResult struct {
	SMLVersion        *string            `json:"smlVersion"`
	GameVersion       *string            `json:"gameVersion"`
	CommandLine       *string            `json:"commandLine"`
	Path              *string            `json:"path"`
	Cl                *string            `json:"cl"`
	ModList           []string           `json:"modList"`
	PiracyInfo        *PiracyInformation `json:"piracyInfo"`
	DesiredSMLVersion *string            `json:"desiredSMLVersion"`
	CrashMatches      []CrashMatch       `json:"crashMatches"`
}

type PiracyInformation struct {
	IsPirated bool   `json:"isPirated"`
	Reason    string `json:"reason"`
}
