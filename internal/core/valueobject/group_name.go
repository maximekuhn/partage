package valueobject

type Groupname struct {
	name string
}

func NewGroupname(name string) (Groupname, error) {
	gn := Groupname{name}
	return gn, nil
}

func (g Groupname) String() string {
	return g.name
}
