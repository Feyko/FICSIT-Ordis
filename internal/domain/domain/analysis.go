package domain

type AnalysisResult struct {
	SMLVersion        *string            `json:"smlVersion"`
	GameVersion       *string            `json:"gameVersion"`
	CommandLine       *string            `json:"commandLine"`
	Path              *string            `json:"path"`
	ModList           []string           `json:"modList"`
	LauncherID        *string            `json:"launcherID"`
	LauncherArtifact  *string            `json:"launcherArtifact"`
	PiracyInfo        *PiracyInformation `json:"piracyInfo"`
	DesiredSMLVersion *string            `json:"desiredSMLVersion"`
	CrashMatches      []CrashMatch       `json:"crashMatches"`
}

type PiracyInformation struct {
	IsPirated bool   `json:"isPirated"`
	Reason    string `json:"reason"`
}
