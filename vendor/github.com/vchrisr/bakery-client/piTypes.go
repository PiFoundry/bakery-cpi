package bakeryclient

type PiInfo struct {
	Id             string   `json:"id"`
	Status         piStatus `json:"status"`
	Disks          []Disk   `json:"disks,omitempty"`
	SourceBakeform Bakeform `json:"sourceBakeform,omitempty"`
}

type piStatus int

const (
	NOTINUSE  piStatus = 1
	INUSE     piStatus = 2
	PREPARING piStatus = 3
)

type Bakeform struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Disk struct {
	ID       string `json:"id"`
	Location string `json:"location"`
}
