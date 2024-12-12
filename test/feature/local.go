package feature

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bool64/httpdog"
	"github.com/bool64/shared"
	"github.com/cucumber/godog"
	"github.com/swaggest/assertjson/json5"
	"reflect"
)

// Local is step-driven HTTP client for application local HTTP service.
type Local struct {
	*httpdog.Local
}

// NewLocal creates an instance of step-driven HTTP client.
func NewLocal(baseURL string) *Local {
	return &Local{
		Local: httpdog.NewLocal(baseURL),
	}
}

func (l *Local) RegisterSteps(s *godog.ScenarioContext) {
	l.Local.RegisterSteps(s)

	s.Step(`^I should have response with body like$`, l.iShouldHaveResponseWithBody)
}

func (l *Local) iShouldHaveResponseWithBody(bodyDoc *godog.DocString) error {
	body, err := loadBody([]byte(bodyDoc.Content), l.JSONComparer.Vars)
	if err != nil {
		return err
	}

	body, err = l.alignBodyIfPossible(body)
	if err != nil {
		return err
	}

	return l.ExpectResponseBody(body)
}

func loadBody(body []byte, vars *shared.Vars) ([]byte, error) {
	var err error

	if json5.Valid(body) {
		if body, err = json5.Downgrade(body); err != nil {
			return nil, fmt.Errorf("failed to downgrade JSON5 to JSON: %w", err)
		}
	}

	if vars != nil {
		for k, v := range vars.GetAll() {
			jv, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal var %s (%v): %w", k, v, err)
			}

			body = bytes.ReplaceAll(body, []byte(`"`+k+`"`), jv)
		}
	}

	return body, nil
}

func (l *Local) alignBodyIfPossible(body []byte) ([]byte, error) {
	if len(body) == 0 {
		if len(body) == 0 {
			return nil, nil
		}

		return nil, errors.New("received empty body")
	}

	if body != nil && !json5.Valid(body) {
		return nil, errors.New("received invalid JSON5 body")
	}

	received := l.Client.Details().RespBody

	// sort both arrays to make comparison predictable
	var bDecoded, recDecoded any

	err := json.Unmarshal(body, &bDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body:\n%+v", err)
	}

	err = json.Unmarshal(received, &recDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal received:\n%+v", err)
	}

	bDecoded, _ = AlignLeftRightIfPossible(bDecoded, recDecoded)

	body, err = json.Marshal(bDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body:\n%+v", err)
	}

	return body, nil
}

func AlignLeftRightIfPossible(left, right any) (any, bool) {
	if left == nil && right == nil {
		return left, true
	}

	if left == nil || right == nil {
		return left, false
	}

	switch l := left.(type) {
	case []any:
		r, ok := right.([]any)
		if !ok {
			return left, false
		}

		return alignLeftRightSlicePossible(l, r)

	case map[string]any:
		r, ok := right.(map[string]any)
		if !ok {
			return left, false
		}

		return alignLeftRightMapPossible(l, r)

	default:
		// Check if type are equals
		if reflect.TypeOf(left) != reflect.TypeOf(right) {
			return left, false
		}

		// Skip if left is a placeholder for ignoring diff <ignore-diff>
		if left == "<ignore-diff>" {
			println("AlignLeftRightIfPossible.return true")
			return right, true
		}

		if left == right {
			return left, true
		}

		return left, false
	}
}

func alignLeftRightSlicePossible(left, right []any) ([]any, bool) {
	if len(left) != len(right) {
		return left, false
	}

	var (
		// comm contains the common elements between left and right.
		comm = make([]any, len(left))
		// notComm contains the elements that are not common between left and right, but from left.
		notComm = make([]any, 0, len(left))
	)

	for i, l := range left {
		if reflect.DeepEqual(l, right[i]) {
			comm[i] = l

			continue
		}

		aligned, ok := AlignLeftRightIfPossible(l, right[i])
		if ok {
			comm[i] = aligned

			continue
		}

		notComm = append(notComm, l)
	}

	// Check if all elements are equal left and right
	if len(notComm) == 0 {
		return comm, true
	}

	var reo bool

	for _, n := range notComm {
		reo = false

		for i, r := range right {
			if reflect.DeepEqual(n, r) {
				comm[i] = n

				reo = true

				break
			}
		}

		if !reo {
			return left, false
		}
	}

	return comm, true
}

func alignLeftRightMapPossible(l map[string]any, r map[string]any) (any, bool) {
	if len(l) != len(r) {
		return l, false
	}

	comm := make(map[string]any, len(l))

	for lk, lv := range l {
		if rv, ok := r[lk]; ok {
			if reflect.DeepEqual(lv, rv) {
				comm[lk] = lv

				continue
			}

			aligned, ok := AlignLeftRightIfPossible(lv, rv)
			if ok {
				comm[lk] = aligned

				continue
			}
		}
	}

	if len(comm) == len(l) {
		return comm, true
	}

	return l, false
}
