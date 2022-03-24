package discovery

import "net/url"

const (
	OpenApiV2Source = "OpenApi2"
	OpenApiV3Source = "OpenApi3"
)

// WebApiDescription description of a web API
// holds information about:
// - endpoints (requests, responses, parameters, etc.)
// - global information (URL, version, description, title, etc.)
// - how to authenticate with the API (security information)
type WebApiDescription struct {
	DiscoverySource string
	DiscoveryDoc    url.URL
	Host            string
	BaseUrl         url.URL
	Version         string
	Title           string
	Description     string
	Endpoints       []Endpoint
	SecuritySchemes []SecuritySchema
}

// Endpoint information about an endpoint inside an API
type Endpoint struct {
	Path        string
	Method      string // POST, GET, etc.
	Body        Body
	Parameters  []Parameter
	Consumes    string
	Produces    string
	Responses   []Response
	Summary     string
	Description string
}

// Body for a http request
type Body struct {
	Description string
	Required    bool
	Content     []Content
}

// Parameter a parameter inside an endpoint
type Parameter struct {
	Name        string
	In          string // body, header, etc.
	Required    bool
	Description string
	ParamType   string
	Schema      Schema
}

// Response a response for a request to an endpoint
type Response struct {
	Code        int
	Description string
	Content     []Content
}

// Content content in a request
// example:
// ContentType could for example be 'application/json
// And Schema holds information about the content like its format (keys and types) and if it's nullable
type Content struct {
	ContentType string
	Schema      Schema
}

// Schema holds information about an object
// like a json object, xml object, html document, etc.
type Schema struct {
	Key        string
	Type       string
	Format     string
	Nullable   bool
	AllowEmpty bool
	Example    any
	Properties []Schema
}

// SecuritySchema defines a security scheme
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#security-scheme-object
type SecuritySchema struct {
	Key string

	Type             string
	Description      string
	Name             string
	In               string
	BearerFormat     string
	OpenIdConnectUrl string
	Flows            []OAuthFlow
	// Scheme string?
}

// OAuthFlow configuration for a OAuth flow
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#oauthFlowsObject
type OAuthFlow struct {
	GrantType        string
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

const (
	Implicit          = "OAuth2Implicit"
	AuthorizationCode = "OAuth2AuthCode"
	Password          = "OAuth2Password"
	ClientCredentials = "OAuth2ClientCredentials"

	BasicSecSchemaType  = "http"
	ApiKeySecSchemaType = "apiKey"
	OAuth2SecSchemaType = "oauth2"
)

/* type SecuritySchemeReference struct {
	Type string
	Description string
	Name string
	In string
	BearerFormat string
	OpenIdConnectUrl string
	// Flows?
	// Scheme string?
} */
