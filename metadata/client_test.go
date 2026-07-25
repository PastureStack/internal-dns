package metadata

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientReadsRequiredMetadataResources(t *testing.T) {
	mux := http.NewServeMux()
	responses := map[string]string{
		"/2016-07-29/version":                 `"7"`,
		"/2016-07-29/region_name":             `"region-a"`,
		"/2016-07-29/self/host":               `{"uuid":"host-1","environment_uuid":"environment-1"}`,
		"/2016-07-29/services":                `[{"name":"api","stack_name":"app","uuid":"service-1"}]`,
		"/2016-07-29/containers":              `[{"name":"api-1","uuid":"container-1","primary_ip":"10.0.0.8"}]`,
		"/2016-07-29/stacks/app/services/api": `{"name":"api","stack_name":"app","uuid":"service-1"}`,
		"/2016-07-29/environments":            `[{"name":"environment-a","region_name":"region-a","services":[{"name":"api","stack_name":"app","uuid":"service-1"}]}]`,
	}
	for path, body := range responses {
		path, body := path, body
		mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(body)) })
	}
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := NewClientAndWait(server.URL + "/2016-07-29")
	if err != nil {
		t.Fatal(err)
	}
	if region, err := client.GetRegionName(); err != nil || region != "region-a" {
		t.Fatalf("region = %q, err = %v", region, err)
	}
	if host, err := client.GetSelfHost(); err != nil || host.UUID != "host-1" {
		t.Fatalf("host = %#v, err = %v", host, err)
	}
	if services, err := client.GetServices(); err != nil || len(services) != 1 {
		t.Fatalf("services = %#v, err = %v", services, err)
	}
	if containers, err := client.GetContainers(); err != nil || len(containers) != 1 {
		t.Fatalf("containers = %#v, err = %v", containers, err)
	}
	if service, err := client.GetServiceInLocalEnvironment("app", "api"); err != nil || service.UUID != "service-1" {
		t.Fatalf("local service = %#v, err = %v", service, err)
	}
	if service, err := client.GetServiceInLocalRegion("environment-a", "app", "api"); err != nil || service.UUID != "service-1" {
		t.Fatalf("regional service = %#v, err = %v", service, err)
	}
	if service, err := client.GetServiceFromRegionEnvironment("region-a", "environment-a", "app", "api"); err != nil || service.UUID != "service-1" {
		t.Fatalf("explicit regional service = %#v, err = %v", service, err)
	}
}

func TestClientRejectsInvalidURLAndOversizedResponse(t *testing.T) {
	if _, err := NewClientAndWait("metadata/2016-07-29"); err == nil {
		t.Fatal("expected invalid URL error")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`"7"` + strings.Repeat(" ", maxResponseBytes)))
	}))
	defer server.Close()
	if _, err := NewClientAndWait(server.URL); err == nil {
		t.Fatal("expected oversized response error")
	}
}

func TestClientReportsMissingRegionalService(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`"1"`)) })
	mux.HandleFunc("/environments", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`[]`)) })
	server := httptest.NewServer(mux)
	defer server.Close()
	client, err := NewClientAndWait(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetServiceFromRegionEnvironment("r", "e", "s", "x"); err == nil {
		t.Fatal("expected not found error")
	}
}
