package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type MetricSnapshot struct {
	ID                int64   `json:"id"`
	Timestamp         string  `json:"timestamp"`
	CPUUsage          float64 `json:"cpu_usage"`
	RAMUsage          uint64  `json:"ram_usage"`
	RAMTotal          uint64  `json:"ram_total"`
	DiskUsage         uint64  `json:"disk_usage"`
	DiskTotal         uint64  `json:"disk_total"`
	RunningContainers int     `json:"running_containers"`
	StoppedContainers int     `json:"stopped_containers"`
}

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// Create table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS host_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		cpu_usage REAL,
		ram_usage INTEGER,
		ram_total INTEGER,
		disk_usage INTEGER,
		disk_total INTEGER,
		running_containers INTEGER,
		stopped_containers INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON host_metrics(timestamp);
	`

	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("SQLite database initialized successfully at", dbPath)
	return nil
}

func SaveMetricSnapshot(cpu float64, ramUsage, ramTotal, diskUsage, diskTotal uint64, running, stopped int) error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	query := `INSERT INTO host_metrics (cpu_usage, ram_usage, ram_total, disk_usage, disk_total, running_containers, stopped_containers) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := DB.Exec(query, cpu, ramUsage, ramTotal, diskUsage, diskTotal, running, stopped)
	if err != nil {
		return err
	}

	// Prune metrics older than 24 hours
	_, _ = DB.Exec(`DELETE FROM host_metrics WHERE timestamp < datetime('now', '-24 hours')`)
	return nil
}

func GetMetricsHistory(limit int) ([]MetricSnapshot, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	if limit <= 0 {
		limit = 60
	}

	query := `SELECT id, strftime('%Y-%m-%dT%H:%M:%SZ', timestamp), cpu_usage, ram_usage, ram_total, disk_usage, disk_total, running_containers, stopped_containers 
	          FROM host_metrics ORDER BY timestamp DESC LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []MetricSnapshot
	for rows.Next() {
		var s MetricSnapshot
		if err := rows.Scan(&s.ID, &s.Timestamp, &s.CPUUsage, &s.RAMUsage, &s.RAMTotal, &s.DiskUsage, &s.DiskTotal, &s.RunningContainers, &s.StoppedContainers); err != nil {
			continue
		}
		snapshots = append(snapshots, s)
	}

	// Reverse so oldest is first for charts
	for i, j := 0, len(snapshots)-1; i < j; i, j = i+1, j-1 {
		snapshots[i], snapshots[j] = snapshots[j], snapshots[i]
	}

	return snapshots, nil
}
