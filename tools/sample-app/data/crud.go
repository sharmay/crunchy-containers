package data

import (
	"database/sql"
	"fmt"
)

type Coffee struct {
	Name    string
	Species string
	Regions string
	Comment string
}

func (db *DB) InitDDL() error {
	statement := `
CREATE TABLE public.coffee(
	name TEXT PRIMARY KEY, 
	species TEXT, 
	regions TEXT, 
	comment TEXT
)`
	_, err := db.Exec(statement)
	return err
}

func (db *DB) CoffeeDeleteAll() error {
	statement := "DELETE FROM public.coffee"
	_, err := db.Exec(statement)
	return err
}

func (db *DB) AddCoffee(c Coffee) error {
	statement := "INSERT INTO public.coffee" +
		"(name, species, regions, comment) " +
		"VALUES ($1, $2, $3, $4) " +
		"RETURNING name, species, regions, comment"

	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(c.Name, c.Species, c.Regions,
		c.Comment).Scan(&c.Name, &c.Species, &c.Regions, &c.Comment)
	return err
}

func (db *DB) AllCoffee() ([]Coffee, error) {
	statement := "SELECT * FROM public.coffee"

	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	coffees := []Coffee{}
	for rows.Next() {
		var name, species, regions, comment sql.NullString
		if err := rows.Scan(&name, &species, &regions, &comment); err != nil {
			return nil, err
		}
		c := Coffee{
			Name:    name.String,
			Species: species.String,
			Regions: regions.String,
			Comment: comment.String,
		}
		coffees = append(coffees, c)
	}
	return coffees, nil
}

func (db *DB) Coffee(name string) (Coffee, error) {
	rows, err := db.Query("SELECT * FROM public.coffee WHERE name=$1", name)
	if err != nil {
		return Coffee{}, err
	}
	defer rows.Close()

	c := Coffee{}
	for rows.Next() {
		if err := rows.Scan(&c.Name, &c.Species, &c.Regions, &c.Comment); err != nil {
			return Coffee{}, err
		}
	}
	return c, nil
}

type CRUDResult struct {
	Name    string
	Success bool
	Error   string
}

func (db *DB) RunCRUD() []CRUDResult {
	var results []CRUDResult

	// DDL
	ddl := CRUDResult{
		Name:    "DDL",
		Success: false,
	}

	err := db.InitDDL()
	if err != nil {
		ddl.Error = fmt.Sprintf("%s", err)
	}
	ddl.Success = true
	results = append(results, ddl)

	// Delete
	delete := CRUDResult{
		Name:    "Delete",
		Success: false,
	}

	err = db.CoffeeDeleteAll()
	if err != nil {
		delete.Error = fmt.Sprintf("%s", err)
	}
	delete.Success = true
	results = append(results, delete)

	// Write tests
	write := CRUDResult{
		Name:    "Write",
		Success: false,
	}

	for _, coffee := range coffees {
		err := db.AddCoffee(coffee)
		if err != nil {
			write.Error = fmt.Sprintf("%s", err)
		}
	}
	write.Success = true
	results = append(results, write)

	// Read
	read := CRUDResult{
		Name:    "Read",
		Success: false,
	}

	coffees, err := db.AllCoffee()
	if err != nil {
		read.Error = fmt.Sprintf("%s", err)
	}

	if len(coffees) > 0 {
		read.Success = true
	}
	results = append(results, read)

	return results
}
