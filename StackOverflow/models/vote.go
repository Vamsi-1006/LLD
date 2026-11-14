package models

type VoteType int

const (
	Upvote VoteType = iota
	Downvote
)

// Votes. This Votes struct is present on Question struct and Answer struct.
type Votes struct {
	byUser map[int]VoteType // userID -> their vote
}

func NewVotes() *Votes {
	return &Votes{
		byUser: make(map[int]VoteType),
	}
}

func (v *Votes) HasVoted(userID int) bool {
	_, ok := v.byUser[userID]
	return ok
}

func (v *Votes) AddVote(userID int, vt VoteType) {
	v.byUser[userID] = vt
}

// Score = number of upvotes - number of downvotes
func (v *Votes) Score() int {
	score := 0

	for _, vt := range v.byUser {
		if vt == Upvote {
			score++
		} else if vt == Downvote {
			score--
		}
	}
	return score
}
