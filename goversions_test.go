package main

import (
	"reflect"
	"testing"
)

func TestParseVersion(t *testing.T) {

	cases := []struct {
		in     string
		expect ParsedVersion
	}{
		{
			in: "go1.3rc1",
			expect: ParsedVersion{
				Major: 1,
				Minor: 3,
				RC:    1,
			},
		},
		{
			in: "go1.16",
			expect: ParsedVersion{
				Major: 1,
				Minor: 16,
			},
		},
		{
			in: "go1.16.12",
			expect: ParsedVersion{
				Major: 1,
				Minor: 16,
				Patch: 12,
			},
		},
		{
			in: "go1.17beta1",
			expect: ParsedVersion{
				Major: 1,
				Minor: 17,
				Beta:  1,
			},
		},
		{
			in: "go1.17rc2",
			expect: ParsedVersion{
				Major: 1,
				Minor: 17,
				RC:    2,
			},
		},
		{
			in: "go1.18",
			expect: ParsedVersion{
				Major: 1,
				Minor: 18,
			},
		},
	}

	for i, tc := range cases {
		got, err := parseVersion(tc.in)
		if err != nil {
			t.Errorf("parse %s err: %s", tc.in, err)
			continue
		}
		if !reflect.DeepEqual(tc.expect, *got) {
			t.Errorf("parse failure %s: got: %+v expect: %+v", tc.in, got, tc.expect)
		}

		if i+1 < len(cases) {
			next := cases[i+1].expect
			if !got.Less(next) {
				t.Fatalf("expected %+v to be less than %+v", got, next)
			}
		}

		if got.String() != tc.in {
			t.Errorf("%s to string is wrong: %s", tc.in, got.String())
		}
	}
}
