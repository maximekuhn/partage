package sqlite

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func TestSaveInvitation(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteInvitationStore(db)

	tests := []struct {
		title      string
		invitation valueobject.Invitation
	}{
		{
			title:      "invitation with no update time",
			invitation: createInvitation(uuid.New(), uuid.New(), valueobject.InvitationStatusPending, time.Now(), time.Time{}),
		},
		{
			title:      "invitation with update time",
			invitation: createInvitation(uuid.New(), uuid.New(), valueobject.InvitationStatusPending, time.Now(), time.Now().Add(24*time.Hour)),
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if err := s.Save(context.TODO(), test.invitation); err != nil {
				t.Errorf("Save(): expected ok got error %v", err)
			}
		})
	}
}

func TestSaveDuplicateInvitation(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteInvitationStore(db)
	invitation := createInvitation(uuid.New(), uuid.New(), valueobject.InvitationStatusPending, time.Now(), time.Time{})

	if err := s.Save(context.TODO(), invitation); err != nil {
		panic(err)
	}

	// constraint should be on userid and groupid
	invitation.UpdatedAt = time.Now().Add(23 * time.Minute)
	invitation.Status = valueobject.InvitationStatusAccepted
	err := s.Save(context.TODO(), invitation)
	if err == nil {
		t.Error("Save(): expected err unique constraint violation got ok")
	}
}

func createInvitation(userID, groupID uuid.UUID, status valueobject.InvitationStatus, createdAt, updatedAt time.Time) valueobject.Invitation {
	userid, err := valueobject.NewUserID(userID)
	if err != nil {
		panic(err)
	}
	groupid, err := valueobject.NewGroupID(groupID)
	if err != nil {
		panic(err)
	}

	return valueobject.NewInvitation(userid, groupid, status, createdAt, updatedAt)
}
