package metrics

import (
	"context"
	"log"
	"sync"
	"time"

	"docker-monitoring/db"
	"docker-monitoring/dockerclient"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemStats struct {
	Timestamp         string                       `json:"timestamp"`
	HostCPU           float64                      `json:"host_cpu"`
	HostRAM           RAMStats                     `json:"host_ram"`
	HostDisk          DiskStats                    `json:"host_disk"`
	TotalContainers   int                          `json:"total_containers"`
	RunningContainers int                          `json:"running_containers"`
	StoppedContainers int                          `json:"stopped_containers"`
	Containers        []dockerclient.ContainerInfo `json:"containers"`
}

type RAMStats struct {
	Used        uint64  `json:"used"`
	Total       uint64  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskStats struct {
	Used        uint64  `json:"used"`
	Total       uint64  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

type Collector struct {
	mu           sync.RWMutex
	currentStats SystemStats
}

var GlobalCollector *Collector

func InitCollector() *Collector {
	c := &Collector{}
	GlobalCollector = c
	return c
}

func (c *Collector) Start(ctx context.Context) {
	log.Println("Starting system metrics collector background loop...")
	ticker2s := time.NewTicker(2 * time.Second)
	ticker10s := time.NewTicker(10 * time.Second)
	defer ticker2s.Stop()
	defer ticker10s.Stop()

	// Initial collect
	c.collect(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker2s.C:
			c.collect(ctx)
		case <-ticker10s.C:
			c.saveSnapshot()
		}
	}
}

func (c *Collector) collect(ctx context.Context) {
	var stats SystemStats
	stats.Timestamp = time.Now().Format(time.RFC3339)

	// CPU %
	cpuPercents, err := cpu.PercentWithContext(ctx, 0, false)
	if err == nil && len(cpuPercents) > 0 {
		stats.HostCPU = cpuPercents[0]
	}

	// RAM
	vMem, err := mem.VirtualMemoryWithContext(ctx)
	if err == nil && vMem != nil {
		stats.HostRAM = RAMStats{
			Used:        vMem.Used,
			Total:       vMem.Total,
			UsedPercent: vMem.UsedPercent,
		}
	}

	// Disk
	diskUsage, err := disk.UsageWithContext(ctx, "/")
	if err == nil && diskUsage != nil {
		stats.HostDisk = DiskStats{
			Used:        diskUsage.Used,
			Total:       diskUsage.Total,
			UsedPercent: diskUsage.UsedPercent,
		}
	}

	// Containers
	if dockerclient.Instance != nil {
		containers, err := dockerclient.Instance.ListContainers(ctx)
		if err == nil {
			stats.TotalContainers = len(containers)
			running := 0
			stopped := 0
			for _, cont := range containers {
				if cont.State == "running" {
					running++
				} else {
					stopped++
				}
			}
			stats.RunningContainers = running
			stats.StoppedContainers = stopped
			stats.Containers = containers
		}
	}

	c.mu.Lock()
	c.currentStats = stats
	c.mu.Unlock()
}

func (c *Collector) saveSnapshot() {
	c.mu.RLock()
	s := c.currentStats
	c.mu.RUnlock()

	_ = db.SaveMetricSnapshot(
		s.HostCPU,
		s.HostRAM.Used,
		s.HostRAM.Total,
		s.HostDisk.Used,
		s.HostDisk.Total,
		s.RunningContainers,
		s.StoppedContainers,
	)
}

func (c *Collector) GetCurrentStats() SystemStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentStats
}
