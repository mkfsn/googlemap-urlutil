package urlutil

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var (
	ErrInvalidURL = errors.New("invalid embed url")
)

func ParseEmbed(rawURL string) (coordinate Coordinate, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	parts := strings.Split(u.Query().Get("pb"), "!")
	parts = filter(parts, func(s string) bool { return len(s) > 0 })

	root, err := decode(parts)
	if err != nil {
		return Coordinate{}, err
	}

	defer func() {
		if caughtError := recover(); caughtError != nil {
			err = fmt.Errorf("%s: %w", caughtError, ErrInvalidURL)
		}
	}()

	return Coordinate{
		Lat: root[0].([]interface{})[0].([]interface{})[0].([]interface{})[2].(float64),
		Lng: root[0].([]interface{})[0].([]interface{})[0].([]interface{})[1].(float64),
	}, nil
}

func filter(elements []string, fn func(s string) bool) []string {
	results := make([]string, 0, len(elements))
	for _, e := range elements {
		if !fn(e) {
			continue
		}
		results = append(results, e)
	}
	return results
}

// Following https://stackoverflow.com/questions/47017387/decoding-the-google-maps-embedded-parameters
func decode(parts []string) ([]interface{}, error) {
	root := make([]interface{}, 0)

	for i := 0; i < len(parts); {
		part := parts[i]

		kind := part[1:2]
		value := part[2:]

		switch kind {
		case "m":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to decode kind: %s: %w", kind, ErrInvalidURL)
			}

			children, err := decode(parts[i+1 : i+1+v])
			if err != nil {
				return nil, err
			}
			root = append(root, children)
			i += v

		case "b":
			root = append(root, value == "1")
			i++
		case "d", "f":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to decode kind: %s: %w", kind, ErrInvalidURL)
			}
			root = append(root, v)
			i++
		case "i", "u", "e":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to decode kind: %s: %w", kind, ErrInvalidURL)
			}
			root = append(root, v)
			i++
		default:
			root = append(root, value)
			i++
		}
	}

	return root, nil
}
