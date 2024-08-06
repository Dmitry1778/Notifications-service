package storage

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"notify/internal/domain"
	"strconv"
)

func (db *DB) GetPublisherList(ctx context.Context, id int) (publisher []int, err error) {
	var publisherID string
	row := db.pgClient.QueryRowContext(ctx, `select publisher from emp where empid = $1 `, id)
	if err != nil {
		panic(err.Error())
	}
	err = row.Scan(&publisherID)
	if err != nil {
		return nil, err
	}
	for i, b := range publisherID {
		if i != 0 && i%2 != 0 {
			publisherID = string(b)
			result, _ := strconv.Atoi(publisherID)
			publisher = append(publisher, result)
		}
	}
	return publisher, nil
}

func (db *DB) Subscribe(ctx context.Context, sub, pub int) error {
	_, err := db.GetEmployeeID(ctx, pub)
	if err != nil {
		log.Printf("ID employer is not found: %v\n", err)
		return err
	}
	_, err = db.pgClient.ExecContext(ctx, `UPDATE emp SET publisher = CASE WHEN NOT $2 = ANY(publisher) THEN array_append(publisher, $2) ELSE publisher END WHERE empid = $1`, sub, pub)
	if err != nil {
		log.Printf("Unable to insert data (publisher):%v\n", err)
		return err
	}
	return nil
}

func (db *DB) Unsubscribe(ctx context.Context, sub, pub int) error {
	_, err := db.pgClient.ExecContext(ctx, `UPDATE emp SET publisher = array_remove(publisher, $2) WHERE empid = $1`, sub, pub)
	if err != nil {
		log.Printf("Unable to delete data (publisher):%v\n", err)
		return err
	}
	return nil
}

func (db *DB) GetID(ctx context.Context, username interface{}) (*int, error) {
	var emp domain.Employee
	row := db.pgClient.QueryRowContext(ctx, `select empid from emp where username = $1`, username)
	err := row.Scan(&emp.UserID)
	if err != nil {
		return nil, err
	}
	return &emp.UserID, nil
}

func (db *DB) Get(ctx context.Context, username string) (*domain.Employee, error) {
	var emp domain.Employee
	row := db.pgClient.QueryRowContext(ctx, `select username from emp where username = $1`, username)
	err := row.Scan(&emp.Username)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (db *DB) GetEmployeeID(ctx context.Context, id int) (*domain.ResponseEmployee, error) {
	row := db.pgClient.QueryRowContext(ctx, `select firstname, lastname, day, month, year from emp where empid = $1`, id)
	var emp domain.ResponseEmployee
	err := row.Scan(&emp.Firstname, &emp.Lastname, &emp.Day, &emp.Month, &emp.Year)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (db *DB) GetAllEmployee(ctx context.Context) (*[]domain.ResponseEmployee, error) {
	rows, err := db.pgClient.QueryContext(ctx, `select firstname, lastname, day, month, year from emp`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	empl := make([]domain.ResponseEmployee, 0)
	for rows.Next() {
		emp := domain.ResponseEmployee{}
		err = rows.Scan(&emp.Firstname, &emp.Lastname, &emp.Day, &emp.Month, &emp.Year)
		if err != nil {
			log.Fatal(err)
		}
		empl = append(empl, emp)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return &empl, nil
}

func (db *DB) Register(ctx context.Context, emp domain.Employee) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(emp.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.pgClient.ExecContext(ctx, `insert into emp (password, username, firstname, lastname, day, month, year) values ($1, $2, $3, $4, $5, $6,$7)`, string(hashedPassword), emp.Username, emp.Firstname, emp.Lastname, emp.Day, emp.Month, emp.Year)
	if err != nil {
		log.Printf("Unable to insert data (employee):%v\n", err)
		return err
	}
	return nil
}

func NewDB(pgClient *sql.DB) *DB {
	return &DB{pgClient: pgClient}
}

type DB struct {
	pgClient *sql.DB
}
