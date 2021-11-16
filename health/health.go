package health

import (
	"context"
	"sync"
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
	Status            Status            `json:"status"`
	AffectedEndpoints []string          `json:"affectedEndpoints,omitempty"`
	Time              string            `json:"time,omitempty"`
	Output            string            `json:"output,omitempty"`
	Links             map[string]string `json:"links,omitempty"`
}

const (
	Pass Status = "pass"
	Fail Status = "fail"
	Warn Status = "warn"
)

type healthCheckResult struct {
	SystemId    string
	ComponentId string
	Check       Check
}
type HealthCheck func(ctx context.Context, result chan healthCheckResult)

type HealthChecker struct {
	mu          sync.Mutex
	version     string
	releaseId   string
	serviceId   string
	description string
	checks      []HealthCheck
	results     map[string]map[string]Check
	available   bool
}

func (hc *HealthChecker) AddCheck(checks ...HealthCheck) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.checks = append(hc.checks, checks...)
}
func (hc *HealthChecker) Reset() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	hc.results = make(map[string]map[string]Check)
}
func (hc *HealthChecker) Checks() []HealthCheck {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	return hc.checks
}
func (hc *HealthChecker) IsAvailable() bool {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	return hc.available
}
func (hc *HealthChecker) Health() Health {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	var serviceStatus = Pass

	var newResults = make(map[string][]Check, len(hc.results))
	for k1, v1 := range hc.results {
		for k2, v2 := range v1 {
			v2.ComponentID = k2
			newResults[k1] = append(newResults[k1], v2)

			if v2.Status == Fail {
				serviceStatus = Fail
			}
		}
	}

	h := Health{
		Status:      serviceStatus,
		Version:     hc.version,
		ReleaseID:   hc.releaseId,
		Notes:       []string{},
		Output:      "",
		ServiceID:   hc.serviceId,
		Description: hc.description,
		Checks:      newResults,
		Links:       map[string]string{},
	}

	return h
}
func (hc *HealthChecker) Check(ctx context.Context) {
	checks := hc.Checks()

	var channel = make(chan healthCheckResult)
	defer close(channel)
	for _, check := range checks {
		go check(ctx, channel)
	}

	for i := 0; i < len(checks); i++ {
		select {
		case result := <-channel:

			if result.SystemId == "" || result.ComponentId == "" {
				continue
			}

			if result.Check.Time == "" {
				result.Check.Time = time.Now().UTC().Format(time.RFC3339Nano)
			}

			hc.mu.Lock()
			if hc.results[result.SystemId] == nil {
				hc.results[result.SystemId] = make(map[string]Check)
			}
			hc.results[result.SystemId][result.ComponentId] = result.Check
			hc.mu.Unlock()
		}
	}

	hc.mu.Lock()
	hc.available = true
	hc.mu.Unlock()
}

func (hc *HealthChecker) SetOption(options ...HealthCheckerOption) {
	for _, opt := range options {
		opt(hc)
	}
}

type HealthCheckerOption func(e *HealthChecker)

func New(options ...HealthCheckerOption) *HealthChecker {
	hc := &HealthChecker{
		results: make(map[string]map[string]Check),
	}

	for _, opt := range options {
		opt(hc)
	}

	return hc
}

func WithVersion(version string) HealthCheckerOption {
	return func(hc *HealthChecker) {
		hc.version = version
	}
}
func WithReleaseId(rId string) HealthCheckerOption {
	return func(hc *HealthChecker) {
		hc.releaseId = rId
	}
}
func WithServiceId(sId string) HealthCheckerOption {
	return func(hc *HealthChecker) {
		hc.serviceId = sId
	}
}
func WithDescription(desc string) HealthCheckerOption {
	return func(hc *HealthChecker) {
		hc.description = desc
	}
}
