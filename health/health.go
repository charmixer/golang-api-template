package health

import (
	"context"
	"syscall"
	"time"
)

type Status string

type Health struct {
	Status      Status             `json:"status"`
	Version     string             `json:"version,omitempty"`
	ReleaseID   string             `json:"releaseId,omitempty"`
	Notes       []string           `json:"notes,omitempty"`
	Output      string             `json:"output,omitempty"`
	ServiceID   string             `json:"serviceId,omitempty"`
	Description string             `json:"description,omitempty"`
	Checks      map[string][]Check `json:"checks,omitempty"`
	Links       map[string]string  `json:"links,omitempty"`
}

type Check struct {
	ComponentID       string            `json:"componentId,omitempty"`
	ComponentType     string            `json:"componentType,omitempty"`
	ObservedValue     interface{}       `json:"observedValue,omitempty"`
	ObservedUnit      string            `json:"observedUnit,omitempty"`
	Status            Status            `json:"status" example:"pass"`
	AffectedEndpoints []string          `json:"affectedEndpoints,omitempty"`
	Time              string            `json:"time,omitempty" example:"2019-02-20T22:01:44,654015561+00:00"`
	Output            string            `json:"output,omitempty"`
	Links             map[string]string `json:"links,omitempty"`
}

const (
	Pass Status = "pass"
	Fail Status = "fail"
	Warn Status = "warn"
)

type HealthCheck func() (string, Check)

func WithUptimeCheck(ctx context.Context) HealthCheck {
	si := &syscall.Sysinfo_t{}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return func() (string, Check) {
		return "uptime", Check{
			ComponentType: "system",
			ObservedValue: si.Uptime,
			ObservedUnit:  "s",
			Status:        Pass,
			Time:          now,
		}
	}
}

type HealthChecker struct {
	checks []HealthCheck
}

func (hc *HealthChecker) Check() {

}

func New(ctx context.Context) *HealthChecker {
	hc := &HealthChecker{}

	return hc
}
