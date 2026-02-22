package templates

import "github.com/jgoodhcg/mindmeld/internal/contentrating"

func contentRatingLabel(id int16) string {
	return contentrating.PoliteModeLabel(id)
}

func politeModeEnabled(id int16) bool {
	return contentrating.PoliteModeEnabled(id)
}
