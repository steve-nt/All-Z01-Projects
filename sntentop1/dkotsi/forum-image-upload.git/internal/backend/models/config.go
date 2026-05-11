package models

type Config struct {
	Tls             TLS       `json:"Tls"`
	Api             API       `json:"API"`
	Durations       Durations `json:"Durations"`
	CookieName      string    `json:"CookieName"`
	HandlersTimeOut int       `json:"HandlersTimeout"`
}

type TLS struct {
	Certification string `json:"CERTIFICATION"`
	Key           string `json:"KEY"`
}

type API struct {
	Api_base_url  string `json:"API_BASE_URL"`
	Api_key       string `json:"API_KEY"`
	Header_format string `json:"HEADER_FORMAT"`
}

type Durations struct {
	ClientTimeOut        int     `json:"ClientTimeOut"`
	IdleTTL              int     `json:"IdleSessionExpiration"`
	AbsoluteTTL          int     `json:"AbsoluteSessionExpiration"`
	DbCleanUpRate        int     `json:"DBCleanUpRate"`
	RateLimit            float64 `json:"RateLimit"`
	RateLimitRefreshRate int     `json:"RateLimiterRefreshRate"`
}
