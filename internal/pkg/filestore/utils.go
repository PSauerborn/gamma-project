package filestore

type SearchType int

const (
	CompleteMatch SearchType = iota
	PartialMatch
	NoMatch
)

func (t SearchType) String() string {
	return [...]string{"CompleteMatch", "PartialMatch", "NoMatch"}[t]
}

// function used to search a given map with a given list of terms
func MapMatchesTerms(items, terms map[string]interface{},
	t SearchType) bool {
	matches := []bool{}
	// iterate over search terms and compare key:val pairs
	for key, value := range terms {
		if val, ok := items[key]; ok {
			// if value is present and matches, add true to
			if val == value {
				matches = append(matches, true)
				// else add false
			} else {
				matches = append(matches, false)
			}
			// if key is not present in items, add false
		} else {
			matches = append(matches, false)
		}
	}

	switch t {
	// return true only if all terms are matched
	case CompleteMatch:
		for _, match := range matches {
			if !match {
				return false
			}
		}
		return true
	// return true if one or more terms matched
	case PartialMatch:
		for _, match := range matches {
			if match {
				return true
			}
		}
		return false
	// return true if no terms matched
	case NoMatch:
		for _, match := range matches {
			if match {
				return false
			}
		}
		return true
	}
	return false
}
