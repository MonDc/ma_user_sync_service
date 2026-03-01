package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
    UserSyncTotal     *prometheus.CounterVec
    UserSyncDuration  *prometheus.HistogramVec
    UserSyncErrors    *prometheus.CounterVec
    ActiveSyncs       prometheus.Gauge
    DatabaseLatency   *prometheus.HistogramVec
}

func NewMetrics(namespace string) *Metrics {
    return &Metrics{
        UserSyncTotal: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Name:      "user_sync_total",
                Help:      "Total number of user sync operations",
            },
            []string{"status"},
        ),
        UserSyncDuration: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Name:      "user_sync_duration_seconds",
                Help:      "Duration of user sync operations",
                Buckets:   prometheus.DefBuckets,
            },
            []string{"operation"},
        ),
        UserSyncErrors: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Namespace: namespace,
                Name:      "user_sync_errors_total",
                Help:      "Total number of user sync errors",
            },
            []string{"error_type"},
        ),
        ActiveSyncs: promauto.NewGauge(
            prometheus.GaugeOpts{
                Namespace: namespace,
                Name:      "active_syncs",
                Help:      "Number of active sync operations",
            },
        ),
        DatabaseLatency: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Namespace: namespace,
                Name:      "database_latency_seconds",
                Help:      "Database operation latency",
                Buckets:   prometheus.DefBuckets,
            },
            []string{"operation", "database"},
        ),
    }
}