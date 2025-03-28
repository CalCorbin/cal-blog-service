package models

import (
	"testing"
)

func TestUser_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "SecurePassword123!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{}
			err := u.HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the password was actually hashed
				if u.Password == tt.password {
					t.Errorf("HashPassword() failed to hash password, still matches original")
				}

				// A proper hash should be stored
				if len(u.Password) == 0 {
					t.Errorf("HashPassword() didn't store a hash")
				}
			}
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	originalPassword := "SecurePassword123!"

	// Create user with hashed password
	u := &User{}
	err := u.HashPassword(originalPassword)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}

	// Test cases
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Correct password",
			password: originalPassword,
			want:     true,
		},
		{
			name:     "Incorrect password",
			password: "WrongPassword",
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.CheckPassword(tt.password)
			if got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
