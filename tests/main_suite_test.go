package tests

import (
	"avito-tech-go-task/internal/application/service"
	"avito-tech-go-task/internal/clients/postgres"
	"avito-tech-go-task/internal/infrastructure/http/controller"
	"avito-tech-go-task/internal/infrastructure/storage"
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	dbDSN string
)

type TestSuite struct {
	suite.Suite
	db *postgres.Client
	*controller.ApiService
}

func init() {
	dbDSN = os.Getenv("DSN")
}

func TestSuiteFunc(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	db, err := postgres.Connect(dbDSN)
	if err != nil {
		s.FailNow("failed to connect", err)
	}
	s.db = db
	s.initDeps()
	if err = populateDB(s.db); err != nil {
		s.FailNow("failed to populate db", err)
	}
}

func (s *TestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *TestSuite) initDeps() {
	team := storage.NewTeamRepo(s.db)
	user := storage.NewUserRepo(s.db)
	pr := storage.NewPRRepo(s.db)
	prService := service.NewPRService(pr, user, team)
	s.ApiService = controller.NewApiService(prService)
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

// очистить таблицу
func truncateTable(db *postgres.Client, tableName string) error {
	sqlStatement := `TRUNCATE TABLE ` + tableName
	_, err := db.Exec(context.Background(), sqlStatement)
	if err != nil {
		return err
	}
	return nil
}

// заполение бд тестовыми данными
func populateDB(db *postgres.Client) (err error) {
	cleanDB(db)
	//teamRepo := storage.NewTeamRepo(db)
	//userRepo := storage.NewUserRepo(db)
	//prRepo := storage.NewPRRepo(db)

	//for _, team := range testTeams {
	//	err = teamRepo.Save(context.Background(), *domain.NewTeam(team.teamName), team.members)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

// очистить бд
func cleanDB(db *postgres.Client) {
	err := truncateTable(db, "users")
	if err != nil {
		log.Print("failed to truncate users", err)
	}

	err = truncateTable(db, "pull_requests")
	if err != nil {
		log.Print("failed to truncate pull_requests", err)
	}

	err = truncateTable(db, "teams")
	if err != nil {
		log.Print("failed to truncate teams", err)
	}

	err = truncateTable(db, "user_review_stats")
	if err != nil {
		log.Print("failed to truncate teams", err)
	}
}
