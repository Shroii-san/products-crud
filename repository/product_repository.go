package repository

import (
	"crud-lengkap/model"
	"database/sql"
	"fmt"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAllProduct() ([]model.Product, error) {
	rows, err := r.DB.Query(`SELECT id, name, price, stock FROM products ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}

	var result []model.Product

	for rows.Next() {
		var each model.Product
		err := rows.Scan(&each.ID, &each.Name, &each.Price, &each.Stock)
		if err != nil {
			return nil, err
		}
		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ProductRepository) GetProductById(id int) (model.Product, error) {
	var result model.Product
	err := r.DB.QueryRow(`SELECT id, name, price, stock FROM products WHERE id=$1`, id).Scan(&result.ID, &result.Name, &result.Price, &result.Stock)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
		}
		return model.Product{}, err
	}

	return result, nil

}

func (r *ProductRepository) AddProduct(p model.Product) (model.Product, error) {
	query := `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	err := r.DB.QueryRow(query, p.Name, p.Price, p.Stock).Scan(&p.ID)
	if err != nil {
		return model.Product{}, err
	}

	return p, nil
}

func (r *ProductRepository) UpdateProduct(p model.Product) error {
	stmt, err := r.DB.Prepare(`UPDATE products SET name=$1, price=$2, stock=$3 WHERE id=$4`)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("produk dengan ID %d tidak ditemukan", p.ID)
	}

	return nil
}

func (r *ProductRepository) DeleteProduct(id int) error {
	result, err := r.DB.Exec(`DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
	}

	return nil
}
