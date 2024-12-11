package feature_test

import (
	"encoding/json"
	"github.com/dohernandez/go-grpc-service/test/feature"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type testCase struct {
	name  string
	left  []byte
	right []byte
}

type testCaseFunc = func(*testing.T) testCase

func TestAlignLeftRightIfPossible(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		alignLeftRightDetails(t),
		alignLeftRightMapSlice(t),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decLeft, decRight any

			err := json.Unmarshal(tt.left, &decLeft)
			require.NoError(t, err)

			err = json.Unmarshal(tt.right, &decRight)
			require.NoError(t, err)

			got, ok := feature.AlignLeftRightIfPossible(decLeft, decRight)
			require.True(t, ok)
			require.True(t, reflect.DeepEqual(got, decRight))
		})
	}
}

func alignLeftRightDetails(t *testing.T) testCase {
	t.Helper()

	return testCase{
		name: "align left right details",
		left: []byte(`{
			  "code": 400,
			  "message": "Bad Request",
			  "error": "<ignore-diff>",
			  "details": [
				  {"field": "first_name", "description": "must not be empty"},
				  {"field": "last_name", "description": "must not be empty"},
				  {"field": "password_hash", "description": "invalid hash"},
				  {"field": "email", "description": "must be a valid email"},
				  {"field": "country", "description": "must have 2 characters"}
			  ]
		}`),
		right: []byte(`{
			  "code": 400,
			  "error": "<ignore-diff>",
			  "message": "Bad Request",
			  "details": [
				  {"field": "password_hash", "description": "invalid hash"},
				  {"field": "last_name", "description": "must not be empty"},
				  {"field": "email", "description": "must be a valid email"},
				  {"field": "country", "description": "must have 2 characters"},
				  {"field": "first_name", "description": "must not be empty"}
			  ]
			}`),
	}
}

func alignLeftRightMapSlice(t *testing.T) testCase {
	t.Helper()

	return testCase{
		name: "align left right map slice",
		left: []byte(`{
			  "users":[
				{
				  "id": "26ef0140-c436-4838-a271-32652c72f6f2",
				  "first_name": "Alice",
				  "last_name": "Bob",
				  "nickname": "",
				  "password_hash": "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
				  "email": "alice@bob.com",
				  "country": "UK"
				},{
				  "id": "29d7fe1d-6d03-4c52-9880-d39788f9c227",
				  "first_name": "Lina",
				  "last_name": "Lowe",
				  "nickname": "magna",
				  "password_hash": "41eeaa061fa11f084957d4522cb4b408dbe4b16f446c513883d8c81e66da33f6",
				  "email": "linalowe@beadzza.com",
				  "country": "UK"
				},{
				  "id": "f1ec4c49-2166-45d2-988f-cb632bd380f9",
				  "first_name": "Roman",
				  "last_name": "Keith",
				  "nickname": "dolor",
				  "password_hash": "80e967e6c166120fc14badb021298fdb9ae5f20224d4c6c416d9898cfcc3b7e7",
				  "email": "romankeith@beadzza.com",
				  "country": "UK"
				}
			  ],
			  "nextPageToken":""
    	}`),
		right: []byte(`{
			  "users": [
				{
				  "id": "26ef0140-c436-4838-a271-32652c72f6f2",
				  "first_name": "Alice",
				  "last_name": "Bob",
				  "nickname": "",
				  "password_hash": "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
				  "email": "alice@bob.com",
				  "country": "UK"
				},
				{
				  "id": "29d7fe1d-6d03-4c52-9880-d39788f9c227",
				  "first_name": "Lina",
				  "last_name": "Lowe",
				  "nickname": "magna",
				  "password_hash": "41eeaa061fa11f084957d4522cb4b408dbe4b16f446c513883d8c81e66da33f6",
				  "email": "linalowe@beadzza.com",
				  "country": "UK"
				},
				{
				  "id": "f1ec4c49-2166-45d2-988f-cb632bd380f9",
				  "first_name": "Roman",
				  "last_name": "Keith",
				  "nickname": "dolor",
				  "password_hash": "80e967e6c166120fc14badb021298fdb9ae5f20224d4c6c416d9898cfcc3b7e7",
				  "email": "romankeith@beadzza.com",
				  "country": "UK"
				}
			  ],
			  "nextPageToken": ""
		}`),
	}

}
