package entity

type (
	// InstallationHeader ... structure of a installation header
	InstallationHeader struct {
		ID string
	}

	// PoppoHeader ... x original header
	PoppoHeader struct {
		SessionID  string
		CSRFToken  string
		AppVersion string
	}
)
