package repository_test

import (
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCheckSession(t *testing.T) {
	tests := []struct {
		name       string
		sessions   func() map[string]*repo.Session
		sessionID  string
		want       *repo.UserSessionInfo
		wantAbsent bool
	}{
		{
			name: "valid session",
			sessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID1": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 1,
							Login:  "user1",
						},
						ExpirationDate: time.Now().Add(15 * time.Minute),
					},
				}
			},
			sessionID: "sessionID1",
			want: &repo.UserSessionInfo{
				UserID: 1,
				Login:  "user1",
			},
		},
		{
			name: "expired session",
			sessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID2": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 2,
							Login:  "user2",
						},
						ExpirationDate: time.Now().Add(-5 * time.Minute),
					},
				}
			},
			sessionID:  "sessionID2",
			wantAbsent: true,
		},
		{
			name: "nonexistent session",
			sessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{}
			},
			sessionID:  "sessionID3",
			wantAbsent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionTable := &repo.SessionTable{Sessions: tt.sessions()}

			got := sessionTable.CheckSession(tt.sessionID)
			if tt.wantAbsent {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
