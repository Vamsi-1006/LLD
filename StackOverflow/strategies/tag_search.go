package strategies

import (
	"stackoverflow/models"
	"strings"
)

type TagSearch struct{}

func (t TagSearch) Search(questions []*models.Question, query string) []*models.Question {
	query = strings.ToLower(query)

	var result []*models.Question

	for _, q := range questions {
		if q.HasTag(query) {
			result = append(result, q)
		}
	}

	return result
}
