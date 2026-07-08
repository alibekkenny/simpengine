package postgres

import (
	"context"
	"database/sql"
	"time"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/shared/db"
)

type EventViewRepository struct {
	db *sql.DB
}

func NewEventViewRepository(db *sql.DB) EventViewRepository {
	return EventViewRepository{db: db}
}

func (r EventViewRepository) InsertView(ctx context.Context, v rmodel.EventView) error {
	stmt := `INSERT INTO event_views(event_id, visitor_id, device, os, browser, ip)
	VALUES($1, $2, $3, $4, $5, $6)`
	if _, err := r.db.ExecContext(ctx, stmt, v.EventID, v.VisitorID, v.Device, v.OS, v.Browser, v.IP); err != nil {
		return db.MapPQError(err)
	}
	return nil
}

func (r EventViewRepository) LastViewAt(ctx context.Context, eventID int64, visitorID string) (*time.Time, error) {
	stmt := `SELECT MAX(created_at) FROM event_views WHERE event_id = $1 AND visitor_id = $2`
	var t sql.NullTime
	if err := r.db.QueryRowContext(ctx, stmt, eventID, visitorID).Scan(&t); err != nil {
		return nil, db.MapPQError(err)
	}
	if !t.Valid {
		return nil, nil
	}
	return &t.Time, nil
}

func (r EventViewRepository) StatsByEventID(ctx context.Context, eventID int64, recentLimit int) (rmodel.EventViewStats, error) {
	stats := rmodel.EventViewStats{RecentOpens: []rmodel.EventViewSummary{}}
	var last sql.NullTime
	head := `SELECT COUNT(DISTINCT visitor_id), COUNT(*), MAX(created_at) FROM event_views WHERE event_id = $1`
	if err := r.db.QueryRowContext(ctx, head, eventID).Scan(&stats.Views, &stats.Opens, &last); err != nil {
		return stats, db.MapPQError(err)
	}
	if last.Valid {
		stats.LastOpenedAt = &last.Time
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT device, os, browser, created_at FROM event_views WHERE event_id = $1 ORDER BY created_at DESC LIMIT $2`,
		eventID, recentLimit)
	if err != nil {
		return stats, db.MapPQError(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s rmodel.EventViewSummary
		if err := rows.Scan(&s.Device, &s.OS, &s.Browser, &s.OpenedAt); err != nil {
			return stats, err
		}
		stats.RecentOpens = append(stats.RecentOpens, s)
	}
	return stats, rows.Err()
}
