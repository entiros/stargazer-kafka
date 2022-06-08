package starlify

import (
	"gopkg.in/h2non/gock.v1"
	"reflect"
	"testing"
)

func createStarlifyClient() *Client {
	starlify := &Client{
		BaseUrl:  "http://127.0.0.1:8080/hypermedia",
		ApiKey:   "api-key-123",
		AgentId:  "agent-id-123",
		SystemId: "system-id-123",
	}

	// Intercept Starlify client
	gock.InterceptClient(starlify.getRestyClient().GetClient())

	return starlify
}

func createGock() *gock.Request {
	return gock.New("http://127.0.0.1:8080/hypermedia")
}

func TestClient_ClearError(t *testing.T) {
	defer gock.Off()

	starlify := createStarlifyClient()

	tests := []struct {
		gock    func(*gock.Request)
		name    string
		wantErr bool
	}{
		{
			func(gock *gock.Request) {
				gock.Patch("/agents/agent-id-123").
					MatchType("json").
					JSON(AgentRequest{Error: ""}).
					Reply(200).
					JSON(Agent{Id: "agent-id-123", Name: "Test agent", AgentType: "kafka"})
			},
			"Clear error",
			false,
		},
		{
			func(gock *gock.Request) {
				gock.Patch("/agents/agent-id-123").
					Reply(404)
			},
			"Not found",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gock(createGock())
			if err := starlify.ClearError(); (err != nil) != tt.wantErr {
				t.Errorf("ClearError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateService(t *testing.T) {
	defer gock.Off()

	starlify := createStarlifyClient()

	type args struct {
		name string
	}
	tests := []struct {
		gock    func(*gock.Request)
		name    string
		args    args
		want    *Service
		wantErr bool
	}{
		{
			func(gock *gock.Request) {
				gock.Post("/systems/system-id-123/services").
					MatchType("json").
					JSON(ServiceRequest{Name: "Test service"}).
					Reply(201).
					JSON(Service{Id: "service-id-123", Name: "Test service"})
			},
			"Create service",
			args{name: "Test service"},
			&Service{Id: "service-id-123", Name: "Test service"},
			false,
		},
		{
			func(gock *gock.Request) {
				gock.Post("/systems/system-id-123/services").
					Reply(404)
			},
			"System not found",
			args{name: "Test service"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gock(createGock())
			got, err := starlify.CreateService(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAgent(t *testing.T) {
	defer gock.Off()

	starlify := createStarlifyClient()

	tests := []struct {
		gock    func(*gock.Request)
		name    string
		want    *Agent
		wantErr bool
	}{
		{
			func(gock *gock.Request) {
				gock.Get("/agents/agent-id-123").
					Reply(200).
					JSON(Agent{Id: "agent-id-123", Name: "Test agent", AgentType: "kafka"})
			},
			"Get agent",
			&Agent{Id: "agent-id-123", Name: "Test agent", AgentType: "kafka"},
			false,
		},
		{
			func(gock *gock.Request) {
				gock.Get("/agents/agent-id-123").
					Reply(404)
			},
			"Not found",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gock(createGock())
			got, err := starlify.GetAgent()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAgent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetServices(t *testing.T) {
	defer gock.Off()

	starlify := createStarlifyClient()

	tests := []struct {
		gock    func(*gock.Request)
		name    string
		want    []Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.gock(createGock())
			got, err := starlify.GetServices()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	type fields struct {
		BaseUrl  string
		ApiKey   string
		AgentId  string
		SystemId string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starlify := &Client{
				BaseUrl:  tt.fields.BaseUrl,
				ApiKey:   tt.fields.ApiKey,
				AgentId:  tt.fields.AgentId,
				SystemId: tt.fields.SystemId,
			}
			if err := starlify.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ReportError(t *testing.T) {
	type fields struct {
		BaseUrl  string
		ApiKey   string
		AgentId  string
		SystemId string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starlify := &Client{
				BaseUrl:  tt.fields.BaseUrl,
				ApiKey:   tt.fields.ApiKey,
				AgentId:  tt.fields.AgentId,
				SystemId: tt.fields.SystemId,
			}
			if err := starlify.ReportError(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("ReportError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdateDetails(t *testing.T) {
	type fields struct {
		BaseUrl  string
		ApiKey   string
		AgentId  string
		SystemId string
	}
	type args struct {
		details Details
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starlify := &Client{
				BaseUrl:  tt.fields.BaseUrl,
				ApiKey:   tt.fields.ApiKey,
				AgentId:  tt.fields.AgentId,
				SystemId: tt.fields.SystemId,
			}
			if err := starlify.UpdateDetails(tt.args.details); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
