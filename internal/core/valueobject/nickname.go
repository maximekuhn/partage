package valueobject

type Nickname struct {
	nickname string
}

func NewNickname(nickname string) (Nickname, error) {
	n := Nickname{nickname}
	return n, nil
}

func (n Nickname) String() string {
	return n.nickname
}
