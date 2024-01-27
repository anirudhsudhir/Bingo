package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/anirudhsudhir/Bingo/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

type application struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	snipModel     *models.SnipModel
	templateCache map[string]*template.Template
	formDecoder   *schema.Decoder
	sessionStore  *sessions.CookieStore
}

func main() {
	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	defaultDsn, sessionKey, err := parseSecrets(errorLogger)
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

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	formDecoder := schema.NewDecoder()
	app := &application{
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
		snipModel:     &models.SnipModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	app.sessionStore = store

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

func parseSecrets(errorLogger *log.Logger) (dsn, sessionLey string, err error) {
	inFile, err := os.Open("secrets.env")
	if err != nil {
		return "", "", err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	var secrets string
	for scanner.Scan() {
		secrets += scanner.Text()
	}
	if err = scanner.Err(); err != nil {
		return "", "", err
	}

	reUser := regexp.MustCompile(`DBuser = "(.*?)"`)
	rePass := regexp.MustCompile(`DBpass = "(.*?)"`)
	reSessionKey := regexp.MustCompile(`SessionKey = "(.*?)"`)
	dbUser, _ := strings.CutPrefix(reUser.FindString(secrets), "DBuser = \"")
	dbPass, _ := strings.CutPrefix(rePass.FindString(secrets), "DBpass = \"")
	sessionKey, _ := strings.CutPrefix(reSessionKey.FindString(secrets), "SessionKey = \"")
	dbUser, _ = strings.CutSuffix(dbUser, "\"")
	dbPass, _ = strings.CutSuffix(dbPass, "\"")
	sessionKey, _ = strings.CutSuffix(sessionKey, "\"")

	if dbUser == "" {
		return "", "", errors.New("no database user present in secrets.env")
	}
	if dbPass == "" {
		return "", "", errors.New("no database password present in secrets.env")
	}
	if sessionKey == "" {
		sessionKey, err = generateSessionKey(32)
		if err != nil {
			errorLogger.Fatal(err)
		}

		outFile, err := os.OpenFile("secrets.env", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			errorLogger.Fatal(err)
		}
		defer outFile.Close()

		key := fmt.Sprintf("\nSessionKey = \"%s\"", sessionKey)
		_, err = outFile.WriteString(key)
		if err != nil {
			errorLogger.Fatal(err)
		}
	}
	dsn = fmt.Sprintf("%s:%s@/bingo?parseTime=true", dbUser, dbPass)
	return dsn, sessionKey, nil
}
