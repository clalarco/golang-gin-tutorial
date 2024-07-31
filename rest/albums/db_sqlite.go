package albums

import (
	"database/sql"
	"errors"

	"clalarco.io/helpers"
)

type db_sqlite struct {
	instance helpers.DbSqlite
}

func (sqlite *db_sqlite) Init() error {
	instance, err := helpers.GetSqlite3Connection()
	sqlite.instance = instance
	if err != nil {
		return err
	}
	schemaErr := sqlite.createSchema()
	if schemaErr != nil {
		return schemaErr
	}
	return nil
}

func (sqlite *db_sqlite) GetAlbums() ([]album, error) {
	rows, err := sqlite.instance.DB.Query("SELECT * FROM albums")
	if err != nil {
		return []album{}, err
	}

	var albums []album = processRows(rows)
	defer rows.Close()
	return albums, nil
}

func (sqlite *db_sqlite) GetAlbum(id string) (album, error) {
	rows, err := sqlite.instance.DB.Query("SELECT * FROM albums WHERE id = ?", id)
	if err != nil {
		return album{}, err
	}

	var albums []album = processRows(rows)
	defer rows.Close()
	if len(albums) == 0 {
		return album{}, errors.New("ID not found")
	}
	return albums[0], nil
}

func (sqlite *db_sqlite) AddAlbum(album album) error {
	_, err := sqlite.instance.DB.Exec("INSERT INTO albums VALUES (?, ?, ?, ?)", album.ID, album.Title, album.Artist, album.Price)
	return err
}

func (sqlite *db_sqlite) DeleteAlbum(id string) error {
	_, err := sqlite.instance.DB.Exec("DELETE FROM albums WHERE id = ?", id)
	return err
}

func (sqlite *db_sqlite) createSchema() error {
	_, err := sqlite.instance.DB.Exec("CREATE TABLE IF NOT EXISTS albums (id TEXT PRIMARY KEY, title TEXT, artist TEXT, price REAL)")
	return err
}

func processRows(rows *sql.Rows) []album {
	var albums []album
	for rows.Next() {
		var id string
		var title string
		var artist string
		var price float64
		err := rows.Scan(&id, &title, &artist, &price)
		if err != nil {
			return []album{}
		}
		albums = append(albums, album{ID: id, Title: title, Artist: artist, Price: price})
	}
	return albums
}
