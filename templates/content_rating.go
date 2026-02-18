package templates

import "github.com/jgoodhcg/mindmeld/internal/contentrating"

func contentRatingLabel(id int16) string {
	return contentrating.Label(id)
}
