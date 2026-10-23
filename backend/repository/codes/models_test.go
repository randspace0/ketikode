package codes_test

import (
	"strings"
	"testing"

	"github.com/randspace0/ketikode/backend/repository/codes"
	"github.com/randspace0/ketikode/backend/repository/repotest"
)

func create(t *testing.T, language string, owner uint64) *codes.Codes {
	t.Helper()
	c := &codes.Codes{Language: language, Code: "x", Owner: owner}
	if err := c.Create(); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestLongTextRoundTrips(t *testing.T) {
	repotest.Open(t)
	long := strings.Repeat("a", 10000)
	c := &codes.Codes{Language: "cpp", Code: long, Stdout: long, Stderr: long, Stdin: long, Owner: 1}
	if err := c.Create(); err != nil {
		t.Fatal(err)
	}
	got, err := (&codes.Codes{ID: c.ID}).FindByID()
	if err != nil || got == nil || got.Code != long || got.Stdout != long {
		t.Fatalf("err=%v got=%v", err, got)
	}
}

func TestFindByIDMissing(t *testing.T) {
	repotest.Open(t)
	got, err := (&codes.Codes{ID: 42}).FindByID()
	if got != nil || err != nil {
		t.Fatalf("want nil, nil; got %v, %v", got, err)
	}
}

func TestDeleteChecksOwner(t *testing.T) {
	repotest.Open(t)
	c := create(t, "cpp", 1)

	if err := (&codes.Codes{ID: c.ID, Owner: 2}).Delete(); err == nil {
		t.Fatal("another user must not be able to delete the code")
	}
	if got, _ := (&codes.Codes{ID: c.ID}).FindByID(); got == nil {
		t.Fatal("code was deleted by a non-owner")
	}
	if err := (&codes.Codes{ID: c.ID, Owner: 1}).Delete(); err != nil {
		t.Fatal(err)
	}
	if got, _ := (&codes.Codes{ID: c.ID}).FindByID(); got != nil {
		t.Fatal("owner delete did not remove the code")
	}
}

func TestListAndTotal(t *testing.T) {
	repotest.Open(t)
	create(t, "cpp", 1)
	create(t, "python3", 1)
	last := create(t, "cpp", 2)

	all, err := codes.List("", 10, 0)
	if err != nil || len(all) != 3 || all[0].ID != last.ID {
		t.Fatalf("want 3 newest first, got %+v (%v)", all, err)
	}
	if page, _ := codes.List("", 2, 2); len(page) != 1 {
		t.Fatalf("offset page: %+v", page)
	}
	if cpp, _ := codes.List("cpp", 10, 0); len(cpp) != 2 {
		t.Fatalf("filter: %+v", cpp)
	}
	if n, _ := (&codes.Codes{}).Total(); n != 3 {
		t.Fatalf("total=%d", n)
	}
	if n, _ := (&codes.Codes{Language: "cpp"}).Total(); n != 2 {
		t.Fatalf("cpp total=%d", n)
	}
}

func TestOwnedCodesAreScopedToUser(t *testing.T) {
	repotest.Open(t)
	create(t, "cpp", 1)
	create(t, "cpp", 2)
	mine, err := (&codes.Codes{Owner: 1}).FindAllOwnedCodesByUserID()
	if err != nil || len(mine) != 1 {
		t.Fatalf("%+v %v", mine, err)
	}
}
