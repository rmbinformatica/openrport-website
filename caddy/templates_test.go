package caddy

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldMakeNewRouteRequestJSON(t *testing.T) {
	tmpl := template.New("NRR")
	tmpl, err := tmpl.Parse(NewRouteRequestTemplate)
	require.NoError(t, err)

	nrr := &NewRouteRequest{
		RouteID:                 "new_route_id",
		TargetTunnelHost:        "target_tunnel_host",
		TargetTunnelPort:        "target_tunnel_port",
		UpstreamProxySubdomain:  "upstream_proxy_subdomain",
		UpstreamProxyBaseDomain: "upstream_proxy_basedomain",
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, nrr)
	require.NoError(t, err)

	templateText := b.String()

	assert.Contains(t, templateText, `"@id": "new_route_id"`)
	assert.Contains(t, templateText, `"handler": "reverse_proxy"`)
	assert.Contains(t, templateText, `"dial": "target_tunnel_host:target_tunnel_port"`)
	assert.Contains(t, templateText, `"upstream_proxy_subdomain.upstream_proxy_basedomain"`)
}

func TestShouldParseTemplates(t *testing.T) {
	cases := []struct {
		name     string
		template string
	}{
		{
			name:     "global settings",
			template: globalSettingsTemplate,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := template.New(tc.name)
			_, err := tmpl.Parse(tc.template)
			assert.NoError(t, err)
		})
	}
}

func TestShouldMakeGlobalSettingsText(t *testing.T) {
	tmpl := template.New("GS")
	tmpl, err := tmpl.Parse(globalSettingsTemplate)
	require.NoError(t, err)

	gs := &GlobalSettings{
		LogLevel:    "ERROR",
		AdminSocket: "/tmp/caddyadmin.sock",
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, gs)
	require.NoError(t, err)

	templateText := b.String()
	// fmt.Printf("templateText = %+v\n", templateText)

	assert.Contains(t, templateText, "level ERROR")
	assert.Contains(t, templateText, "admin unix//tmp/caddyadmin.sock")
}

func TestShouldMakeDefaultVirtualHostText(t *testing.T) {
	tmpl := template.New("DVH")
	tmpl, err := tmpl.Parse(defaultVirtualHost)
	require.NoError(t, err)

	dvh := &DefaultVirtualHost{
		ListenAddress: "listen_address",
		ListenPort:    "listen_port",
		CertsFile:     "certs_file",
		KeyFile:       "key_file",
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, dvh)
	require.NoError(t, err)

	templateText := b.String()

	assert.Contains(t, templateText, "https://listen_address:listen_port")
	assert.Contains(t, templateText, "tls certs_file key_file")
}

func TestShouldMakeAPIReverseProxySettingsText(t *testing.T) {
	tmpl := template.New("ARP")
	tmpl, err := tmpl.Parse(apiReverseProxySettingsTemplate)
	require.NoError(t, err)

	arp := &APIReverseProxySettings{
		CertsFile:    "certs_file",
		KeyFile:      "key_file",
		ProxyDomain:  "proxy_domain",
		ProxyPort:    "proxy_port",
		APIDomain:    "api_domain",
		APIScheme:    "api_scheme",
		APIIPAddress: "api_ip_address",
		APIPort:      "api_port",
		ProxyLogFile: "proxy_log_file",
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, arp)
	require.NoError(t, err)

	templateText := b.String()

	assert.Contains(t, templateText, "https://proxy_domain:proxy_port")
	assert.Contains(t, templateText, "tls certs_file key_file")
	assert.Contains(t, templateText, "reverse_proxy api_scheme://api_ip_address:api_port")
	assert.Contains(t, templateText, "output file proxy_log_file")
}

func TestShouldMakeExternalReverseProxyText(t *testing.T) {
	tmpl := template.New("ERP")
	tmpl, err := tmpl.Parse(externalReverseProxyTemplate)
	require.NoError(t, err)

	erp := &ExternalReverseProxy{
		CertsFile:        "certs_file",
		KeyFile:          "key_file",
		BaseDomain:       "base_domain",
		Subdomain:        "sub_domain",
		AllowedIPAddress: "allowed_ip_address",
		TunnelScheme:     "tunnel_scheme",
		TunnelIPAddress:  "tunnel_ip_address",
		TunnelPort:       "tunnel_port",
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, erp)
	require.NoError(t, err)

	templateText := b.String()

	assert.Contains(t, templateText, "https://sub_domain.base_domain")
	assert.Contains(t, templateText, "tls certs_file key_file")
	assert.Contains(t, templateText, "@onlyif remote_ip allowed_ip_address")
	assert.Contains(t, templateText, "reverse_proxy @onlyif tunnel_scheme://tunnel_ip_address:tunnel_port")
}

func TestShouldMakeAll(t *testing.T) {
	tmpl := template.New("ALL")

	tmpl, err := tmpl.Parse(globalSettingsTemplate)
	require.NoError(t, err)

	tmpl, err = tmpl.Parse(defaultVirtualHost)
	require.NoError(t, err)

	tmpl, err = tmpl.Parse(apiReverseProxySettingsTemplate)
	require.NoError(t, err)

	tmpl, err = tmpl.Parse(externalReverseProxyTemplate)
	require.NoError(t, err)

	tmpl, err = tmpl.Parse(allTemplate)
	require.NoError(t, err)

	gs := &GlobalSettings{
		LogLevel:    "ERROR",
		AdminSocket: "/tmp/caddyadmin.sock",
	}

	dvh := &DefaultVirtualHost{
		ListenAddress: "listen_address",
		ListenPort:    "listen_port",
		CertsFile:     "certs_file",
		KeyFile:       "key_file",
	}

	arp := &APIReverseProxySettings{
		CertsFile:    "certs_file",
		KeyFile:      "key_file",
		ProxyDomain:  "proxy_domain",
		ProxyPort:    "proxy_port",
		APIDomain:    "api_domain",
		APIScheme:    "api_scheme",
		APIIPAddress: "api_ip_address",
		APIPort:      "api_port",
		ProxyLogFile: "proxy_log_file",
	}

	erp := &ExternalReverseProxy{
		CertsFile:        "certs_file",
		KeyFile:          "key_file",
		BaseDomain:       "base_domain",
		Subdomain:        "sub_domain",
		AllowedIPAddress: "allowed_ip_address",
		TunnelScheme:     "tunnel_scheme",
		TunnelIPAddress:  "tunnel_ip_address",
		TunnelPort:       "tunnel_port",
	}

	c := ExecBaseConfig{
		GlobalSettings:          gs,
		DefaultVirtualHost:      dvh,
		APIReverseProxySettings: arp,
		ReverseProxies:          []ExternalReverseProxy{*erp},
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, c)
	require.NoError(t, err)

	templateText := b.String()

	assert.Contains(t, templateText, "https://listen_address:listen_port")
	assert.Contains(t, templateText, "https://proxy_domain:proxy_port")
	assert.Contains(t, templateText, "https://sub_domain.base_domain")
}
