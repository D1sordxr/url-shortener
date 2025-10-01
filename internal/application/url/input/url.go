package input

type Create struct {
	URL   string  `json:"url"`
	Alias *string `json:"alias"`
}

type GetForRedirect struct {
	Alias     string  `json:"alias"`
	UserID    *string `json:"user_id"`
	UserAgent *string `json:"user_agent"`
	IpAddress *string `json:"ip_address"`
	Referer   *string `json:"referer"`
}
