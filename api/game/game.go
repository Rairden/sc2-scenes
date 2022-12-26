package game

type Game struct {
	IsReplay    bool    `json:"isReplay"`
	DisplayTime float32 `json:"displayTime"`
	Players     []any   `json:"players"`
}

type Players struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Race   string `json:"race"`
	Result string `json:"result"`
}
