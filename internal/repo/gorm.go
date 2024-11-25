package repo

import (
	"errors"
	"fmt"
	"strings"
)

type OrderBy struct {
	Field     string
	Direction string
}

var allowedFields = map[string]string{
	"group":        "group_name",
	"song":         "song",
	"release_date": "release_date",
}

var allowedDirections = map[string]bool{
	"asc":  true,
	"desc": true,
}

func parseOrderBy(orderByParam string) ([]OrderBy, error) {
	var orderBys []OrderBy

	if orderByParam == "" {
		orderBys = append(orderBys, OrderBy{
			Field:     "id",
			Direction: "asc",
		})
		return orderBys, nil
	}

	pairs := strings.Split(orderByParam, ",")
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, errors.New("invalid_orderBy_format")
		}

		field := parts[0]
		direction := parts[1]

		mappedField, fieldExists := allowedFields[field]
		if !fieldExists {
			return nil, fmt.Errorf("invalid_field_for_orderBy: %s", field)
		}
		if !allowedDirections[direction] {
			return nil, fmt.Errorf("invalid_direction_for_orderBy: %s", direction)
		}

		orderBys = append(orderBys, OrderBy{
			Field:     mappedField, // Use the mapped field name
			Direction: direction,
		})
	}
	return orderBys, nil
}
