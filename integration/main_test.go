package integration

import (
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Irori235/system-design-2023-v2/internal/handler"
	"github.com/Irori235/system-design-2023-v2/internal/migration"
	"github.com/Irori235/system-design-2023-v2/internal/pkg/config"
	"github.com/Irori235/system-design-2023-v2/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

var (
	db        *sqlx.DB
	engine    *gin.Engine
	r         *repository.Repository
	h         *handler.Handler
	userIDMap = make(map[string]uuid.UUID)
	taskMap   = make(map[string]handler.GetTaskResponse)
	jwtMap    = make(map[string]string)
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("connect to docker: ", err)
	}

	if err := pool.Client.Ping(); err != nil {
		log.Fatal("ping docker: ", err)
	}

	mysqlConfig := config.MySQL()

	resource, err := pool.Run("mysql", "latest", []string{
		"MYSQL_ROOT_PASSWORD=" + mysqlConfig.Passwd,
		"MYSQL_DATABASE=" + mysqlConfig.DBName,
	})
	if err != nil {
		log.Fatal("run docker: ", err)
	}

	mysqlConfig.Addr = "localhost:" + resource.GetPort("3306/tcp")

	if err := pool.Retry(func() error {
		_db, err := sqlx.Connect("mysql", mysqlConfig.FormatDSN())
		if err != nil {
			return err
		}
		db = _db
		return _db.Ping()
	}); err != nil {
		log.Fatal("connect to database container: ", err)
	}

	// migrate tables
	if err := migration.MigrateTables(db.DB); err != nil {
		log.Fatal("migrate tables: ", err)
	}

	// setup dependencies
	r = repository.New(db)
	h = handler.New(r)
	engine = gin.New()
	engine.Use(gin.Recovery())
	// engine.Use(gin.Logger())

	h.SetupRoutes(engine.Group("/api/v1"))

	log.Println("start integration test")
	m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatal("purge docker: ", err)
	}
}

func doRequest(t *testing.T, method, path string, bodystr string, headers ...map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(bodystr))
	req.Header.Set("Content-Type", "application/json")

	if len(headers) > 0 {
		header := headers[0]
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	return rec
}

func assert(t *testing.T, expected any, actual any) {
	t.Helper()

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("diff: %v", diff)
		t.Fail()
	}
}
