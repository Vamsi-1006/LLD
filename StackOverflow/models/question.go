package models

type Question struct {
	ID      int
	Title   string
	Body    string
	Author  *User
	Tags    []*Tag
	Answers []*Answer
	Votes   *Votes
}

func NewQuestion(id int, title, body string, author *User, tags []*Tag) *Question {
	return &Question{
		ID:     id,
		Title:  title,
		Body:   body,
		Author: author,
		Tags:   tags,
		Votes:  NewVotes(),
	}
}

func (q *Question) AddAnswer(a *Answer) {
	q.Answers = append(q.Answers, a)
}

func (q *Question) HasTag(name string) bool {
	for _, t := range q.Tags {
		if t.Name == name {
			return true
		}
	}
	return false
}
