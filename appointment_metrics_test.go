package main

import "testing"

func TestNotificationDecision(t *testing.T) {
	cases := []struct{ name, status, want string }{{"confirmed", "confirmed", "send_confirmation"}, {"cancelled", "cancelled", "send_cancellation"}, {"waitlist", "waitlisted", "send_waitlist_update"}, {"draft", "draft", "hold"}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := notificationFor(Appointment{Status: tc.status}); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
