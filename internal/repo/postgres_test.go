package repo

import (
	"database/sql"
	"errors"
	"rsoi/internal/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_Add(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	tests := []struct {
		name    string
		person  *model.Person
		mock    func()
		wantErr bool
	}{
		{
			name: "successful add",
			person: &model.Person{
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
				Work:    "Developer",
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO persons").
					WithArgs("John Doe", 30, "123 Main St", "Developer").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "database error",
			person: &model.Person{
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
				Work:    "Developer",
			},
			mock: func() {
				mock.ExpectQuery("INSERT INTO persons").
					WithArgs("John Doe", 30, "123 Main St", "Developer").
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.Add(tt.person)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, tt.person.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	tests := []struct {
		name    string
		person  *model.Person
		mock    func()
		wantErr bool
	}{
		{
			name: "successful update",
			person: &model.Person{
				ID:      1,
				Name:    "John Doe Updated",
				Age:     31,
				Address: "456 Oak St",
				Work:    "Senior Developer",
			},
			mock: func() {
				mock.ExpectExec("UPDATE persons").
					WithArgs("John Doe Updated", 31, "456 Oak St", "Senior Developer", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "person not found",
			person: &model.Person{
				ID:      999,
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
				Work:    "Developer",
			},
			mock: func() {
				mock.ExpectExec("UPDATE persons").
					WithArgs("John Doe", 30, "123 Main St", "Developer", 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "database error",
			person: &model.Person{
				ID:      1,
				Name:    "John Doe",
				Age:     30,
				Address: "123 Main St",
				Work:    "Developer",
			},
			mock: func() {
				mock.ExpectExec("UPDATE persons").
					WithArgs("John Doe", 30, "123 Main St", "Developer", 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.Update(tt.person)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	tests := []struct {
		name    string
		id      int
		mock    func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   1,
			mock: func() {
				mock.ExpectExec("DELETE FROM persons WHERE id = \\$1").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "person not found",
			id:   999,
			mock: func() {
				mock.ExpectExec("DELETE FROM persons WHERE id = \\$1").
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "database error",
			id:   1,
			mock: func() {
				mock.ExpectExec("DELETE FROM persons WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.Delete(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	tests := []struct {
		name     string
		mock     func()
		expected []model.Person
		wantErr  bool
	}{
		{
			name: "successful get all",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"}).
					AddRow(1, "John Doe", 30, "123 Main St", "Developer").
					AddRow(2, "Jane Smith", 25, "456 Oak St", "Designer")
				mock.ExpectQuery("SELECT id, name, age, address, work FROM persons ORDER BY id").
					WillReturnRows(rows)
			},
			expected: []model.Person{
				{ID: 1, Name: "John Doe", Age: 30, Address: "123 Main St", Work: "Developer"},
				{ID: 2, Name: "Jane Smith", Age: 25, Address: "456 Oak St", Work: "Designer"},
			},
			wantErr: false,
		},
		{
			name: "empty result",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"})
				mock.ExpectQuery("SELECT id, name, age, address, work FROM persons ORDER BY id").
					WillReturnRows(rows)
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name: "database error",
			mock: func() {
				mock.ExpectQuery("SELECT id, name, age, address, work FROM persons ORDER BY id").
					WillReturnError(errors.New("database error"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.GetAll()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRepo_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	tests := []struct {
		name     string
		id       int
		mock     func()
		expected *model.Person
		wantErr  bool
	}{
		{
			name: "successful get by id",
			id:   1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "address", "work"}).
					AddRow(1, "John Doe", 30, "123 Main St", "Developer")
				mock.ExpectQuery("SELECT id, name, age, address, work FROM persons WHERE id = \\$1").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expected: &model.Person{ID: 1, Name: "John Doe", Age: 30, Address: "123 Main St", Work: "Developer"},
			wantErr:  false,
		},
		{
			name: "person not found",
			id:   999,
			mock: func() {
				mock.ExpectQuery("SELECT id, name, age, address, work FROM persons WHERE id = \\$1").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "database error",
			id:   1,
			mock: func() {
				mock.ExpectQuery("SELECT id, name, age, address, work FROM persons WHERE id = \\$1").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			result, err := repo.GetById(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
