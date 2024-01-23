package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/anirudhsudhir/Bingo/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	snipModel   *models.SnipModel
}

func main() {
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	defaultDsn, err := parseDBCredentials()
	if err != nil {
		errorLogger.Fatal(err)
	}
	dsn := flag.String("dsn", defaultDsn, "MySQL data source name")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		snipModel:   &models.SnipModel{DB: db},
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLogger,
		Handler:  app.routes(),
	}

	infoLogger.Printf("Listening on port %s", *addr)
	err = server.ListenAndServe()
	if err != nil {
		errorLogger.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func parseDBCredentials() (dsn string, err error) {
	file, err := os.Open("secrets.env")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var secrets string
	for scanner.Scan() {
		secrets += scanner.Text()
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}

	reUser := regexp.MustCompile(`(DBuser = "\w+")`)
	rePass := regexp.MustCompile(`(DBpass = "\w+")`)
	dbUser, _ := strings.CutPrefix(reUser.FindString(secrets), "DBuser = \"")
	dbPass, _ := strings.CutPrefix(rePass.FindString(secrets), "DBpass = \"")
	dbUser, _ = strings.CutSuffix(dbUser, "\"")
	dbPass, _ = strings.CutSuffix(dbPass, "\"")

	if dbUser == "" {
		return "", errors.New("no database user present in secrets.env")
	}
	if dbPass == "" {
		return "", errors.New("no database password present in secrets.env")
	}

	dsn = fmt.Sprintf("%s:%s@/bingo?parseTime=true", dbUser, dbPass)
	return dsn, nil
}
