package domain

type LatestInformation struct {
	Revision int
	Text     string
}

type LatestInformationUpdate struct {
	Revision *int
	Text     *string
}
