package repository_test

import (
	repo "github.com/go-park-mail-ru/2024_1_CyberHedgehogs/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteSession(t *testing.T) {
	tests := []struct {
		name              string
		initialSessions   func() map[string]*repo.Session
		sessionIDToDelete string
		expectedError     bool
		expectedSessions  func() map[string]*repo.Session
	}{
		{
			name: "delete existing session",
			initialSessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID1": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 1,
							Login:  "user1",
						},
					},
					"sessionID2": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 2,
							Login:  "user2",
						},
					},
				}
			},
			sessionIDToDelete: "sessionID1",
			expectedError:     false,
			expectedSessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID2": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 2,
							Login:  "user2",
						},
					},
				}
			},
		},
		{
			name: "delete non-existing session",
			initialSessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID1": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 1,
							Login:  "user1",
						},
					},
				}
			},
			sessionIDToDelete: "sessionID3",
			expectedError:     false,
			expectedSessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID1": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 1,
							Login:  "user1",
						},
					},
				}
			},
		},
		{
			name: "delete last session",
			initialSessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{
					"sessionID1": {
						UserInfo: &repo.UserSessionInfo{
							UserID: 1,
							Login:  "user1",
						},
					},
				}
			},
			sessionIDToDelete: "sessionID1",
			expectedError:     false,
			expectedSessions: func() map[string]*repo.Session {
				return map[string]*repo.Session{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionTable := &repo.SessionTable{Sessions: tt.initialSessions()}

			err := sessionTable.DeleteSession(tt.sessionIDToDelete)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedSessions(), sessionTable.Sessions)
		})
	}
}
