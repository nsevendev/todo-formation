package job

import "time"

type Job struct {
	Name     string            `json:"name"`
	Payload  map[string]string `json:"payload"`
	Retry    int               `json:"retry"`
	MaxRetry int               `json:"max_retry"`
	Created  time.Time         `json:"created"`
}
