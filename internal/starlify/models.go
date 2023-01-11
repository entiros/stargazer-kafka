package starlify

import "time"

type EndpointResponse struct {
	Type           string    `json:"type"`
	Id             string    `json:"id"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	GoLiveDateTime time.Time `json:"goLiveDateTime"`
	Name           string    `json:"name"`
	CreatedByAgent bool      `json:"createdByAgent"`
	Provider       struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
		Name    string    `json:"name"`
		Links   []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"provider"`
	Attributes           []interface{} `json:"attributes"`
	Engagements          []interface{} `json:"engagements"`
	EndpointInteractions []interface{} `json:"endpointInteractions"`
	Domain               struct {
		Type  string `json:"type"`
		Id    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"domain"`
	Network struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
		Name    string    `json:"name"`
		Links   []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"network"`
	Links []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

type EndpointRequest struct {
	Name string `json:"name"`
}

type Middleware struct {
	Type           string    `json:"type"`
	Id             string    `json:"id"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	GoLiveDateTime time.Time `json:"goLiveDateTime"`
	Name           string    `json:"name"`
	KafkaPrefix    string    `json:"kafkaPrefix"`
	Endpoints      []struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
		Name    string    `json:"name"`
		Links   []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"endpoints"`
	Attributes   []interface{} `json:"attributes"`
	Engagements  []interface{} `json:"engagements"`
	Observations []interface{} `json:"observations"`
	Domain       struct {
		Type  string `json:"type"`
		Id    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"domain"`
	Network struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
		Name    string    `json:"name"`
		Links   []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"network"`
	Contents []interface{} `json:"contents"`
	Links    []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

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
	Type    string    `json:"type"`
	Id      string    `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Domain  struct {
		Type  string `json:"type"`
		Id    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"domain"`
	MiddlewareId string `json:"middlewareId"`
	Person       struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
		Name    string    `json:"name"`
		Agent   bool      `json:"agent"`
		Links   []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"person"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AgentType   string `json:"agentType"`
	Links       []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
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
