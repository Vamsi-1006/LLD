package models

type Answer struct {
	ID         int
	Body       string
	Author     *User
	QuestionID int
	Votes      *Votes
}

func NewAnswer(id int, body string, author *User, questionID int) *Answer {
	return &Answer{
		ID:         id,
		Body:       body,
		Author:     author,
		QuestionID: questionID,
		Votes:      NewVotes(),
	}
}
