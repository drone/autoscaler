// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package engine

import "testing"

func TestAbs(t *testing.T) {
	tests := []struct {
		x, want int
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
	}
	for _, test := range tests {
		if got, want := abs(test.x), test.want; got != want {
			t.Errorf("Want abs value %d, got %d", want, got)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		x, y, want int
	}{
		{0, 1, 1},
		{0, 0, 0},
		{1, 1, 1},
		{-1, 0, 0},
		{-1, 1, 1},
	}
	for _, test := range tests {
		if got, want := max(test.x, test.y), test.want; got != want {
			t.Errorf("Want max value %d, got %d", want, got)
		}
	}
}

func TestServerDiff(t *testing.T) {
	tests := []struct {
		pending, // count pending builds
		available, // count available capacity
		concurrency, // per-server concurrency
		want int
	}{
		// use 2 of 2 existing
		{
			pending:     2,
			available:   2,
			concurrency: 2,
			want:        0,
		},
		// use 1 of 2 existing
		{
			pending:     1,
			available:   2,
			concurrency: 2,
			want:        0,
		},
		// want 1 server
		{
			pending:     4,
			available:   2,
			concurrency: 2,
			want:        1,
		},
		// want 2 servers
		{
			pending:     4,
			available:   2,
			concurrency: 1,
			want:        2,
		},
		// want 2 servers (round-up)
		{
			pending:     5,
			available:   2,
			concurrency: 2,
			want:        2,
		},

		//
		// the following test cases check for instances when
		// we have exceess server capacity and want to remove
		// server instances.
		//

		// want 0 servers removed, at capacity
		{
			pending:     1,
			available:   1,
			concurrency: 2,
			want:        0,
		},
		// want 0 servers removed, at capacity (server partially used)
		{
			pending:     1,
			available:   2,
			concurrency: 2,
			want:        0,
		},
		// want 1 server removed, pending builds, but excess capacity
		{
			pending:     2,
			available:   4,
			concurrency: 2,
			want:        -1,
		},
		// want 2 servers removed (round down)
		{
			pending:     0,
			available:   5,
			concurrency: 2,
			want:        -2,
		},
		// want 10 servers removed
		{
			pending:     4,
			available:   24,
			concurrency: 2,
			want:        -10,
		},
	}
	for _, test := range tests {
		diff := serverDiff(
			test.pending,
			test.available,
			test.concurrency,
		)
		if got, want := diff, test.want; got != want {
			t.Errorf("Got server diff %d, want %d", got, want)
		}
	}
}

func TestSeverCeil(t *testing.T) {
	tests := []struct {
		curr, // count of servers running
		diff, // count of servers to add
		ceil, // max number of servers
		want int
	}{
		// add 0 servers
		{
			curr: 2,
			diff: 0,
			ceil: 2,
			want: 0,
		},
		// add 0 servers, handle 0 current count
		{
			curr: 0,
			diff: 0,
			ceil: 1,
			want: 0,
		},
		// add 1 server
		{
			curr: 2,
			diff: 1,
			ceil: 4,
			want: 1,
		},
		// add 1 server, handle 0 current count
		{
			curr: 0,
			diff: 2,
			ceil: 1,
			want: 1,
		},
		// add 2 servers
		{
			curr: 2,
			diff: 2,
			ceil: 4,
			want: 2,
		},
		// add 2 servers, adjust to ceil
		{
			curr: 2,
			diff: 4,
			ceil: 4,
			want: 2,
		},
		// add 4 servers, adjust to ceil
		{
			curr: 0,
			diff: 10,
			ceil: 4,
			want: 4,
		},
	}
	for _, test := range tests {
		diff := serverCeil(
			test.curr,
			test.diff,
			test.ceil,
		)
		if got, want := diff, test.want; got != want {
			t.Errorf("Got server diff %d, want %d", got, want)
		}
	}
}

func TestSeverFloor(t *testing.T) {
	tests := []struct {
		curr, // count of servers running
		diff, // count of servers to remove
		floor, // min number of servers
		want int
	}{
		// remove 0 servers
		{
			curr:  2,
			diff:  0,
			floor: 2,
			want:  0,
		},
		// remove 1 server
		{
			curr:  4,
			diff:  1,
			floor: 2,
			want:  1,
		},
		// remove 2 servers
		{
			curr:  4,
			diff:  2,
			floor: 2,
			want:  2,
		},
		// remove 2 servers, adjust to floor
		{
			curr:  4,
			diff:  3,
			floor: 2,
			want:  2,
		},
		// remove 0 servers, adjust to floor
		{
			curr:  2,
			diff:  1,
			floor: 2,
			want:  0,
		},
		// should not remove non-existent servers
		{
			curr:  0,
			diff:  4,
			floor: 2,
			want:  0,
		},
		{
			curr:  1,
			diff:  4,
			floor: 2,
			want:  0,
		},
	}
	for _, test := range tests {
		diff := serverFloor(
			test.curr,
			test.diff,
			test.floor,
		)
		if got, want := diff, test.want; got != want {
			t.Errorf("Got server diff %d, want %d", got, want)
		}
	}
}
