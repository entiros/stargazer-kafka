package starlify

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

// Client Starlify system client configuration
type Client struct {
	BaseUrl  string
	ApiKey   string
	AgentId  string
	SystemId string
}

// get performs GET request to path and return parsed response
func (starlify *Client) get(path string, returnType any) error {
	httpClient := resty.New()

	// GET request
	response, err := httpClient.R().
		SetHeader("X-API-KEY", starlify.ApiKey).
		Get(starlify.BaseUrl + path)
	if err != nil {
		return err
	}

	// Parse response
	err = json.Unmarshal(response.Body(), &returnType)
	if err != nil {
		return err
	}

	return err
}

// post performs POST request to path and return parsed response
func (starlify *Client) post(path string, body any, returnType any) error {
	httpClient := resty.New()

	// POST request
	response, err := httpClient.R().
		SetHeader("X-API-KEY", starlify.ApiKey).
		SetBody(body).
		Post(starlify.BaseUrl + path)
	if err != nil {
		return err
	}

	// Parse response
	err = json.Unmarshal(response.Body(), &returnType)
	if err != nil {
		return err
	}

	return err
}

// patch performs PATCH request to path and return parsed response
func (starlify *Client) patch(path string, body any, returnType any) error {
	httpClient := resty.New()

	// POST request
	response, err := httpClient.R().
		SetHeader("X-API-KEY", starlify.ApiKey).
		SetBody(body).
		Patch(starlify.BaseUrl + path)
	if err != nil {
		return err
	}

	// Parse response
	err = json.Unmarshal(response.Body(), &returnType)
	if err != nil {
		return err
	}

	return err
}

// GetServices will get all services for system
func (starlify *Client) GetServices() ([]Service, error) {
	log.Printf("Get services for system %s", starlify.SystemId)

	var services []Service = nil

	var totalPages = 1

	var serviceIndex = 0
	for page := 0; page < totalPages; page++ {
		log.Printf("Fetching services page %d", page)

		// Get services page
		var servicesPage ServicesPage
		path := fmt.Sprintf("/systems/%s/services?page=%d", starlify.SystemId, page)
		err := starlify.get(path, &servicesPage)
		if err != nil {
			return nil, err
		}

		// Update total pages
		totalPages = servicesPage.Page.TotalPages

		// Initialize return array
		if services == nil {
			log.Printf("%d services to be fetched (%d pages)", servicesPage.Page.TotalElements, totalPages)
			services = make([]Service, servicesPage.Page.TotalElements)
		}

		// Add services to response
		for _, service := range servicesPage.Services {
			services[serviceIndex] = service
			serviceIndex++
		}
	}

	log.Printf("%d services fetched", len(services))
	return services, nil
}

// CreateService will create and return new service
func (starlify *Client) CreateService(name string) (*Service, error) {
	log.Printf("Create service '%s' in system %s", name, starlify.SystemId)

	var service Service
	err := starlify.post("/systems/"+starlify.SystemId+"/services", &ServiceRequest{Name: name}, &service)
	if err != nil {
		return nil, err
	}

	log.Printf("Service '%s' (%s) created in system %s", service.Name, service.Id, starlify.SystemId)
	return &service, nil
}

// Ping will perform an agent update without data - nothing will be updated except 'lastSeen'
func (starlify *Client) Ping() error {
	var agent Agent
	err := starlify.patch("/agents/"+starlify.AgentId, AgentRequest{}, &agent)
	if err != nil {
		return err
	}

	log.Printf("Ping")
	return nil
}

// GetAgent will return the Starlify agent
func (starlify *Client) GetAgent() (*Agent, error) {
	log.Printf("Get agent %s", starlify.AgentId)

	var agent Agent
	err := starlify.get("/agents/"+starlify.AgentId, &agent)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// UpdateDetails will update the agent details
func (starlify *Client) UpdateDetails(details Details) error {
	log.Println("Updating details")
	var agent Agent
	err := starlify.patch("/agents/"+starlify.AgentId, AgentRequest{Details: &details}, &agent)
	if err != nil {
		return err
	}
	return nil
}

// ReportError will update the error field in the Starlify agent
func (starlify *Client) ReportError(message string) error {
	var agent Agent
	err := starlify.patch("/agents/"+starlify.AgentId, AgentRequest{Error: message}, &agent)
	if err != nil {
		return err
	}
	return nil
}

// ClearError will Starlify agent error field to empty string
func (starlify *Client) ClearError() error {
	return starlify.ReportError("")
}
