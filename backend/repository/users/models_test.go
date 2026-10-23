package users_test

import (
	"testing"

	"github.com/randspace0/ketikode/backend/repository"
	"github.com/randspace0/ketikode/backend/repository/migrate"
	"github.com/randspace0/ketikode/backend/repository/programminglang"
	"github.com/randspace0/ketikode/backend/repository/repotest"
	"github.com/randspace0/ketikode/backend/repository/users"
)

func TestCreateEnforcesUniqueUsernameAndEmail(t *testing.T) {
	repotest.Open(t)
	if err := (&users.Users{Username: "a", Email: "a@x.io", Password: "p"}).Create(); err != nil {
		t.Fatal(err)
	}
	if err := (&users.Users{Username: "a", Email: "b@x.io", Password: "p"}).Create(); err == nil {
		t.Error("duplicate username must fail")
	}
	if err := (&users.Users{Username: "b", Email: "a@x.io", Password: "p"}).Create(); err == nil {
		t.Error("duplicate email must fail")
	}
}

func TestFindByIDMissing(t *testing.T) {
	repotest.Open(t)
	if got, err := (&users.Users{ID: 7}).FindByID(); got != nil || err != nil {
		t.Fatalf("want nil, nil; got %v, %v", got, err)
	}
}

func TestSigninLookupAndPasswordUpdate(t *testing.T) {
	repotest.Open(t)
	u := &users.Users{Username: "a", Email: "a@x.io", Password: "old"}
	if err := u.Create(); err != nil {
		t.Fatal(err)
	}

	if got, _ := (&users.Users{Username: "a", Password: "wrong"}).GetUserByUsernamePassword(); got != nil {
		t.Fatal("wrong password must not match")
	}
	if err := (&users.Users{Password: "new"}).UpdateWhereID(u.ID); err != nil {
		t.Fatal(err)
	}
	got, _ := (&users.Users{Username: "a", Password: "new"}).GetUserByUsernamePassword()
	if got == nil || got.Email != "a@x.io" {
		t.Fatalf("update must change only the password: %+v", got)
	}
}

func TestUsernamesByID(t *testing.T) {
	repotest.Open(t)
	a := &users.Users{Username: "a", Email: "a@x.io", Password: "p"}
	b := &users.Users{Username: "b", Email: "b@x.io", Password: "p"}
	a.Create()
	b.Create()

	names, err := users.UsernamesByID([]uint64{a.ID, b.ID, 999})
	if err != nil || len(names) != 2 || names[a.ID] != "a" || names[b.ID] != "b" {
		t.Fatalf("%v %v", names, err)
	}
	if names, err := users.UsernamesByID(nil); err != nil || len(names) != 0 {
		t.Fatalf("empty input: %v %v", names, err)
	}
}

func TestMigrateSeedsLanguagesOnce(t *testing.T) {
	repotest.Open(t)
	if err := migrate.Run(repository.DB); err != nil { // second run must be a no-op
		t.Fatal(err)
	}
	if n := len(programminglang.AllList()); n != 6 {
		t.Fatalf("want 6 languages, got %d", n)
	}
}
