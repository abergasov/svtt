package entities

type DutyRequest struct {
	Validator int      `json:"validator"`
	Duty      DutyType `json:"duty"`
	Height    int      `json:"height"`
}
