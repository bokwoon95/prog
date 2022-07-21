// https://go.dev/play/p/S0jptJmsn70
package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/bokwoon95/sq"
	"github.com/bokwoon95/sqddl/ddl"
)

//go:embed prog.go
var embedFS embed.FS

/*
This is a program that reads its own source code (via //go:embed), parses the
table struct definitions and generates migrations for them. Try commenting and
uncommenting various table structs and struct fields and observing how the
generated migrations change accordingly. Also try changing the dialect between
"sqlite", "postgres", "mysql" and "sqlserver" and see the migrations change to
fit the dialect.

- If you uncomment a table struct, it is added to the generated migrations.

- If you comment out a table struct, it is removed from the generated
migrations.

- If you comment out a table struct that is referenced by another table struct,
the program fails with an error notifying you that no such table exists
(foreign key references are validated).

More about table structs: https://bokwoon.neocities.org/sqddl.html#table-structs

More about ddl struct tags: https://bokwoon.neocities.org/sqddl.html#ddl-struct-tags
*/

func main() {
	generateCmd := &ddl.GenerateCmd{
		Dialect:   "postgres", // "sqlite", "postgres", "mysql", "sqlserver"
		DirFS:     embedFS,
		Filenames: []string{"prog.go"},
	}
	files, warnings, err := generateCmd.Results()
	if err != nil {
		log.Fatal(err)
	}
	for _, warning := range warnings {
		fmt.Println(warning)
	}
	for _, file := range files {
		fileinfo, _ := file.Stat()
		fmt.Printf("\n-- %s\n", fileinfo.Name())
		os.Stdout.ReadFrom(file)
		file.Close()
	}
}

type FILM struct {
	sq.TableStruct
	FILM_ID              sq.NumberField `ddl:"primarykey auto_increment identity"`
	TITLE                sq.StringField `ddl:"notnull len=255 index"`
	DESCRIPTION          sq.StringField
	RELEASE_YEAR         sq.NumberField `ddl:"postgres:type=year"`
	LANGUAGE_ID          sq.NumberField `ddl:"notnull references={language onupdate=cascade ondelete=restrict index}"`
	ORIGINAL_LANGUAGE_ID sq.NumberField `ddl:"references={language.language_id onupdate=cascade ondelete=restrict index}"`
	RENTAL_DURATION      sq.NumberField `ddl:"notnull default=3"`
	RENTAL_RATE          sq.NumberField `ddl:"type=DECIMAL(4,2) sqlite:type=REAL notnull default=4.99"`
	LENGTH               sq.NumberField
	REPLACEMENT_COST     sq.NumberField `ddl:"type=DECIMAL(5,2) sqlite:type=REAL notnull default=19.99"`
	RATING               sq.StringField `ddl:"postgres:type=mpaa_rating default='G'"`
	SPECIAL_FEATURES     sq.ArrayField
	LAST_UPDATE          sq.TimeField `ddl:"notnull default=CURRENT_TIMESTAMP onupdatecurrenttimestamp"`
	FULLTEXT             sq.AnyField  `ddl:"dialect=postgres type=TSVECTOR index={. using=GIST}"`
}

type LANGUAGE struct {
	sq.TableStruct
	LANGUAGE_ID sq.NumberField `ddl:"primarykey auto_increment identity"`
	NAME        sq.StringField `ddl:"notnull len=20"`
	LAST_UPDATE sq.TimeField   `ddl:"notnull default=CURRENT_TIMESTAMP onupdatecurrenttimestamp"`
}

/*
type ACTOR struct {
	sq.TableStruct
	ACTOR_ID    sq.NumberField `ddl:"primarykey auto_increment identity"`
	FIRST_NAME  sq.StringField `ddl:"notnull len=45"`
	LAST_NAME   sq.StringField `ddl:"notnull len=45 index"`
	LAST_UPDATE sq.TimeField   `ddl:"notnull default=CURRENT_TIMESTAMP onupdatecurrenttimestamp"`
}
*/

/*
type FILM_ACTOR struct {
	sq.TableStruct `ddl:"primarykey=actor_id,film_id"`
	ACTOR_ID       sq.NumberField `ddl:"notnull references={actor onupdate=cascade ondelete=restrict}"`
	FILM_ID        sq.NumberField `ddl:"notnull references={film onupdate=cascade ondelete=restrict index}"`
	LAST_UPDATE    sq.TimeField   `ddl:"notnull default=CURRENT_TIMESTAMP onupdatecurrenttimestamp"`
}
*/

/*
type CATEGORY struct {
	sq.TableStruct
	CATEGORY_ID sq.NumberField `ddl:"primarykey auto_increment identity"`
	NAME        sq.StringField `ddl:"notnull len=45"`
	LAST_UPDATE sq.TimeField   `ddl:"notnull default=CURRENT_TIMESTAMP onupdatecurrenttimestamp"`
}
*/

/*
type FILM_CATEGORY struct {
	sq.TableStruct `ddl:"primarykey=film_id,category_id"`
	FILM_ID        sq.NumberField `ddl:"notnull references={film onupdate=cascade ondelete=restrict}"`
	CATEGORY_ID    sq.NumberField `ddl:"notnull references={category onupdate=cascade ondelete=restrict}"`
	LAST_UPDATE    sq.TimeField   `ddl:"notnull default=CURRENT_TIMESTAMP onupdatecurrenttimestamp"`
}
*/
