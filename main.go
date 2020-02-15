package core

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type DB struct {
	DB *sql.DB
}

// ErrInvalidPass ...
var ErrInvalidPass = errors.New("Неверный пароль")

// QueryError ...
type QueryError struct { // alt + enter
	Query string
	Err   error
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}
func (receiver *QueryError) Error() string {
	return fmt.Sprintf("Не удалось выполнить запрос %s: %s", loginSQL, receiver.Err.Error())
}

// DbError ...
type DbError struct {
	Err error
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

func (receiver DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

// for managers =================================

// AddAccountToUser ...
func AddAccountToUser(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET account = ? WHERE id = ?;")
	if err != nil {
		log.Println(err)
	}
	return err
}

// CreateService ...
func CreateService(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO services(name, price) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
	}
	return err
}

// CreateATM ...
func CreateATM(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO atm_list(name, address) VALUES (?, ?)")
	if err != nil {
		log.Println(err)
	}
	return err
}

// for users ==========================================

// ошибки - это тоже часть API

const loginSQL = `SELECT login, password FROM managers WHERE login = ?`

// Login manager
func Login(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string
	err := db.QueryRow(
		loginSQL,
		login).Scan(&dbLogin, &dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, queryError(loginSQL, err)
	}
	if dbPassword != password {
		return false, ErrInvalidPass
	}
	return true, nil
}

var loginUser = `SELECT login, password FROM users WHERE login = ?`

// LoginUser for login users in cli
func LoginUser(login, password string, db *sql.DB) (bool, error) {
	var ursLogin, usrPassword string
	err := db.QueryRow(
		loginUser, login).Scan(&ursLogin, &usrPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, queryError(loginUser, err)
	}
	if usrPassword != password {
		return false, ErrInvalidPass
	}
	return true, err
}

const getServicesSQL = `SELECT name, price FROM services;`

// ShowAllServices ...
func ShowAllServices(db *sql.DB) ([]Service, error) {
	services := make([]Service, 0)

	rows, err := db.Query(getServicesSQL)
	if err != nil {
		return nil, queryError(getServicesSQL, err)
	}

	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			services, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		s := new(Service)
		err := rows.Scan(&s.Name, &s.Price)
		if err != nil {
			return nil, dbError(err)
		}
		services = append(services, *s)
	}

	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return services, nil
}

const allUserAccounts = `SELECT account FROM users;`

// AllUserAccounts ...
func AllUserAccounts(db *sql.DB, userID string) ([]UserType, error) {
	accounts := make([]UserType, 0)
	rows, err := db.Query(allUserAccounts, userID)
	if err != nil {
		return nil, queryError(allUserAccounts, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			accounts, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		a := new(UserType)
		err := rows.Scan(&a.Name, &a.Account)
		if err != nil {
			return nil, dbError(err)
		}
		accounts = append(accounts, *a)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return accounts, nil
}

// UseService ... NEED CHECK THIS >>>>>
func UseService(serviceID string, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	var (
		currentPrice int64
		name         string
	)

	err = tx.QueryRow(`SELECT name, price from services WHERE id = ?`, serviceID).Scan(
		&currentPrice, &name)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO sales(user_id, service_id, price) VALUES (:user_id, :service_id, :price)`,
		sql.Named("user_id", 1),
		sql.Named("service_id", serviceID),
		sql.Named("price", currentPrice))

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

const getATMSQL = `SELECT name, price FROM services;`

// ShowAllATMs ...
func ShowAllATMs(db *sql.DB) ([]ATM, error) {
	atms := make([]ATM, 0)

	rows, err := db.Query(getATMSQL)
	if err != nil {
		return nil, queryError(getATMSQL, err)
	}

	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		a := new(ATM)
		err := rows.Scan(&a.Name, &a.Address)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, *a)
	}

	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

// CreateNewUser ...
func CreateNewUser(db *sql.DB, data *UserType) error {
	log.Printf("%+v", data)
	sqlQuery := `INSERT INTO users (name, phone, account, login, password, balance) VALUES (?, ?, ?, ?, ?, ?);`
	_, err := db.Exec(sqlQuery, data.Name, data.Phone, data.Account, data.Login, data.Password, data.Balance)
	return err
}

func checkUserBalance(db *sql.DB, userID string) (int, error) {
	var balance int
	query := "SELECT balance FROM users WHERE id = ?"
	if len(userID) == 9 {
		query = "SELECT balance FROM users WHERE phone = ?"
	}
	err := db.QueryRow(query, userID).Scan(&balance)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return balance, err
}

// SendMoney ...
func SendMoneyByID(db *sql.DB, firstID, secondID, value string) error {
	balance, err := checkUserBalance(db, firstID)
	if err != nil {
		log.Println(err)
	}
	val, err := strconv.Atoi(value)
	if balance < val {
		return errors.New("Недостаточно средств на счете")
	}

	tx, err := db.Begin()
	handleError(err, "Transaction begin error: ")

	// insert a record into table1
	res, err := tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ?", value, firstID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Println("Transaction Rollback Error:", err)
		}
		log.Fatalf("Transaction execute error: %s", err)
	}

	// fetch the auto incremented id
	_, err = res.LastInsertId()
	handleError(err, "LastInsertId error")

	// insert record into table2, referencing the first record from table1
	res, err = tx.Exec("UPDATE users SET balance = balance + ? WHERE id = ?", value, secondID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Println("Ошибка базы данных:", err)
		}
		log.Printf("Ошибка запроса: %s", err)
	}

	// commit the transaction
	handleError(tx.Commit(), "Ошибка завершения транзакции")

	return err
}

// SendMoneyByPhoneNumber ...
func SendMoneyByPhoneNumber(db *sql.DB, firstID, secondUserPhone, value string) error {
	balance, err := checkUserBalance(db, firstID)
	if err != nil {
		log.Println(err)
	}
	val, err := strconv.Atoi(value)
	if balance < val {
		return errors.New("Недостаточно средств на счете")
	}

	tx, err := db.Begin()
	handleError(err, "Transaction begin error: ")

	// insert a record into table1
	res, err := tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ?", value, firstID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Println("Transaction Rollback Error:", err)
		}
		log.Fatalf("Transaction execute error: %s", err)
	}

	// fetch the auto incremented id
	_, err = res.LastInsertId()
	handleError(err, "LastInsertId error")

	// insert record into table2, referencing the first record from table1
	res, err = tx.Exec("UPDATE users SET balance = balance + ? WHERE phone = ?", value, secondUserPhone)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			log.Println("Ошибка базы данных:", err)
		}
		log.Printf("Ошибка запроса: %s", err)
	}

	// commit the transaction
	handleError(tx.Commit(), "Ошибка завершения транзакции")

	return err
}

func ExportAtm(db *sql.DB) ([]ATM, error) {
	atms := make([]ATM, 0)
	query := `SELECT name, addres FROM atm_list;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, queryError(allUserAccounts, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()
	for rows.Next() {
		a := new(ATM)
		err := rows.Scan(&a.Name, &a.Address)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, *a)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

// ExportData ...
func ExportData(db *sql.DB, query string) ([]UserType, error) {
	users := make([]UserType, 0)
	rows, err := db.Query(query)

	if err != nil {
		return nil, queryError(query, err)
	}

	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			users, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		a := new(UserType)
		err := rows.Scan(&a.Name, &a.Account)
		if err != nil {
			return nil, dbError(err)
		}
	}

	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return users, nil
}

func handleError(err error, message string) {
	if err != nil {
		log.Println(message, err)
	}
}
