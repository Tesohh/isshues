package config

import "testing"

func TestPriority(t *testing.T) {
	priorities := Priorities{
		"crit": Priority{
			Value:    10,
			ColorKey: "red",
		},
		"high": Priority{
			Value:    5,
			ColorKey: "yellow",
		},
		"med": Priority{
			Value:    3,
			ColorKey: "cyan",
		},
		"default": Priority{
			Value: 1,
		},
		"low": Priority{
			Value:    0,
			ColorKey: "black",
		},
	}

	got := priorities.FindClosest(2)
	want := priorities["default"]
	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}

	got = priorities.FindClosest(1)
	want = priorities["default"]
	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}

	got = priorities.FindClosest(1112)
	want = priorities["crit"]
	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}

	got = priorities.FindClosest(4)
	want = priorities["med"]
	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}

	got = priorities.FindClosest(0)
	want = priorities["low"]
	if got != want {
		t.Errorf("got: %v != want: %v", got, want)
	}
}
