package main

import (
	"database/sql"
	"fmt"
)

type product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []product{}

	for rows.Next() {
		var p product
		err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *product) getProductsById(db *sql.DB) error {
	query := fmt.Sprintf("SELECT id,name,quantity,price from products where id=%v", p.Id)
	rows := db.QueryRow(query)
	err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
	if err != nil {

		return err
	}

	return nil
}

func (p *product) insertProducts(db *sql.DB) error {
	query := "Insert into products (name,quantity,price)values(?,?,?)"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	p.Id = int(id)

	return nil
}

func (p *product) updateDetails(db *sql.DB) error {

	query := "UPDATE products SET name =?, quantity =?, price =? WHERE id =?"

	_, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.Id)

	return err
}

func (p *product) deleteProducts(db *sql.DB) error {
	query := "DELETE  from products where id =?"
	_, err := db.Exec(query, p.Id)
	return err
}
