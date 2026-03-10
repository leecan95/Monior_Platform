package model

import (
	"context"
	"time"
)

// DeploymentReplicaStatus stores available replica count for one deployment at a point in time.
type DeploymentReplicaStatus struct {
	CollectedAt       time.Time
	PrometheusTS      time.Time
	Namespace         string
	Deployment        string
	AvailableReplicas int64
}

// BatchInsertDeploymentReplicaStatuses stores a deployment status snapshot in batch.
func (c *PostgreDb) BatchInsertDeploymentReplicaStatuses(rows []DeploymentReplicaStatus) error {
	if len(rows) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO deployment_replicas_available (collected_at, prometheus_ts, namespace, deployment, available_replicas) VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, row := range rows {
		if _, err := stmt.ExecContext(ctx, row.CollectedAt, row.PrometheusTS, row.Namespace, row.Deployment, row.AvailableReplicas); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
