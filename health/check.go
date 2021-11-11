package health

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

func WithUptimeCheck(componentId string) HealthCheck {
	return func(ctx context.Context, result chan healthCheckResult) {
		up, err := host.Uptime()

		status := Pass
		if err != nil {
			status = Warn
		}

		result <- healthCheckResult{
			SystemId:    "uptime",
			ComponentId: componentId,
			Check: Check{
				ComponentType: "system",
				ObservedValue: int64(up),
				ObservedUnit:  "s",
				Status:        status,
			},
		}
	}
}

func WithCpuCheck(componentId string) HealthCheck {
	return func(ctx context.Context, result chan healthCheckResult) {
		cpu, err := cpu.Percent(5*time.Second, false)

		status := Pass
		if err != nil {
			status = Warn
		}

		result <- healthCheckResult{
			SystemId:    "cpu:utilalization",
			ComponentId: componentId,
			Check: Check{
				ComponentType: "system",
				ObservedValue: fmt.Sprintf("%.2f", cpu[0]),
				ObservedUnit:  "%",
				Status:        status,
			},
		}
	}
}

func WithNumGcCheck(componentId string) HealthCheck {
	return func(ctx context.Context, result chan healthCheckResult) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		status := Pass

		result <- healthCheckResult{
			SystemId:    "mem:utilalization",
			ComponentId: componentId,
			Check: Check{
				ComponentType: "system",
				ObservedValue: fmt.Sprintf("%d", m.NumGC),
				ObservedUnit:  "num of GC cycles",
				Status:        status,
			},
		}
	}

}

func WithMemObtainedCheck(componentId string) HealthCheck {
	return func(ctx context.Context, result chan healthCheckResult) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		status := Pass

		result <- healthCheckResult{
			SystemId:    "mem:utilalization",
			ComponentId: componentId,
			Check: Check{
				ComponentType: "system",
				ObservedValue: fmt.Sprintf("%d", m.Sys),
				ObservedUnit:  "bytes",
				Status:        status,
			},
		}
	}
}

func WithMemTotalAllocCheck(componentId string) HealthCheck {
	return func(ctx context.Context, result chan healthCheckResult) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		status := Pass

		result <- healthCheckResult{
			SystemId:    "mem:utilalization",
			ComponentId: componentId,
			Check: Check{
				ComponentType: "system",
				ObservedValue: fmt.Sprintf("%d", m.TotalAlloc),
				ObservedUnit:  "bytes",
				Status:        status,
			},
		}
	}
}