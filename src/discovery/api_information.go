package discovery

import "net/url"

const (
	OpenApiV2Source = "OpenApi2"
	OpenApiV3Source = "OpenApi3"
)

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

type Body struct {
	Description string
	Required    bool
	Content     []Content
}

type Parameter struct {
	Name        string
	In          string // body, header, etc.
	Required    bool
	Description string
	ParamType   string
	Schema      Schema
}

type Response struct {
	Code        int
	Description string
	Content     []Content
}

type Content struct {
	ContentType string
	Schema      Schema
}

type Schema struct {
	Key        string
	Type       string
	Format     string
	Nullable   bool
	AllowEmpty bool
	Example    interface{}
	Properties []Schema
}

// SecuritySchema https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#security-scheme-object
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
