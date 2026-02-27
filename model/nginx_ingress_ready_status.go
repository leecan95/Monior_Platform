package model

import (
	"context"
	"time"
)

// NginxIngressReadyStatus stores nginx-ingress ready state snapshots for downtime statistics.
type NginxIngressReadyStatus struct {
	CollectedAt  time.Time
	PrometheusTS *time.Time
	Namespace    string
	ReadyStatus  int16
}

// InsertNginxIngressReadyStatus inserts a single nginx-ingress ready state row.
func (c *PostgreDb) InsertNginxIngressReadyStatus(row NginxIngressReadyStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	_, err := c.db.ExecContext(ctx,
		`INSERT INTO nginx_ingress_ready_status (collected_at, prometheus_ts, namespace, ready_status) VALUES ($1,$2,$3,$4)`,
		row.CollectedAt,
		row.PrometheusTS,
		row.Namespace,
		row.ReadyStatus,
	)

	return err
}
