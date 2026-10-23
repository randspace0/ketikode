package main

import (
	"github.com/go-xorm/xorm"

	"gitlab.com/cikadev/ketide/repository/codes"
	"gitlab.com/cikadev/ketide/repository/programminglang"
	"gitlab.com/cikadev/ketide/repository/users"
)

func MigrateTable(db *xorm.Engine) {
	if err := db.Sync(new(users.Users)); err != nil {
		panic("Table users migrate error")
	}

	if err := db.Sync(new(codes.Codes)); err != nil {
		panic("Table code migrate error")
	}

	if err := db.Sync(new(programminglang.ProgrammingLang)); err != nil {
		panic("Table code migrate error")
	}
}

func MigrateLanguage(db *xorm.Engine) {
	languageLists := []programminglang.ProgrammingLang{
		programminglang.ProgrammingLang{
			LanguageID: "c",
			Name:       "C",
		},
		programminglang.ProgrammingLang{
			LanguageID: "cpp",
			Name:       "C++",
		},
		programminglang.ProgrammingLang{
			LanguageID: "ruby",
			Name:       "Ruby",
		},
		programminglang.ProgrammingLang{
			LanguageID: "python3",
			Name:       "Python 3",
		},
		programminglang.ProgrammingLang{
			LanguageID: "java",
			Name:       "Java",
		},
		programminglang.ProgrammingLang{
			LanguageID: "rust",
			Name:       "Rust",
		},
	}

	for _, language := range languageLists {
		language.Create()
	}
}

// func dropTable(db *xorm.Engine) {
// 	if err := db.DropTables(new(Users)); err != nil {
// 		panic("Table users drop error")
// 	}

// 	if err := db.DropTables(new(Code)); err != nil {
// 		panic("Table code drop error")
// 	}
// }
