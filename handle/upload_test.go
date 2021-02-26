package handle

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Squirrel-Qiu/image-bed/dbb"
)

func TestUpload(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbInstance := dbb.InitDB(db)

	ts := gin.New()

	file, err := os.Open("/tmp/1.txt")
	api := &Implement{
		DB:        dbInstance,
		Generator: testGenerator("b0804ec967"),
		Tool: testStorage(""),
	}

	ts.POST("/upload", api.Upload)

	req, err := http.NewRequest("POST", "http://localhost/upload", file)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	mock.ExpectQuery("SELECT id FROM resource WHERE hash=?").
		WithArgs("a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO resource").WithArgs("b0804ec967", "2021-02",
			"a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447", "2021-02-24 15:04", 41).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Body, "http://localhost/b0804ec967")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

type testGenerator string

func (t testGenerator) GenerateId(idType string) (id string, err error) {
	return string(t), nil
}

type testStorage string

func (t *testStorage) Storage(resourceId, bucket string, reader *bytes.Reader) error {
	return nil
}
