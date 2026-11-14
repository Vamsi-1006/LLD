package models

type User struct {
	ID         int
	Name       string
	Reputation int
}

func NewUser(id int, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

// AddReputation touches only the user's own field.
func (u *User) AddReputation(points int) {
	u.Reputation += points
}
