package repo

import (
	"database/sql"
	"fmt"
	"rsoi/internal/model"
)

type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo создает новый экземпляр репозитория
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Add добавляет новую запись о человеке в базу данных
func (r *PostgresRepo) Add(p *model.Person) error {
	query := `
		INSERT INTO persons (name, age, address, work)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		p.Name,
		p.Age,
		p.Address,
		p.Work,
	).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("failed to add person: %w", err)
	}

	return nil
}

// Update обновляет существующую запись о человеке
func (r *PostgresRepo) Update(p *model.Person) error {
	query := `
		UPDATE persons 
		SET name = $1, age = $2, address = $3, work = $4
		WHERE id = $5
	`

	result, err := r.db.Exec(
		query,
		p.Name,
		p.Age,
		p.Address,
		p.Work,
		p.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete удаляет запись о человеке по ID
func (r *PostgresRepo) Delete(id int) error {
	query := `DELETE FROM persons WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// GetAll возвращает всех людей из базы данных
func (r *PostgresRepo) GetAll() ([]model.Person, error) {
	query := `
		SELECT id, name, age, address, work
		FROM persons
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all persons: %w", err)
	}
	defer rows.Close()

	var persons []model.Person
	for rows.Next() {
		var p model.Person
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Age,
			&p.Address,
			&p.Work,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}
		persons = append(persons, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return persons, nil
}

// GetById возвращает человека по ID
func (r *PostgresRepo) GetById(id int) (*model.Person, error) {
	query := `
		SELECT id, name, age, address, work
		FROM persons
		WHERE id = $1
	`

	var p model.Person
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Age,
		&p.Address,
		&p.Work,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get person by id: %w", err)
	}

	return &p, nil
}
