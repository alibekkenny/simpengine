package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

// seedPublishedEvent inserts a user + simp_target + a published romantic_event
// with the given public_token, and returns the seeded user id, simp_target id,
// and event id.
func seedPublishedEvent(t *testing.T, db *sql.DB, token string) (userID, targetID, eventID int64) {
	t.Helper()
	ctx := context.Background()

	suffix := time.Now().Format("150405.000000")
	if err := db.QueryRowContext(ctx,
		`INSERT INTO users(login, email, password_hash, role) VALUES($1,$2,$3,'user') RETURNING id`,
		"fbpt_"+suffix, "fbpt"+suffix+"@x.io", "x").Scan(&userID); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	if err := db.QueryRowContext(ctx,
		`INSERT INTO simp_targets(name, description, user_id) VALUES('n','d',$1) RETURNING id`,
		userID).Scan(&targetID); err != nil {
		t.Fatalf("seed target: %v", err)
	}

	if err := db.QueryRowContext(ctx,
		`INSERT INTO romantic_events(event_date, title, description, status, public_token, simp_target_id, user_id)
		 VALUES(now(),'t','d','published',$1,$2,$3) RETURNING id`,
		token, targetID, userID).Scan(&eventID); err != nil {
		t.Fatalf("seed event: %v", err)
	}

	return userID, targetID, eventID
}

func TestRomanticEventRepository_FindByPublicToken(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	token := "fbpt-token-" + time.Now().Format("150405.000000")
	userID, targetID, eventID := seedPublishedEvent(t, db, token)

	repo := NewRomanticEventRepository(db)
	event, err := repo.FindByPublicToken(context.Background(), token)
	if err != nil {
		t.Fatalf("FindByPublicToken: %v", err)
	}

	if event.ID != eventID {
		t.Errorf("ID: got %d want %d", event.ID, eventID)
	}
	if event.SimpTargetID != targetID {
		t.Errorf("SimpTargetID: got %d want %d", event.SimpTargetID, targetID)
	}
	if event.UserID != userID {
		t.Errorf("UserID: got %d want %d", event.UserID, userID)
	}
}
