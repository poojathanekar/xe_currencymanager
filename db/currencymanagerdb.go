package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"xe_currencymanager/config"
	"xe_currencymanager/model"

	logger "github.com/sirupsen/logrus"

	_ "github.com/lib/pq" // here
)

type currencies struct {
	Currencies []string `json:"currencies"`
}

//makeConnection create connection
func makeConnection() (*sql.DB, error) {
	pgcon := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		config.GetConfig("pgdb_config.host"),
		config.GetConfig("pgdb_config.port"),
		config.GetConfig("pgdb_config.user"),
		config.GetConfig("pgdb_config.password"),
		config.GetConfig("pgdb_config.dbname"))
	db, err := sql.Open("postgres", pgcon)
	if err != nil {
		logger.WithField("err", err.Error()).Error(" Failed to initialize database")
		return nil, err
	}

	return db, nil
}

func createCurrencyTableIfNotExist() error {
	createTableQuery := `CREATE TABLE IF NOT EXISTS currencydata(amount numeric,
  rate double precision,
  fromcurrency character(3),
  tocurrency character(3),
	timestamp time with time zone,
	CONSTRAINT fromunq UNIQUE (tocurrency,fromcurrency)
		 )`
	db, dberr := makeConnection()
	if dberr != nil {
		logger.WithField("dberr", dberr.Error()).Error(" Failed to initialize database")
		return dberr
	}
	_, err := db.Exec(createTableQuery)

	if err != nil {
		logger.WithField("err", err.Error()).Error(" Failed to create table")
		return err
	}
	defer db.Close()
	return nil
}

//UpdateResponseData to update in postgresql db
func UpdateResponseData(responseData model.XEResponse, queryValues []string, queryFields []interface{}) error {
	pgStatement := fmt.Sprintf(`INSERT INTO currencydata
			(amount,rate,fromcurrency,tocurrency,timestamp)
			VALUES %s`,

		strings.Join(queryValues, ","))

	valuePattern := "(?, ?, ?, ?, ?)"
	valuePattern += "," // mic on lar
	n := 0
	for strings.IndexByte(pgStatement, '?') != -1 {
		n++
		fieldReplaceValue := "$" + strconv.Itoa(n)
		pgStatement = strings.Replace(pgStatement, "?", fieldReplaceValue, 1)
	}
	pgStatement = strings.TrimSuffix(pgStatement, ",)")

	pgStatement += `ON CONFLICT ON CONSTRAINT fromunq DO UPDATE
			SET rate=EXCLUDED.rate, timestamp=EXCLUDED.timestamp, amount=EXCLUDED.amount
			WHERE currencydata.fromcurrency = EXCLUDED.fromcurrency AND currencydata.tocurrency=EXCLUDED.tocurrency`

	execErr := executeQuery(pgStatement, queryFields)
	if execErr != nil {
		logger.WithField("execute statement failed", execErr.Error()).Error(" Failed to execute query")
		return execErr

	}
	return nil
}

func executeQuery(query string, queryFieldValues []interface{}) error {
	db, dbErr := makeConnection()
	if dbErr != nil {
		logger.WithField("dbErr", dbErr.Error()).Error(" Failed to create table")
		return dbErr
	}
	createDbErr := createCurrencyTableIfNotExist()
	if createDbErr != nil {
		logger.WithField("createDbErr", createDbErr.Error()).Error(" Failed to create table")
		return createDbErr
	}
	preparedStatement, err := db.Prepare(query)
	if err != nil {
		logger.WithField("error in prepare query statement", err.Error()).Error("prepare query Failed")
		return err
	}

	result, err := preparedStatement.Exec(queryFieldValues...)
	if err != nil {
		logger.WithField("error in exec query: ", err.Error()).Error("Query Failed")
		return err
	}
	logger.WithField("Executed statement Result: ", result).Info("Output")
	return nil
}
