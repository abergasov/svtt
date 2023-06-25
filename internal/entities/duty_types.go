package entities

type DutyType string

const (
	DutyTypeProposer      DutyType = "PROPOSER"
	DutyTypeAttester      DutyType = "ATTESTER"
	DutyTypeAggregator    DutyType = "AGGREGATOR"
	DutyTypeSyncCommittee DutyType = "SYNC_COMMITTEE"
)

func (d DutyType) String() string {
	return string(d)
}
