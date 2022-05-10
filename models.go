package stargazer_kafka

import "time"

type StarlifySystem struct {
	Type           string    `json:"type"`
	Id             string    `json:"id"`
	Created        time.Time `json:"created"`
	Updated        time.Time `json:"updated"`
	GoLiveDateTime time.Time `json:"goLiveDateTime"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Services       []struct {
		Type           string    `json:"type"`
		Id             string    `json:"id"`
		Created        time.Time `json:"created"`
		Updated        time.Time `json:"updated"`
		GoLiveDateTime time.Time `json:"goLiveDateTime"`
		Name           string    `json:"name"`
		Provider       struct {
			Type           string    `json:"type"`
			Id             string    `json:"id"`
			Created        time.Time `json:"created"`
			Updated        time.Time `json:"updated"`
			GoLiveDateTime time.Time `json:"goLiveDateTime"`
			Name           string    `json:"name"`
			Description    string    `json:"description"`
			Links          []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"provider"`
		CertifiedIntegratorCompliant struct {
			DONE      int `json:"DONE"`
			NOSTATUS  int `json:"NO_STATUS"`
			NOTDONE   int `json:"NOT_DONE"`
			NOTNEEDED int `json:"NOT_NEEDED"`
		} `json:"certifiedIntegratorCompliant"`
		Domain struct {
			Type    string    `json:"type"`
			Id      string    `json:"id"`
			Created time.Time `json:"created"`
			Name    string    `json:"name"`
			Links   []struct {
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
	} `json:"services"`
	Attributes []struct {
		Type           string        `json:"type"`
		Id             string        `json:"id"`
		Name           string        `json:"name"`
		Enabled        bool          `json:"enabled"`
		EntityType     string        `json:"entityType"`
		AttributeType  string        `json:"attributeType"`
		PossibleValues []string      `json:"possibleValues"`
		DefaultValue   []interface{} `json:"defaultValue"`
		Value          []interface{} `json:"value"`
		Links          []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"attributes"`
	Domain struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Name    string    `json:"name"`
		Links   []struct {
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

type Attributes struct {
	Type          string `json:"type"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Enabled       bool   `json:"enabled"`
	EntityType    string `json:"entityType"`
	AttributeType string `json:"attributeType"`
	DefaultValue  string `json:"defaultValue"`
	Value         string `json:"value"`
	Links         []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}

type StarlifyService struct {
	Type     string    `json:"type"`
	Id       string    `json:"id"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Name     string    `json:"name"`
	Provider struct {
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
	CertifiedIntegratorCompliant struct {
		DONE      int `json:"DONE"`
		NOSTATUS  int `json:"NO_STATUS"`
		NOTDONE   int `json:"NOT_DONE"`
		NOTNEEDED int `json:"NOT_NEEDED"`
	} `json:"certifiedIntegratorCompliant"`
	Attributes []struct {
		Type          string      `json:"type"`
		Id            string      `json:"id"`
		Name          string      `json:"name"`
		Enabled       bool        `json:"enabled"`
		EntityType    string      `json:"entityType"`
		AttributeType string      `json:"attributeType"`
		DefaultValue  interface{} `json:"defaultValue"`
		Value         interface{} `json:"value"`
		Links         []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
		PossibleValues []bool `json:"possibleValues,omitempty"`
	} `json:"attributes"`
	Domain struct {
		Type    string    `json:"type"`
		Id      string    `json:"id"`
		Created time.Time `json:"created"`
		Updated time.Time `json:"updated"`
		Name    string    `json:"name"`
		Links   []struct {
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
