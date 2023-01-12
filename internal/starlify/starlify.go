package starlify

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/entiros/stargazer-kafka/internal/log"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
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

func (starlify *Client) RestyClient() *resty.Client {
	if starlify.resty == nil {
		starlify.resty = resty.New()
	}

	return starlify.resty
}

// get performs GET request to path and return parsed response
func (starlify *Client) get(ctx context.Context, path string, returnType any) error {

	RetryCount := 6
	Timeout := 10

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(Timeout)*time.Duration(RetryCount)+5)
	defer cancel()

	requestPath := starlify.BaseUrl + path
	log.Logger.Debugf("Performing get to: '%s'", requestPath)

	var retryCounter int

	response, err := starlify.
		RestyClient().
		SetTimeout(time.Duration(Timeout)*time.Second).
		SetRetryCount(RetryCount).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			retry := len(response.Body()) == 0 && response.StatusCode() == http.StatusOK
			if retry {
				log.Logger.Debugf("Retry %d GET %s: %s/%d: %v ", retryCounter, requestPath, response.Status(), response.StatusCode(), err)
				retryCounter++
			}
			return retry
		}).
		R().
		SetContext(ctx).
		SetHeader("X-API-KEY", starlify.ApiKey).
		Get(requestPath)

	if err != nil {
		log.Logger.Errorf("Failed GET request to %s. error: -->%v<--", requestPath, err)
		return err
	}

	if response.StatusCode() == http.StatusOK && len(response.Body()) >= 2 {
		err = json.Unmarshal(response.Body(), &returnType)
		if err != nil {
			log.Logger.Errorf("Failed to unmarshal response: %s. Err: %v", string(response.Body()), err)
			return err
		}

	} else {
		return fmt.Errorf("error while performing request to %s, status: %s:%d, error: %v", response.Request.URL, response.Status(), response.StatusCode(), response.Error())
	}

	log.Logger.Debugf("Retry count: %d for GET %s", retryCounter, requestPath)

	return nil

}

// post performs POST request to path and return parsed response
func (starlify *Client) post(ctx context.Context, path string, body any, returnType any) error {

	resp, err := starlify.RestyClient().R().
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

	_, err := starlify.RestyClient().R().
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

	response, err := starlify.RestyClient().R().
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
	log.Logger.Debugf("Get services for system %s", starlify.MiddlewareId)

	var services []Service = nil

	var totalPages = 1

	var serviceIndex = 0
	for page := 0; page < totalPages; page++ {
		log.Logger.Debugf("Fetching services page %d", page)

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
			log.Logger.Debugf("%d services to be fetched (%d pages)", servicesPage.Page.TotalElements, totalPages)
			services = make([]Service, servicesPage.Page.TotalElements)
		}

		// Add services to response
		for _, service := range servicesPage.Services {
			services[serviceIndex] = service
			serviceIndex++
		}
	}

	log.Logger.Debugf("%d services fetched", len(services))
	return services, nil
}

// CreateService will create and return new service
func (starlify *Client) CreateService(ctx context.Context, name string) (*Service, error) {
	log.Logger.Debugf("Create service '%s' in system %s", name, starlify.MiddlewareId)

	var service Service
	err := starlify.post(ctx, "/systems/"+starlify.MiddlewareId+"/services", &ServiceRequest{Name: name}, &service)
	if err != nil {
		return nil, err
	}

	log.Logger.Debugf("Service '%s' (%s) created in system %s", service.Name, service.Id, starlify.MiddlewareId)
	return &service, nil
}

// Ping will perform an agent update without data - nothing will be updated except 'lastSeen'
func (starlify *Client) Ping(ctx context.Context) error {
	var agent Agent
	err := starlify.patch(ctx, "/agents/"+starlify.AgentId, AgentRequest{}, &agent)
	if err != nil {
		return err
	}

	log.Logger.Debugf("Ping ok for %s", starlify.AgentId)
	return nil
}

// GetAgent will return the Starlify agent
func (starlify *Client) GetAgent(ctx context.Context) (*Agent, error) {
	log.Logger.Debugf("Get agent %s", starlify.AgentId)

	var agent Agent
	err := starlify.get(ctx, "/agents/"+starlify.AgentId, &agent)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// UpdateDetails will update the agent details
func (starlify *Client) UpdateDetails(ctx context.Context, details Details) error {
	log.Logger.Debug("Updating details")
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
