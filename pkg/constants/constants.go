package constants

// Configuration defaults
const (
	// DefaultTimeoutSeconds is the default timeout for JSON-RPC requests
	DefaultTimeoutSeconds = 30

	// DefaultJSONRPCVersion is the default JSON-RPC version
	DefaultJSONRPCVersion = "2.0"
)

// HTTP headers
const (
	// HeaderContentType is the content type header for JSON-RPC requests
	HeaderContentType = "application/json"
)

// Output formatting
const (
	// MaxNameLength is the maximum length for request names in table output
	MaxNameLength = 25

	// MaxMethodLength is the maximum length for method names in table output
	MaxMethodLength = 30

	// MaxConfigLength is the maximum length for config names in table output
	MaxConfigLength = 15

	// BoxWidth is the width of the detailed output box
	BoxWidth = 78

	// BoxContentWidth is the width of content inside the detailed output box
	BoxContentWidth = 76
)

// HTTP status codes
const (
	// MinClientErrorStatus is the minimum HTTP status code considered a client error
	MinClientErrorStatus = 400
)
