package starlify

type Service struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ServiceRequest struct {
	Name string `json:"name,omitempty"`
}

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

type ServicesPage struct {
	Services []Service `json:"content"`
	Page     Page      `json:"page"`
}

type Agent struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	AgentType string `json:"agentType"`
	LastSeen  string `json:"lastSeen"`
}

type AgentRequest struct {
	Error   string  `json:"error"`
	Details Details `json:"details"`
}

type DetailsPartition struct {
	ID int32 `json:"id"`
}

type DetailsTopic struct {
	Name       string             `json:"name"`
	Partitions []DetailsPartition `json:"partitions"`
}

type Details struct {
	Topics []DetailsTopic `json:"topics"`
}
