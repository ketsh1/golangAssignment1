package data

import (
	"assignment1/internal/validator"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type Module struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	ModuleName     string    `json:"module_name"`
	ModuleDuration int32     `json:"module_duration,omitempty"`
	ExamType       []string  `json:"exam_type,omitempty"`
	Version        int32     `json:"version"`
}

func ValidateModule(v *validator.Validator, module *Module) {
	v.Check(module.ModuleName != "", "module_name", "must be provided")
	v.Check(len(module.ModuleName) <= 500, "module_name", "must not be more than 500 bytes long")
	/*v.Check(module.Year != 0, "year", "must be provided")
	v.Check(module.Year >= 1888, "year", "must be greater than 1888")
	v.Check(module.Year <= int32(time.Now().Year()), "year", "must not be in the future")*/
	v.Check(module.ModuleDuration != 0, "module_duration", "must be provided")
	v.Check(module.ModuleDuration > 0, "module_duration", "must be a positive integer")
	v.Check(module.ExamType != nil, "exam_type", "must be provided")
	v.Check(len(module.ExamType) >= 1, "exam_type", "must contain at least 1 genre")
	v.Check(len(module.ExamType) <= 5, "exam_type", "must not contain more than 5 genres")
	v.Check(validator.Unique(module.ExamType), "exam_type", "must not contain duplicate values")
}

// Define a ModulesModel struct type which wraps a sql.DB connection pool.
type ModuleModel struct {
	DB *sql.DB
}

// method for inserting info to table
func (m ModuleModel) Insert(module *Module) error {
	//SQL query
	query := `
INSERT INTO module_info (module_name, module_duration,exam_type)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at, version`

	args := []any{module.ModuleName, module.ModuleDuration, pq.Array(module.ExamType)}

	return m.DB.QueryRow(query, args...).Scan(&module.ID, &module.CreatedAt, &module.UpdatedAt, &module.Version)
}

// method for fetching
func (m ModuleModel) Get(id int64) (*Module, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//SQL query
	query := `SELECT id, created_at, updated_at, module_name, module_duration, exam_type, version 
FROM module_info 
WHERE id =$1`

	var module Module

	err := m.DB.QueryRow(query, id).Scan(
		&module.ID,
		&module.CreatedAt,
		&module.UpdatedAt,
		&module.ModuleName,
		&module.ModuleDuration,
		pq.Array(&module.ExamType),
		&module.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return &module, nil

}

// method for updating
func (m ModuleModel) Update(module *Module) error {
	query := `
	UPDATE module_info
	SET  module_name =$1, module_duration=$2, exam_type=$3, version = version + 1
	WHERE id=$4
	RETURNING version
`

	args := []any{
		module.ModuleName,
		module.ModuleDuration,
		pq.Array(module.ExamType),
		module.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&module.Version)

}

// method for deleting
func (m ModuleModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM module_info
	WHERE id = $1
`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	RowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if RowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
