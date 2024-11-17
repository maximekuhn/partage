package valueobject

import "time"

type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "Pending"
	InvitationStatusAccepted InvitationStatus = "Accepted"
	InvitationStatusRejected InvitationStatus = "Rejected"
)

type Invitation struct {
	UserID    UserID
	GroupID   GroupID
	Status    InvitationStatus
	CreatedAt time.Time
	UpdatedAt time.Time // accepted or rejected
}

func NewInvitation(userID UserID, groupID GroupID, status InvitationStatus, createdAt, updatedAt time.Time) Invitation {
	return Invitation{userID, groupID, status, createdAt.UTC(), updatedAt.UTC()}
}
