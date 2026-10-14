package main

import "fmt"

type Appointment struct {
	Status    string
	PatientID string
}

func notificationFor(a Appointment) string {
	switch a.Status {
	case "confirmed":
		return "send_confirmation"
	case "cancelled":
		return "send_cancellation"
	case "waitlisted":
		return "send_waitlist_update"
	default:
		return "hold"
	}
}
func reportAppointment(a Appointment) error {
	d := notificationFor(a)
	return infrai.metrics.report(map[string]any{"type": "counter", "name": "appointments.workflow", "value": 1, "tags": map[string]string{"status": a.Status, "notification": d}})
}
func main() {
	a := Appointment{Status: "confirmed", PatientID: "patient-42"}
	if err := reportAppointment(a); err != nil {
		fmt.Println(err)
		return
	}
	usage, err := usageTimeseries("window=24h")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("appointment=%s notification=%s platform_usage=%s\n", a.Status, notificationFor(a), usage)
}
