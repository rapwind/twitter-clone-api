package constant

import "time"

const (
	// UserSessionPrefix ...
	UserSessionPrefix = "session@"
	// SessionExpires is 3 days
	SessionExpires = 3 * 24 * time.Hour

	//CsrfTokenPrefix ...
	CsrfTokenPrefix = "csrftkn@"
	// CsrfTokenExpires is 12 hours
	CsrfTokenExpires = 12 * time.Hour
)
