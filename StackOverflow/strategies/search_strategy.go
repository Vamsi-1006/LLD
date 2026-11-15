package strategies

import "stackoverflow/models"

type SearchStrategy interface {
	Search(questions []*models.Question, query string) []*models.Question
}
