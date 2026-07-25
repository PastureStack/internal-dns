package metadata

type HealthCheck struct {
	Port int `json:"port"`
}

type Service struct {
	Name               string            `json:"name"`
	StackName          string            `json:"stack_name"`
	Kind               string            `json:"kind"`
	Hostname           string            `json:"hostname"`
	Vip                string            `json:"vip"`
	UUID               string            `json:"uuid"`
	ExternalIps        []string          `json:"external_ips"`
	Containers         []Container       `json:"containers"`
	Links              map[string]string `json:"links"`
	HealthCheck        HealthCheck       `json:"health_check"`
	PrimaryServiceName string            `json:"primary_service_name"`
	EnvironmentUUID    string            `json:"environment_uuid"`
}

type Container struct {
	Name                     string            `json:"name"`
	PrimaryIp                string            `json:"primary_ip"`
	ServiceName              string            `json:"service_name"`
	StackName                string            `json:"stack_name"`
	HostUUID                 string            `json:"host_uuid"`
	UUID                     string            `json:"uuid"`
	State                    string            `json:"state"`
	HealthState              string            `json:"health_state"`
	Dns                      []string          `json:"dns"`
	DnsSearch                []string          `json:"dns_search"`
	NetworkFromContainerUUID string            `json:"network_from_container_uuid"`
	Links                    map[string]string `json:"links"`
	EnvironmentUUID          string            `json:"environment_uuid"`
}

type Host struct {
	Name            string `json:"name"`
	UUID            string `json:"uuid"`
	EnvironmentUUID string `json:"environment_uuid"`
}

type Environment struct {
	Name       string    `json:"name"`
	RegionName string    `json:"region_name"`
	Services   []Service `json:"services"`
}
