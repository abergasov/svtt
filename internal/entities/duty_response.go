package entities

type DutyResponce struct {
	Validator int      `json:"validator"`
	Duty      DutyType `json:"duty"`
	Height    int      `json:"height"`
	Response  string   `json:"response"`
}
