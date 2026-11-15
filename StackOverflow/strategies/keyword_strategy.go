package strategies

import (
	"stackoverflow/models"
	"strings"
)

type KeyWordSearch struct{}

func (k KeyWordSearch) Search(questions []*models.Question, query string) []*models.Question {
	query = strings.ToLower(query)

	var result []*models.Question

	for _, q := range questions {
		text := strings.ToLower(q.Title + " " + q.Body)
		if strings.Contains(text, query) {
			result = append(result, q)
		}
	}

	return result
}
