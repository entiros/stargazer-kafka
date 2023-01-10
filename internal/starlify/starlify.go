package starlify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"strings"
)

// Client Starlify system client configuration
type Client struct {
	BaseUrl      string
	ApiKey       string
	AgentId      string
	MiddlewareId string
	resty        *resty.Client
}

type TopicEndpoint struct {
	Name   string
	ID     string
	Prefix string
}

func (starlify *Client) GetRestyClient() *resty.Client {
	if starlify.resty == nil {
		starlify.resty = resty.New()
	}

	return starlify.resty
}

// get performs GET request to path and return parsed response
func (starlify *Client) get(ctx context.Context, path string, returnType any) error {
	// GET request
	requestPath := starlify.BaseUrl + path
	log.Printf("Performing get to : %s", requestPath)
	response, err := starlify.GetRestyClient().R().
		SetContext(ctx).
		SetHeader("X-API-KEY", starlify.ApiKey).
		Get(requestPath)
	if err != nil {
		return err
	}

	if response.StatusCode() == http.StatusOK {
		err = json.Unmarshal(response.Body(), &returnType)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("error while performing request to %s, error: %v", response.Request.URL, response.Error())
	}

	return nil

}

// post performs POST request to path and return parsed response
func (starlify *Client) post(ctx context.Context, path string, body any, returnType any) error {

	resp, err := starlify.GetRestyClient().R().
		SetContext(ctx).
		SetHeader("X-API-KEY", starlify.ApiKey).
		SetBody(body).
		SetResult(returnType).
		Post(starlify.BaseUrl + path)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("%s", resp.Status())
	}
	return nil
}

// post performs POST request to path and return parsed response
func (starlify *Client) delete(ctx context.Context, path string) error {

	_, err := starlify.GetRestyClient().R().
		SetContext(ctx).
		SetHeader("X-API-KEY", starlify.ApiKey).
		Delete(starlify.BaseUrl + path)
	if err != nil {
		return err
	}

	return err
}

// patch performs PATCH request to path and return parsed response
func (starlify *Client) patch(ctx context.Context, path string, body any, returnType any) error {

	response, err := starlify.GetRestyClient().R().
		SetContext(ctx).
		SetHeader("X-API-KEY", starlify.ApiKey).
		SetBody(body).
		Patch(starlify.BaseUrl + path)
	if err != nil {
		return err
	}

	if response.IsError() {
		return fmt.Errorf("%s", response.Status())
	}

	// Parse response
	err = json.Unmarshal(response.Body(), &returnType)
	if err != nil {
		return err
	}

	return err
}

func (starlify *Client) GetTopics(ctx context.Context) ([]TopicEndpoint, error) {

	var middleware Middleware
	path := fmt.Sprintf("/middlewares/%s", starlify.MiddlewareId)

	err := starlify.get(ctx, path, &middleware)
	if err != nil {
		return nil, err
	}

	var topics []TopicEndpoint
	for _, endpoint := range middleware.Endpoints {
		e := strings.TrimSpace(endpoint.Name)
		if strings.HasPrefix(e, strings.TrimSpace(middleware.KafkaPrefix)) {

			topics = append(topics, TopicEndpoint{
				Name:   e,
				ID:     endpoint.Id,
				Prefix: middleware.KafkaPrefix,
			})
		}
	}
	return topics, nil
}

func (starlify *Client) CreateTopic(ctx context.Context, topic string) error {

	endpoint := EndpointRequest{
		Name: topic,
	}
	path := fmt.Sprintf("/middlewares/%s/endpoints", starlify.MiddlewareId)

	var response EndpointResponse
	err := starlify.post(ctx, path, &endpoint, &response)
	return err

}

func (starlify *Client) DeleteTopic(ctx context.Context, endpoint TopicEndpoint) error {

	path := fmt.Sprintf("/endpoints/%s", endpoint.ID)

	return starlify.delete(ctx, path)

}

// GetServices will get all services for system
func (starlify *Client) GetServices(ctx context.Context) ([]Service, error) {
	log.Printf("Get services for system %s", starlify.MiddlewareId)

	var services []Service = nil

	var totalPages = 1

	var serviceIndex = 0
	for page := 0; page < totalPages; page++ {
		log.Printf("Fetching services page %d", page)

		// Get services page
		var servicesPage ServicesPage
		path := fmt.Sprintf("/systems/%s/services?page=%d", starlify.MiddlewareId, page)
		err := starlify.get(ctx, path, &servicesPage)
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
func (starlify *Client) CreateService(ctx context.Context, name string) (*Service, error) {
	log.Printf("Create service '%s' in system %s", name, starlify.MiddlewareId)

	var service Service
	err := starlify.post(ctx, "/systems/"+starlify.MiddlewareId+"/services", &ServiceRequest{Name: name}, &service)
	if err != nil {
		return nil, err
	}

	log.Printf("Service '%s' (%s) created in system %s", service.Name, service.Id, starlify.MiddlewareId)
	return &service, nil
}

// Ping will perform an agent update without data - nothing will be updated except 'lastSeen'
func (starlify *Client) Ping(ctx context.Context) error {
	var agent Agent
	err := starlify.patch(ctx, "/agents/"+starlify.AgentId, AgentRequest{}, &agent)
	if err != nil {
		return err
	}

	log.Printf("Ping")
	return nil
}

// GetAgent will return the Starlify agent
func (starlify *Client) GetAgent(ctx context.Context) (*Agent, error) {
	log.Printf("Get agent %s", starlify.AgentId)

	var agent Agent
	err := starlify.get(ctx, "/agents/"+starlify.AgentId, &agent)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// UpdateDetails will update the agent details
func (starlify *Client) UpdateDetails(ctx context.Context, details Details) error {
	log.Println("Updating details")
	var agent Agent
	err := starlify.patch(ctx, "/agents/"+starlify.AgentId, AgentRequest{Details: &details}, &agent)
	if err != nil {
		return err
	}
	return nil
}

// ReportError will update the error field in the Starlify agent
func (starlify *Client) ReportError(ctx context.Context, message string) error {
	var agent Agent
	err := starlify.patch(ctx, "/agents/"+starlify.AgentId, AgentRequest{Error: message}, &agent)
	if err != nil {
		return err
	}
	return nil
}

// ClearError will Starlify agent error field to empty string
func (starlify *Client) ClearError(ctx context.Context) error {
	return starlify.ReportError(ctx, "")
}
