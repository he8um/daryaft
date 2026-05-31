package cmd

import "testing"

func TestDoctorCommandExists(t *testing.T) {
	command, _, err := rootCmd.Find([]string{"doctor"})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if command == nil || command.Name() != "doctor" {
		t.Fatal("doctor command not found")
	}
}
