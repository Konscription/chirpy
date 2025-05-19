package auth

import "testing"

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid password",
			args:    args{password: "password123"},
			wantErr: false,
		},
		{
			name:    "empty password",
			args:    args{password: ""},
			wantErr: true,
		},
		{
			name:    "long password",
			args:    args{password: "thisisaverylongpasswordthatshouldbehashedcorrectly"},
			wantErr: false,
		},
		{
			name:    "special characters",
			args:    args{password: "!@#$%^&	&*()_+"},
			wantErr: false,
		},
		{
			name:    "numeric password",
			args:    args{password: "1234567890"},
			wantErr: false,
		},
		{
			name:    "password with spaces",
			args:    args{password: "password with spaces"},
			wantErr: false,
		},
		{
			name:    "password with unicode",
			args:    args{password: "密码123"},
			wantErr: false,
		},
		{
			name:    "password with emojis",
			args:    args{password: "password😊"},
			wantErr: false,
		},
		{
			name:    "password with control characters",
			args:    args{password: "password\n\t\r"},
			wantErr: false,
		},
		{
			name:    "password with null byte",
			args:    args{password: "password\x00"},
			wantErr: false,
		},
		{
			name:    "password with leading/trailing spaces",
			args:    args{password: "   password   "},
			wantErr: false,
		},
		{
			name:    "password with mixed case",
			args:    args{password: "Password123"},
			wantErr: false,
		},
		{
			name:    "password with repeated characters",
			args:    args{password: "aaaaaaabbbbbbcccccc"},
			wantErr: false,
		},
		{
			name:    "password with only special characters",
			args:    args{password: "!@#$%^&*()"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Validate hash matches the password if no error is expected
			if !tt.wantErr {
				err := CheckPasswordHash(got, tt.args.password)
				if err != nil {
					t.Errorf("CheckPasswordHash() error = %v", err)
				}
			}
		})
	}
}
