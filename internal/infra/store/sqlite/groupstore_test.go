package sqlite

import (
	"context"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maximekuhn/partage/internal/core/entity"
	"github.com/maximekuhn/partage/internal/core/valueobject"
)

func TestSaveGroup(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteGroupStore(db)

	tests := []struct {
		title     string
		groupName string
		ownerID   uuid.UUID
		members   []uuid.UUID
	}{
		{
			title:     "Simple group with no members",
			groupName: "My awesome group where I am alone",
			ownerID:   uuid.New(),
			members:   []uuid.UUID{},
		},
		{
			title:     "A group with one member",
			groupName: "My awesome group with my only friend",
			ownerID:   uuid.New(),
			members:   []uuid.UUID{uuid.New()},
		},
		{
			title:     "A group with multiple members",
			groupName: "My awesome group with all my friends",
			ownerID:   uuid.New(),
			members: []uuid.UUID{
				uuid.New(),
				uuid.New(),
				uuid.New(),
				uuid.New(),
				uuid.New(),
				uuid.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			group := createGroup(test.groupName, test.ownerID, test.members)
			err := s.Save(context.TODO(), group)
			if err != nil {
				t.Errorf("Save(): expected ok got error %v", err)
			}
		})
	}
}

func createGroup(name string, owner uuid.UUID, members []uuid.UUID) *entity.Group {
	id, err := valueobject.NewGroupID(uuid.New())
	if err != nil {
		panic(err)
	}
	ownerID, err := valueobject.NewUserID(owner)
	if err != nil {
		panic(err)
	}
	membersID := make([]valueobject.UserID, 0)
	for _, member := range members {
		memberID, err := valueobject.NewUserID(member)
		if err != nil {
			panic(err)
		}
		membersID = append(membersID, memberID)
	}
	groupname, err := valueobject.NewGroupname(name)
	if err != nil {
		panic(err)
	}
	createdAt := time.Now()
	return entity.NewGroup(id, groupname, membersID, ownerID, createdAt)
}

func TestFindGroupByName(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteGroupStore(db)
	group := createGroup("My first group", uuid.New(), []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	})
	err := s.Save(context.TODO(), group)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		title         string
		groupname     string
		shouldBeFound bool
	}{
		{
			title:         "Wrong name",
			groupname:     "My second group",
			shouldBeFound: false,
		},
		{
			title:         "Correct name",
			groupname:     "My first group",
			shouldBeFound: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			groupname, err := valueobject.NewGroupname(test.groupname)
			if err != nil {
				panic(err)
			}
			g, found, err := s.FindByName(context.TODO(), groupname)
			if err != nil {
				t.Fatalf("FindByName(): should get ok got error %v", err)
			}

			if test.shouldBeFound && !found {
				t.Fatal("FindByName(): expected to find group got nothing")
			} else if !test.shouldBeFound && !found {
				// ok
				return
			}

			// sort members for comparison
			sort.Slice(g.Members, func(i, j int) bool {
				return g.Members[i].String() < g.Members[j].String()
			})

			sort.Slice(group.Members, func(i, j int) bool {
				return group.Members[i].String() < group.Members[j].String()
			})

			if !reflect.DeepEqual(g, group) {
				t.Fatalf("FindByName(): expected group %v got %v", group, g)
			}
		})
	}
}

func TestFindGroupsForUser(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteGroupStore(db)

	userID := uuid.New()

	g1 := createGroup("group - 1", userID, make([]uuid.UUID, 0))
	g2 := createGroup("group - 2", uuid.New(), []uuid.UUID{userID})
	g3 := createGroup("group - 3", uuid.New(), []uuid.UUID{uuid.New()})

	ctx := context.TODO()

	if err := s.Save(ctx, g1); err != nil {
		panic(err)
	}
	if err := s.Save(ctx, g2); err != nil {
		panic(err)
	}
	if err := s.Save(ctx, g3); err != nil {
		panic(err)
	}

	uID, err := valueobject.NewUserID(userID)
	if err != nil {
		panic(err)
	}

	groups, err := s.FindAllForUserID(ctx, uID)
	if err != nil {
		t.Fatalf("FindAllForUserID(): expected ok got error '%s'", err)
	}

	if len(groups) != 2 {
		t.Fatalf("FindAllForUserID(): expected to get 2 groups got %d", len(groups))
	}

	if !groupContains("group - 1", groups) {
		t.Fatal("FindAllForUserID(): expected to find 'group - 1' but it is not the case")
	}
	if !groupContains("group - 2", groups) {
		t.Fatal("FindAllForUserID(): expected to find 'group - 2' but it is not the case")
	}
}

func TestFindGroupByID(t *testing.T) {
	db := CreateTmpDB()
	defer db.Close()

	s := NewSQLiteGroupStore(db)

	tests := []struct {
		title              string
		group              *entity.Group
		shouldFindGroup    bool
		idToUseForQuery    uuid.UUID
		useGroupIDForQuery bool
	}{
		{
			title:              "should not find group",
			group:              createGroup("group 1", uuid.New(), []uuid.UUID{uuid.New(), uuid.New()}),
			shouldFindGroup:    false,
			idToUseForQuery:    uuid.New(),
			useGroupIDForQuery: false,
		},
		{
			title:              "should find group",
			group:              createGroup("group 2", uuid.New(), []uuid.UUID{uuid.New(), uuid.New()}),
			shouldFindGroup:    true,
			idToUseForQuery:    [16]byte{},
			useGroupIDForQuery: true,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			var groupID valueobject.GroupID
			var err error
			if test.useGroupIDForQuery {
				groupID = test.group.ID
			} else {
				groupID, err = valueobject.NewGroupID(test.idToUseForQuery)
			}
			if err != nil {
				panic(err)
			}

			if test.shouldFindGroup {
				if err = s.Save(context.TODO(), test.group); err != nil {
					panic(err)
				}
			}

			g, found, err := s.FindByID(context.TODO(), groupID)
			if err != nil {
				t.Fatalf("FindByID(): expected ok got error %v", err)
			}

			if !found && test.shouldFindGroup {
				t.Fatalf("FindByID(): expected to find group but found nothing")
			}

			if !found && !test.shouldFindGroup {
				return
			}

			// sort members for comparison
			sort.Slice(g.Members, func(i, j int) bool {
				return g.Members[i].String() < g.Members[j].String()
			})

			sort.Slice(test.group.Members, func(i, j int) bool {
				return test.group.Members[i].String() < test.group.Members[j].String()
			})

			if !reflect.DeepEqual(g, test.group) {
				t.Fatalf("FindByID(): expected group %v got %v", test.group, g)
			}

		})
	}

}

func groupContains(groupName string, groups []entity.Group) bool {
	for _, g := range groups {
		if g.Name.String() == groupName {
			return true
		}
	}
	return false
}
