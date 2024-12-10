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
	lcs "github.com/yudai/golcs"
	"sort"
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

	body, err = l.sanitizeBodyJsonDetails(body)
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

func (l *Local) sanitizeBodyJsonDetails(body []byte) ([]byte, error) {
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

	mbDecoded := bDecoded.(map[string]any)
	mrecDecoded := recDecoded.(map[string]any)

	names := sortedKeys(mbDecoded) // stabilize delta order

	for _, name := range names {
		if name != "details" {
			continue
		}

		bit, ok := mbDecoded[name]
		if !ok || bit == nil {
			continue
		}

		eit, ok := bit.([]interface{})
		if !ok {
			continue
		}

		recit, ok := mrecDecoded[name]
		if !ok || recit == nil {
			continue
		}

		com, f := commutator(eit, recit.([]interface{}))
		for !f {
			eit = com

			com, f = commutator(eit, recit.([]interface{}))
		}

		mbDecoded[name] = com
	}

	body, err = json.Marshal(mbDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body:\n%+v", err)
	}

	return body, nil
}

func sortedKeys(m map[string]interface{}) (keys []string) {
	keys = make([]string, 0, len(m))

	for key, _ := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return
}

func commutator(left, right []interface{}) ([]interface{}, bool) {
	lcsPairs := lcs.New(left, right).IndexPairs()

	if len(lcsPairs) == 0 || len(lcsPairs) == len(left) {
		return left, true
	}

	// Check if all elements are equal left and right
	// This is to aovid a bug in lcs.New(left, right).IndexPairs()
	// where it returns all elements equal but not all elements of the left. e.g. [1, 2, 3] and [1, 1], [3, 3]
	var allEq bool

	for _, pair := range lcsPairs {
		if pair.Right == pair.Left {
			allEq = true

			continue
		}

		allEq = false
		break
	}

	if allEq {
		return left, true
	}

	// commutate the first pair of different not equal elements
	com := make([]interface{}, len(left))

	for _, pair := range lcsPairs {
		if pair.Right == pair.Left {
			continue
		}

		com[pair.Left] = left[pair.Right]
		com[pair.Right] = left[pair.Left]

		break
	}

	for i, v := range com {
		if v != nil {
			continue
		}

		com[i] = left[i]
	}

	return com, false
}
