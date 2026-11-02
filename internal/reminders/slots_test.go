package reminders

import (
	"testing"
	"time"
)

func TestSlots(t *testing.T) {
	loc := time.FixedZone("IST", 5*3600+30*60)
	scheduled := time.Date(2026, 11, 3, 10, 0, 0, 0, loc)
	now := time.Date(2026, 11, 1, 9, 0, 0, 0, loc)
	got := Slots(scheduled, now)
	if len(got) != 2 || got[0].Kind != "eve" || got[1].Kind != "soon" {
		t.Fatalf("slots = %#v", got)
	}
	wantEve := time.Date(2026, 11, 2, 18, 0, 0, 0, loc)
	if !got[0].At.Equal(wantEve) {
		t.Fatalf("eve = %s, want %s", got[0].At, wantEve)
	}
	if !got[1].At.Equal(scheduled.Add(-2 * time.Hour)) {
		t.Fatalf("soon = %s", got[1].At)
	}

	late := time.Date(2026, 11, 3, 9, 30, 0, 0, loc)
	if len(Slots(scheduled, late)) != 0 {
		t.Fatal("expected no reminders inside the lead time")
	}
	if len(Slots(scheduled, scheduled.Add(time.Minute))) != 0 {
		t.Fatal("expected no reminders for a past hearing")
	}
}
