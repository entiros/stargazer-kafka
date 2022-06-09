package starlify

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Number        int `json:"number"`
}

type Service struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ServiceRequest struct {
	Name string `json:"name,omitempty"`
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
	Error   string   `json:"error"`
	Details *Details `json:"details"`
}

type PartitionDetails struct {
	ID int32 `json:"id"`
}

type TopicDetails struct {
	Name       string             `json:"name"`
	Partitions []PartitionDetails `json:"partitions"`
}

type Details struct {
	Topics []TopicDetails `json:"topics"`
}
