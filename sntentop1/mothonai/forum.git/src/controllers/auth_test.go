package controllers

import (
	"testing"

	"forum/src/models"
	"forum/src/utils"
)

func TestAuth(t *testing.T) {
	// Initialize database for tests
	if err := models.InitDB(":memory:"); err != nil {
		t.Skip("Database initialization failed: ", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
		want     error
	}{
		{
			name:     "tester",
			email:    "newuser@tester.er",
			password: "testASDF12!@",
			want:     nil,
		},
		{
			name:     "wrong password",
			email:    "newuser@tester.er",
			password: "wrongpassword",
			want:     models.ErrorWrongPassword,
		},
		{
			name:     "user not found",
			email:    "nonexistent@tester.er",
			password: "notfound",
			want:     models.ErrorNotRegistered,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create user for valid test case
			if tt.name == "tester" {
				hash, _ := utils.HashPassword(tt.password)
				user := models.User{
					Username: "tester",
					Email:    tt.email,
					Hash:     hash,
				}
				if err := user.Add(); err != nil {
					t.Fatalf("Failed to create test user: %v", err)
				}
			}

			err := Auth(tt.email, tt.password)
			if tt.want != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.want)
			}
			if tt.want == nil && err != nil {
				t.Errorf("expected nil, got error: %v", err)
			}
		})
	}
}
