// model.go

package cicd_golang

import (
	"database/sql"
)

type product struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price, available FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price, &p.Available)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2, available=$3 WHERE id=$4",
			p.Name, p.Price, p.Available, p.ID)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price, available) VALUES($1, $2, $3) RETURNING id",
		p.Name, p.Price, p.Available).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name, price, available FROM products WHERE available=TRUE LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
