package models

type StartSimualtionReq struct {
	InitEpoch int64 `json:"init_epoch"`
	EndEpoch  int64 `json:"end_epoch"`
	Seed      int64 `json:"seed"`
}
