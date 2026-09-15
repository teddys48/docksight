package dockerclient

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	ID         string            `json:"id"`
	Names      []string          `json:"names"`
	Image      string            `json:"image"`
	ImageID    string            `json:"image_id"`
	Command    string            `json:"command"`
	Created    int64             `json:"created"`
	State      string            `json:"state"`
	Status     string            `json:"status"`
	Ports      []PortMapping     `json:"ports"`
	Mounts     []MountPoint      `json:"mounts"`
	Labels     map[string]string `json:"labels"`
	CPUUsage   float64           `json:"cpu_usage"`
	RAMUsage   uint64            `json:"ram_usage"`
	RAMLimit   uint64            `json:"ram_limit"`
	RAMPercent float64           `json:"ram_percent"`
	NetInput   uint64            `json:"net_input"`
	NetOutput  uint64            `json:"net_output"`
	BlockRead  uint64            `json:"block_read"`
	BlockWrite uint64            `json:"block_write"`
}

type PortMapping struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint16 `json:"private_port"`
	PublicPort  uint16 `json:"public_port,omitempty"`
	Type        string `json:"type"`
}

type MountPoint struct {
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

type ImageInfo struct {
	ID          string   `json:"id"`
	RepoTags    []string `json:"repo_tags"`
	RepoDigests []string `json:"repo_digests"`
	Created     int64    `json:"created"`
	Size        int64    `json:"size"`
	VirtualSize int64    `json:"virtual_size"`
	IsDangling  bool     `json:"is_dangling"`
}

type VolumeInfo struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Scope      string            `json:"scope"`
	Mountpoint string            `json:"mountpoint"`
	Created    string            `json:"created"`
	Options    map[string]string `json:"options"`
}

type DockerClient struct {
	cli *client.Client
}

var Instance *DockerClient

func InitDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Ping docker daemon
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		log.Printf("Warning: Unable to ping docker daemon: %v", err)
	} else {
		log.Println("Successfully connected to Docker Daemon")
	}

	Instance = &DockerClient{cli: cli}
	return Instance, nil
}

func (dc *DockerClient) GetClient() *client.Client {
	return dc.cli
}

func (dc *DockerClient) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := dc.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		ports := make([]PortMapping, 0, len(c.Ports))
		for _, p := range c.Ports {
			ports = append(ports, PortMapping{
				IP:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}

		mounts := make([]MountPoint, 0, len(c.Mounts))
		for _, m := range c.Mounts {
			mounts = append(mounts, MountPoint{
				Type:        string(m.Type),
				Name:        m.Name,
				Source:      m.Source,
				Destination: m.Destination,
				Mode:        m.Mode,
				RW:          m.RW,
			})
		}

		info := ContainerInfo{
			ID:      c.ID,
			Names:   c.Names,
			Image:   c.Image,
			ImageID: c.ImageID,
			Command: c.Command,
			Created: c.Created,
			State:   c.State,
			Status:  c.Status,
			Ports:   ports,
			Mounts:  mounts,
			Labels:  c.Labels,
		}

		// Fetch quick stats if container is running
		if c.State == "running" {
			stats, err := dc.getContainerStatsQuick(ctx, c.ID)
			if err == nil {
				info.CPUUsage = stats.CPUUsage
				info.RAMUsage = stats.RAMUsage
				info.RAMLimit = stats.RAMLimit
				info.RAMPercent = stats.RAMPercent
				info.NetInput = stats.NetInput
				info.NetOutput = stats.NetOutput
				info.BlockRead = stats.BlockRead
				info.BlockWrite = stats.BlockWrite
			}
		}

		result = append(result, info)
	}

	return result, nil
}

type QuickStats struct {
	CPUUsage   float64
	RAMUsage   uint64
	RAMLimit   uint64
	RAMPercent float64
	NetInput   uint64
	NetOutput  uint64
	BlockRead  uint64
	BlockWrite uint64
}

func (dc *DockerClient) getContainerStatsQuick(ctx context.Context, id string) (QuickStats, error) {
	var qs QuickStats
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := dc.cli.ContainerStatsOneShot(ctxTimeout, id)
	if err != nil {
		return qs, err
	}
	defer resp.Body.Close()

	var stats types.StatsJSON
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return qs, err
	}

	// Calculate CPU Percent
	var cpuPercent float64
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
	onlineCPUs := float64(stats.CPUStats.OnlineCPUs)
	if onlineCPUs == 0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	if onlineCPUs == 0 {
		onlineCPUs = 1
	}

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	// Memory
	ramUsage := stats.MemoryStats.Usage
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		if ramUsage > cache {
			ramUsage -= cache
		}
	}
	ramLimit := stats.MemoryStats.Limit
	var ramPercent float64
	if ramLimit > 0 {
		ramPercent = (float64(ramUsage) / float64(ramLimit)) * 100.0
	}

	// Network I/O
	var netInput, netOutput uint64
	for _, net := range stats.Networks {
		netInput += net.RxBytes
		netOutput += net.TxBytes
	}

	// Block I/O
	var blockRead, blockWrite uint64
	for _, bio := range stats.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(bio.Op) {
		case "read":
			blockRead += bio.Value
		case "write":
			blockWrite += bio.Value
		}
	}

	qs.CPUUsage = cpuPercent
	qs.RAMUsage = ramUsage
	qs.RAMLimit = ramLimit
	qs.RAMPercent = ramPercent
	qs.NetInput = netInput
	qs.NetOutput = netOutput
	qs.BlockRead = blockRead
	qs.BlockWrite = blockWrite

	return qs, nil
}

func (dc *DockerClient) StartContainer(ctx context.Context, id string) error {
	return dc.cli.ContainerStart(ctx, id, container.StartOptions{})
}

func (dc *DockerClient) StopContainer(ctx context.Context, id string) error {
	timeout := 10
	return dc.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (dc *DockerClient) RestartContainer(ctx context.Context, id string) error {
	timeout := 10
	return dc.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (dc *DockerClient) RemoveContainer(ctx context.Context, id string, force bool) error {
	return dc.cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: force})
}

func (dc *DockerClient) InspectContainer(ctx context.Context, id string) (interface{}, error) {
	inspectData, err := dc.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, err
	}
	return inspectData, nil
}

func (dc *DockerClient) ListImages(ctx context.Context) ([]ImageInfo, error) {
	images, err := dc.cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	result := make([]ImageInfo, 0, len(images))
	for _, img := range images {
		isDangling := false
		if len(img.RepoTags) == 0 || (len(img.RepoTags) == 1 && img.RepoTags[0] == "<none>:<none>") {
			isDangling = true
		}

		result = append(result, ImageInfo{
			ID:          img.ID,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Created:     img.Created,
			Size:        img.Size,
			VirtualSize: img.VirtualSize,
			IsDangling:  isDangling,
		})
	}

	return result, nil
}

func (dc *DockerClient) ListVolumes(ctx context.Context) ([]VolumeInfo, error) {
	vols, err := dc.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]VolumeInfo, 0, len(vols.Volumes))
	for _, v := range vols.Volumes {
		result = append(result, VolumeInfo{
			Name:       v.Name,
			Driver:     v.Driver,
			Scope:      v.Scope,
			Mountpoint: v.Mountpoint,
			Created:    v.CreatedAt,
			Options:    v.Options,
		})
	}

	return result, nil
}

// StreamContainerLogs streams stdout/stderr lines with docker multiplex header stripped
func (dc *DockerClient) StreamLogs(ctx context.Context, id string, tail string, stdout, stderr, timestamps, follow bool, outChan chan<- string) error {
	opts := container.LogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Timestamps: timestamps,
		Follow:     follow,
		Tail:       tail,
	}

	reader, err := dc.cli.ContainerLogs(ctx, id, opts)
	if err != nil {
		return err
	}
	defer reader.Close()

	bufReader := bufio.NewReader(reader)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			header := make([]byte, 8)
			_, err := io.ReadFull(bufReader, header)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				line, readErr := bufReader.ReadString('\n')
				if readErr == nil && len(line) > 0 {
					outChan <- line
					continue
				}
				return err
			}

			frameSize := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
			if frameSize <= 0 {
				continue
			}

			payload := make([]byte, frameSize)
			_, err = io.ReadFull(bufReader, payload)
			if err != nil {
				return err
			}

			logPrefix := ""
			if header[0] == 2 {
				logPrefix = "[STDERR] "
			}

			lines := strings.Split(string(payload), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					outChan <- logPrefix + line
				}
			}
		}
	}
}
