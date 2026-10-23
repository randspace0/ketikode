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

func TestPasswordHashAndUpdate(t *testing.T) {
	repotest.Open(t)
	hash, err := users.HashPassword("old-password")
	if err != nil || hash == "old-password" {
		t.Fatalf("password must be hashed: %q %v", hash, err)
	}
	u := &users.Users{Username: "a", Email: "a@x.io", Password: hash}
	if err := u.Create(); err != nil {
		t.Fatal(err)
	}

	found, err := users.FindByUsername("a")
	if err != nil || found == nil || !found.CheckPassword("old-password") || found.CheckPassword("wrong") {
		t.Fatalf("lookup/check failed: %+v %v", found, err)
	}
	if missing, err := users.FindByUsername("nobody"); missing != nil || err != nil {
		t.Fatalf("want nil, nil: %v %v", missing, err)
	}

	newHash, _ := users.HashPassword("new-password")
	if err := users.UpdatePassword(u.ID, newHash); err != nil {
		t.Fatal(err)
	}
	found, _ = users.FindByUsername("a")
	if !found.CheckPassword("new-password") || found.CheckPassword("old-password") || found.Email != "a@x.io" {
		t.Fatalf("update must change only the password: %+v", found)
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
