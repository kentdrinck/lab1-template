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

func TestPostgresRepo_UpdateField(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepo(db)

	tests := []struct {
		name        string
		personId    int
		field       string
		value       any
		mock        func()
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "successful update name",
			personId: 1,
			field:    "name",
			value:    "John Updated",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET name = \\$1 WHERE id = \\$2").
					WithArgs("John Updated", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "successful update age",
			personId: 1,
			field:    "age",
			value:    35,
			mock: func() {
				mock.ExpectExec("UPDATE persons SET age = \\$1 WHERE id = \\$2").
					WithArgs(35, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "successful update address",
			personId: 1,
			field:    "address",
			value:    "New Address 123",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET address = \\$1 WHERE id = \\$2").
					WithArgs("New Address 123", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "successful update work",
			personId: 1,
			field:    "work",
			value:    "Senior Developer",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET work = \\$1 WHERE id = \\$2").
					WithArgs("Senior Developer", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "person not found",
			personId: 999,
			field:    "name",
			value:    "Not Found",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET name = \\$1 WHERE id = \\$2").
					WithArgs("Not Found", 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr:     true,
			expectedErr: ErrNotFound,
		},
		{
			name:     "database error",
			personId: 1,
			field:    "name",
			value:    "John",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET name = \\$1 WHERE id = \\$2").
					WithArgs("John", 1).
					WillReturnError(errors.New("database connection failed"))
			},
			wantErr: true,
		},
		{
			name:     "invalid field name",
			personId: 1,
			field:    "invalid_field",
			value:    "value",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET invalid_field = \\$1 WHERE id = \\$2").
					WithArgs("value", 1).
					WillReturnError(errors.New("column \\\"invalid_field\\\" of relation \\\"persons\\\" does not exist"))
			},
			wantErr: true,
		},
		{
			name:     "rows affected error",
			personId: 1,
			field:    "name",
			value:    "John",
			mock: func() {
				result := sqlmock.NewErrorResult(errors.New("rows affected error"))
				mock.ExpectExec("UPDATE persons SET name = \\$1 WHERE id = \\$2").
					WithArgs("John", 1).
					WillReturnResult(result)
			},
			wantErr: true,
		},
		{
			name:     "update with nil value",
			personId: 1,
			field:    "age",
			value:    nil,
			mock: func() {
				mock.ExpectExec("UPDATE persons SET age = \\$1 WHERE id = \\$2").
					WithArgs(nil, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:     "update with special characters",
			personId: 1,
			field:    "address",
			value:    "123 O'Reilly Street, Apt #4",
			mock: func() {
				mock.ExpectExec("UPDATE persons SET address = \\$1 WHERE id = \\$2").
					WithArgs("123 O'Reilly Street, Apt #4", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := repo.UpdateField(tt.personId, tt.field, tt.value)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				assert.NoError(t, err)
			}

			// Проверяем что все ожидания выполнены
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
