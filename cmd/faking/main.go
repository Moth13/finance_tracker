package main

import (
	"context"
	"log"
	"time"

	db "github.com/moth13/finance_tracker/db/sqlc"
	"github.com/moth13/finance_tracker/util"
	"github.com/shopspring/decimal"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}

	conn, err := db.CreateDBConnection(config.DBSource)
	defer conn.Close()

	store := db.NewStore(conn)

	argUser := db.CreateUserParams{
		Username:       "jose",
		HashedPassword: "secret",
		FullName:       "Jose Marcel",
		Email:          "jose.marcel@gmail.com",
		Currency:       util.EUR,
	}

	user, err := store.CreateUser(context.Background(), argUser)
	if err != nil {
		log.Fatal("Can't create user", err)
	}

	argAccount := db.CreateAccountParams{
		Owner:       user.Username,
		Title:       "BoursoBank",
		Description: "Main account",
	}

	argAccount.InitBalance, err = decimal.NewFromString("207.34")
	if err != nil {
		log.Fatal("Can't create Amount", err)
	}

	account, err := store.CreateAccount(context.Background(), argAccount)
	if err != nil {
		log.Fatal("Can't create account", err)
	}

	argYear := db.CreateYearParams{
		Title:       "2024",
		Owner:       user.Username,
		Description: "current year",
		StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	year, err := store.CreateYear(context.Background(), argYear)
	if err != nil {
		log.Fatal("Can't create year", err)
	}

	argMonth := db.CreateMonthParams{
		Title:       "December 2024",
		Owner:       user.Username,
		Description: "nowel",
		YearID:      year.ID,
		StartDate:   time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	month, err := store.CreateMonth(context.Background(), argMonth)
	if err != nil {
		log.Fatal("Can't create month", err)
	}

	argCategoryCourse := db.CreateCategoryParams{
		Title: "Course",
		Owner: user.Username,
	}

	catCourse, err := store.CreateCategory(context.Background(), argCategoryCourse)
	if err != nil {
		log.Fatal("Can't create category course", err)
	}

	argCategoryFun := db.CreateCategoryParams{
		Title: "Fun",
		Owner: user.Username,
	}

	catFun, err := store.CreateCategory(context.Background(), argCategoryFun)
	if err != nil {
		log.Fatal("Can't create category fun", err)
	}

	argCategoryAbo := db.CreateCategoryParams{
		Title: "Abonnement",
		Owner: user.Username,
	}

	catAbo, err := store.CreateCategory(context.Background(), argCategoryAbo)
	if err != nil {
		log.Fatal("Can't create category abo", err)
	}

	argCategorySalaire := db.CreateCategoryParams{
		Title: "Salaire",
		Owner: user.Username,
	}

	catSalaire, err := store.CreateCategory(context.Background(), argCategorySalaire)
	if err != nil {
		log.Fatal("Can't create category salaire", err)
	}

	argLine1 := db.AddLineTxParams{
		Title:       "SuperU",
		Owner:       user.Username,
		AccountID:   account.ID,
		MonthID:     month.ID,
		CategoryID:  catCourse.ID,
		YearID:      year.ID,
		DueDate:     time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC),
		Checked:     false,
		Description: "",
	}
	argLine1.Amount, err = decimal.NewFromString("-57.30")
	if err != nil {
		log.Fatal("Can't create Amount", err)
	}

	_, err = store.AddLineTx(context.Background(), argLine1)
	if err != nil {
		log.Fatal("Can't create line1", err)
	}

	argLine2 := db.AddLineTxParams{
		Title:       "Netflix",
		Owner:       user.Username,
		AccountID:   account.ID,
		MonthID:     month.ID,
		CategoryID:  catAbo.ID,
		YearID:      year.ID,
		DueDate:     time.Date(2024, 12, 14, 0, 0, 0, 0, time.UTC),
		Checked:     true,
		Description: "",
	}
	argLine2.Amount, err = decimal.NewFromString("-13.49")
	if err != nil {
		log.Fatal("Can't create Amount", err)
	}

	_, err = store.AddLineTx(context.Background(), argLine2)
	if err != nil {
		log.Fatal("Can't create line2", err)
	}

	argLine3 := db.AddLineTxParams{
		Title:       "Decathlon ski",
		Owner:       user.Username,
		AccountID:   account.ID,
		MonthID:     month.ID,
		CategoryID:  catFun.ID,
		YearID:      year.ID,
		DueDate:     time.Date(2024, 12, 28, 0, 0, 0, 0, time.UTC),
		Checked:     false,
		Description: "",
	}
	argLine3.Amount, err = decimal.NewFromString("-124")
	if err != nil {
		log.Fatal("Can't create Amount", err)
	}

	_, err = store.AddLineTx(context.Background(), argLine3)
	if err != nil {
		log.Fatal("Can't create line3", err)
	}

	argLine4 := db.AddLineTxParams{
		Title:       "Salaire",
		Owner:       user.Username,
		AccountID:   account.ID,
		MonthID:     month.ID,
		CategoryID:  catSalaire.ID,
		YearID:      year.ID,
		DueDate:     time.Date(2024, 12, 3, 0, 0, 0, 0, time.UTC),
		Checked:     false,
		Description: "",
	}
	argLine4.Amount, err = decimal.NewFromString("2124.98")
	if err != nil {
		log.Fatal("Can't create Amount", err)
	}

	_, err = store.AddLineTx(context.Background(), argLine4)
	if err != nil {
		log.Fatal("Can't create line4", err)
	}
}
