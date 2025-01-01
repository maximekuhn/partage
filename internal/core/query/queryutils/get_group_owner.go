package queryutils

import "github.com/maximekuhn/partage/internal/core/entity"

// GetGroupOwner accepts a group and a list of associated members
// It returns the group owner if found, nil otherwise
func GetGroupOwner(group *entity.Group, members []*entity.User) *entity.User {
	for _, member := range members {
		if member.ID == group.Owner {
			return member
		}
	}
	return nil
}
