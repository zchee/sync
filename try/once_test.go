package try

import (
	"errors"
	"fmt"
	"testing"
)

func TestOnce(t *testing.T) {
	expectedErr := errors.New("error")
	cases := map[string]struct {
		once Once
		ok   bool
		err  error
		f    func(o *Once) (bool, error)
	}{
		"do once": {ok: true, f: func(o *Once) (ok bool, err error) {
			o.Do(func() { ok = true })
			return
		}},
		"do twice": {ok: false, f: func(o *Once) (ok bool, err error) {
			o.Do(func() {})
			o.Do(func() { ok = true })
			return
		}},
		"try once": {ok: true, f: func(o *Once) (ok bool, err error) {
			err = o.Try(func() error { ok = true; return nil })
			return
		}},
		"try twice": {ok: false, f: func(o *Once) (ok bool, err error) {
			err = o.Try(func() error { return nil })
			o.Try(func() error { ok = true; return nil })
			return
		}},
		"try error": {ok: true, f: func(o *Once) (ok bool, err error) {
			err = o.Try(func() error { return expectedErr })
			o.Try(func() error { ok = true; return nil })
			return
		}, err: expectedErr},
		"do and try": {ok: false, f: func(o *Once) (ok bool, err error) {
			o.Do(func() {})
			err = o.Try(func() error { ok = true; return nil })
			return
		}},
		"try and do": {ok: false, f: func(o *Once) (ok bool, err error) {
			err = o.Try(func() error { return nil })
			o.Do(func() { ok = true })
			return
		}},
		"try error and do": {ok: true, f: func(o *Once) (ok bool, err error) {
			err = o.Try(func() error { return expectedErr })
			o.Do(func() { ok = true })
			return
		}, err: expectedErr},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			ok, err := tc.f(&tc.once)
			if ok != tc.ok {
				t.Errorf("want %v got %v", tc.ok, ok)
			}
			if err != tc.err {
				t.Errorf("want %v got %v", tc.err, err)
			}
		})
	}
}

func Example_Once_Try() {
	var once Once
	for i := 1; i <= 3; i++ {
		i := i
		err := once.Try(func() error {
			if i < 3 {
				return errors.New("error")
			}
			return nil
		})
		if err != nil {
			fmt.Printf("try %d %v\n", i, err)
		} else {
			fmt.Printf("try %d success\n", i)
		}
	}
	// Output:
	// try 1 error
	// try 2 error
	// try 3 success
}
