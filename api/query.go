package api

import (
	"regexp"
	"strings"
)

// The QueryType dictates how a query should be carried out.
type QueryType int

const (
	_ QueryType = iota

	// Contain searches for equations where the equation or description
	// contains any of the words in <query>.
	// e.g. "sin cos" would match the equation "sin(A + B)"
	QueryContain

	// ContainExact searches for equations where <query> is a subset
	// of the equation or description.
	// e.g. "sin(x)" would match the equation "sin(x)^2 + cos(x)^2"
	// e.g. "sin cos" would *not* match the equation "sin(A + B)"
	QueryContainExact

	// NotContain searches for equations where <query> is not a subset
	// of the equation or description.
	// e.g. "sin(x)" would match the equation "cos(x)"
	// e.g. "sin(x)" would *not* match the equation "sin(x)^2"
	QueryNotContain

	// Regex searches for equations where <query> is interpreted as a
	// regex and has any matches in the equation or description.
	// e.g. "F = ." would match the equation "F = m a"
	// e.g. "^F = .$" would *not* match the equation "m F = m^2 a"
	QueryRegex
)

// Query runs a query through the database with given limit and query type.
func (a *API) Query(query string, qt QueryType, limit int64) ([]*Equation, error) {
	// Initialise a capacity of 32 to save time later if more than 32
	// equations are found.
	filtered := make([]*Equation, 0, 32)

	potential, err := a.FetchAllEquations()
	if err != nil {
		return filtered, err
	}

	total := int64(0)
	for _, p := range potential {
		if total >= limit {
			break
		}

		match, err := p.filter(query, qt)
		if err != nil {
			return filtered, err
		}

		if match {
			total++
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

func (e *Equation) filter(query string, qt QueryType) (bool, error) {
	l := strings.ToLower

	switch qt {
	case QueryContain:
		words := strings.Fields(query)
		for _, word := range words {
			if strings.Contains(l(e.Source), l(word)) ||
				strings.Contains(l(e.Description), l(word)) {
				return true, nil

			}
		}
		return false, nil

	case QueryContainExact:
		return strings.Contains(l(e.Source), l(query)) ||
			strings.Contains(l(e.Description), l(query)), nil

	case QueryNotContain:
		return !strings.Contains(l(e.Source), l(query)) &&
			!strings.Contains(l(e.Description), l(query)), nil

	case QueryRegex:
		matchSource, err := regexp.MatchString(query, e.Source)
		if err != nil {
			return false, ErrQueryInvalidRegex
		}

		matchDesc, err := regexp.MatchString(query, e.Description)
		if err != nil {
			return false, ErrQueryInvalidRegex
		}

		return matchSource || matchDesc, nil
	}

	return false, nil
}
