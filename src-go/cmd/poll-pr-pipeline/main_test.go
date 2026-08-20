package main

import "testing"

func TestStatus(t *testing.T) {
	cases := []struct {
		name         string
		checks       []check
		done, failed bool
	}{
		{"pending", []check{{Bucket: "pass"}, {Bucket: "pending"}}, false, false},
		{"success", []check{{Bucket: "pass"}, {Bucket: "skipping"}}, true, false},
		{"failure", []check{{Bucket: "pass"}, {Bucket: "cancel"}}, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			done, failed := status(tc.checks)
			if done != tc.done || failed != tc.failed {
				t.Fatalf("status() = (%t, %t), want (%t, %t)", done, failed, tc.done, tc.failed)
			}
		})
	}
}
