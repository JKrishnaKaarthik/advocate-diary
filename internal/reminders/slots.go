package reminders

import "time"

type Slot struct {
	Kind string
	At   time.Time
}

func Slots(scheduled, now time.Time) []Slot {
	if !scheduled.After(now) {
		return nil
	}
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*3600+30*60)
	}
	local := scheduled.In(loc)
	eve := time.Date(local.Year(), local.Month(), local.Day(), 18, 0, 0, 0, loc).AddDate(0, 0, -1)
	soon := scheduled.Add(-2 * time.Hour)
	var out []Slot
	if eve.After(now) {
		out = append(out, Slot{Kind: "eve", At: eve.UTC()})
	}
	if soon.After(now) {
		out = append(out, Slot{Kind: "soon", At: soon.UTC()})
	}
	return out
}
