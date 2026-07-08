package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	_ "github.com/lib/pq"
)

func testDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Skip("TEST_DSN not set; skipping DB integration test")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
	return db
}

// seedEvent inserts a minimal user + romantic_event and returns the event id.
func seedEvent(t *testing.T, db *sql.DB) int64 {
	ctx := context.Background()
	var userID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO users(login, email, password_hash, role) VALUES($1,$2,$3,'user') RETURNING id`,
		"viewtest_"+time.Now().Format("150405.000000"), "vt"+time.Now().Format("150405.000000")+"@x.io", "x").Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	var eventID int64
	err = db.QueryRowContext(ctx,
		`INSERT INTO romantic_events(event_date, title, description, status, simp_target_id, user_id)
		 VALUES(now(),'t','d','published',NULL,$1) RETURNING id`, userID).Scan(&eventID)
	if err != nil {
		// simp_target_id is NOT NULL in schema; fall back to seeding a target.
		var targetID int64
		if e2 := db.QueryRowContext(ctx,
			`INSERT INTO simp_targets(name, description, user_id) VALUES('n','d',$1) RETURNING id`, userID).Scan(&targetID); e2 != nil {
			t.Fatalf("seed target: %v (event err: %v)", e2, err)
		}
		if e3 := db.QueryRowContext(ctx,
			`INSERT INTO romantic_events(event_date, title, description, status, simp_target_id, user_id)
			 VALUES(now(),'t','d','published',$1,$2) RETURNING id`, targetID, userID).Scan(&eventID); e3 != nil {
			t.Fatalf("seed event: %v", e3)
		}
	}
	return eventID
}

func TestEventViewRepository(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	repo := NewEventViewRepository(db)
	ctx := context.Background()
	eventID := seedEvent(t, db)

	mk := func(vid, dev string) rmodel.EventView {
		return rmodel.EventView{EventID: eventID, VisitorID: vid, Device: dev, OS: "iOS", Browser: "Safari", IP: "1.1.1.1"}
	}
	for _, v := range []rmodel.EventView{mk("dev-A", "iPhone"), mk("dev-A", "iPhone"), mk("dev-B", "Desktop")} {
		if err := repo.InsertView(ctx, v); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	last, err := repo.LastViewAt(ctx, eventID, "dev-A")
	if err != nil || last == nil {
		t.Fatalf("LastViewAt dev-A: %v %v", last, err)
	}
	if missing, _ := repo.LastViewAt(ctx, eventID, "nope"); missing != nil {
		t.Fatalf("LastViewAt missing: want nil got %v", missing)
	}

	stats, err := repo.StatsByEventID(ctx, eventID, 10)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.Views != 2 {
		t.Errorf("Views: got %d want 2", stats.Views)
	}
	if stats.Opens != 3 {
		t.Errorf("Opens: got %d want 3", stats.Opens)
	}
	if stats.LastOpenedAt == nil {
		t.Error("LastOpenedAt: want non-nil")
	}
	if len(stats.RecentOpens) != 3 {
		t.Errorf("RecentOpens: got %d want 3", len(stats.RecentOpens))
	}
}
