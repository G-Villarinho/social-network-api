package utils

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

func GenerateSuggestions(username string, count int) []string {
	suggestions := make([]string, 0, count)
	currentYear := time.Now().Year()

	for len(suggestions) < count {
		switch len(suggestions) {
		case 0:
			suggestions = append(suggestions, fmt.Sprintf("%s%d", username, rand.Intn(1000)))
		case 1:
			suggestions = append(suggestions, fmt.Sprintf("%s_%d", username, currentYear))
		case 2:
			suggestions = append(suggestions, fmt.Sprintf("%s_official", username))
		case 3:
			suggestions = append(suggestions, fmt.Sprintf("the_real_%s", username))
		case 4:
			suggestions = append(suggestions, fmt.Sprintf("%s_%d", username, rand.Intn(100)))
		default:
			suggestions = append(suggestions, fmt.Sprintf("%s_%d", username, rand.Intn(1000)))
		}
	}

	return suggestions
}
