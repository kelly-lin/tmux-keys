package stack

import (
	"errors"
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	for desc, tc := range map[string]struct {
		init []int
		exec func(t *testing.T, s *Stack[int]) []int
		want []int
	}{
		"returns the correct 1 item": {
			init: []int{1},
			exec: func(t *testing.T, s *Stack[int]) []int {
				got := []int{}
				got = append(got, safePop(t, s))
				return got
			},
			want: []int{1},
		},
		"returns the correct 2 items": {
			init: []int{},
			exec: func(t *testing.T, s *Stack[int]) []int {
				s.Push(1, 2)

				got := []int{}
				got = append(got, safePop(t, s))
				got = append(got, safePop(t, s))
				return got
			},
			want: []int{2, 1},
		},
		"returns the correct 3 items": {
			init: []int{},
			exec: func(t *testing.T, s *Stack[int]) []int {
				got := []int{}
				s.Push(1)
				s.Push(2)
				s.Push(3)
				got = append(got, safePop(t, s))
				got = append(got, safePop(t, s))
				s.Push(4)
				s.Push(5)
				got = append(got, safePop(t, s))
				got = append(got, safePop(t, s))
				got = append(got, safePop(t, s))
				return got
			},
			want: []int{3, 2, 5, 4, 1},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			stack := NewStack(tc.init)
			got := tc.exec(t, &stack)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("wanted %v but got %v", tc.want, got)
			}
		})
	}

	t.Run("hasItems", func(t *testing.T) {
		for desc, tc := range map[string]struct {
			items []int
			exec  func(t *testing.T, s *Stack[int])
			want  bool
		}{
			"no items initialised": {
				items: []int{},
				exec:  func(t *testing.T, s *Stack[int]) {},
				want:  false,
			},
			"has items initialised": {
				items: []int{1},
				exec:  func(t *testing.T, s *Stack[int]) {},
				want:  true,
			},
			"has no items after being popped": {
				items: []int{1},
				exec: func(t *testing.T, s *Stack[int]) {
					_, _ = s.Pop()
				},
				want: false,
			},
			"has items after being popped": {
				items: []int{1, 2},
				exec: func(t *testing.T, s *Stack[int]) {
					_, _ = s.Pop()
				},
				want: true,
			},
		} {
			t.Run(desc, func(t *testing.T) {
				stack := NewStack(tc.items)

				tc.exec(t, &stack)

				got := stack.HasItems()

				if got != tc.want {
					t.Fatalf("wanted %t but got %t", tc.want, got)
				}
			})
		}
	})

	t.Run("error when popping an empty stack", func(t *testing.T) {
		stack := NewStack[int]([]int{})

		_, got := stack.Pop()

		want := errors.New("stack is empty")
		if got.Error() != want.Error() {
			t.Fatalf("wanted error %q but got %q", want, got)
		}
	})
}

func safePop(t *testing.T, stack *Stack[int]) int {
	t.Helper()

	item, err := stack.Pop()
	if err != nil {
		t.Fatalf("got an error while popping but we did not expect one: %s", err)
	}
	return item
}
