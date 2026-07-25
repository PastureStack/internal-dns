package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/PastureStack/internal-dns/internal/logging"
)

const maxResponseBytes = 32 << 20

type Client interface {
	GetRegionName() (string, error)
	GetSelfHost() (Host, error)
	GetServiceInLocalEnvironment(string, string) (Service, error)
	GetServiceInLocalRegion(string, string, string) (Service, error)
	GetServiceFromRegionEnvironment(string, string, string, string) (Service, error)
	GetServices() ([]Service, error)
	GetContainers() ([]Container, error)
	OnChange(int, func(string))
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewClientAndWait(baseURL string) (Client, error) {
	parsed, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid metadata URL")
	}
	client := &HTTPClient{
		baseURL: parsed.String(),
		client:  &http.Client{Timeout: 15 * time.Second},
	}
	var version string
	if err := client.get("/version", &version); err != nil {
		return nil, fmt.Errorf("metadata service is unavailable: %w", err)
	}
	return client, nil
}

func (c *HTTPClient) get(path string, target any) error {
	request, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))
		return fmt.Errorf("metadata request %s returned %s", path, response.Status)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return err
	}
	if len(body) > maxResponseBytes {
		return fmt.Errorf("metadata response %s exceeds %d bytes", path, maxResponseBytes)
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode metadata response %s: %w", path, err)
	}
	return nil
}

func (c *HTTPClient) GetRegionName() (value string, err error) {
	err = c.get("/region_name", &value)
	return
}

func (c *HTTPClient) GetSelfHost() (value Host, err error) {
	err = c.get("/self/host", &value)
	return
}

func (c *HTTPClient) GetServices() (value []Service, err error) {
	err = c.get("/services", &value)
	return
}

func (c *HTTPClient) GetContainers() (value []Container, err error) {
	err = c.get("/containers", &value)
	return
}

func (c *HTTPClient) GetServiceInLocalEnvironment(stackName, serviceName string) (value Service, err error) {
	path := "/stacks/" + url.PathEscape(stackName) + "/services/" + url.PathEscape(serviceName)
	err = c.get(path, &value)
	return
}

func (c *HTTPClient) getEnvironments() (value []Environment, err error) {
	err = c.get("/environments", &value)
	return
}

func (c *HTTPClient) GetServiceFromRegionEnvironment(regionName, environmentName, stackName, serviceName string) (Service, error) {
	environments, err := c.getEnvironments()
	if err != nil {
		return Service{}, err
	}
	for _, environment := range environments {
		if !strings.EqualFold(regionName, environment.RegionName) || !strings.EqualFold(environmentName, environment.Name) {
			continue
		}
		for _, service := range environment.Services {
			if strings.EqualFold(stackName, service.StackName) && strings.EqualFold(serviceName, service.Name) {
				return service, nil
			}
		}
	}
	return Service{}, fmt.Errorf("metadata service not found in requested region and environment")
}

func (c *HTTPClient) GetServiceInLocalRegion(environmentName, stackName, serviceName string) (Service, error) {
	regionName, err := c.GetRegionName()
	if err != nil {
		return Service{}, err
	}
	return c.GetServiceFromRegionEnvironment(regionName, environmentName, stackName, serviceName)
}

func (c *HTTPClient) OnChange(intervalSeconds int, callback func(string)) {
	version := "init"
	for {
		path := fmt.Sprintf("/version?wait=true&value=%s&maxWait=%d", url.QueryEscape(version), intervalSeconds)
		var next string
		if err := c.get(path, &next); err != nil {
			log.Errorf("metadata version check failed: %v", err)
			time.Sleep(time.Duration(intervalSeconds) * time.Second)
			continue
		}
		if next != version {
			version = next
			callback(version)
		}
	}
}
