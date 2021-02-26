package handle

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Squirrel-Qiu/image-bed/client"
	"github.com/Squirrel-Qiu/image-bed/dbb"
	"github.com/Squirrel-Qiu/image-bed/store"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbInstance := dbb.InitDB(db)

	ts := gin.New()

	api := &Implement{
		DB:        dbInstance,
		Tool:
	}

	ts.GET("/get/:resourceId", api.Get)

	oParams := "b0804ec967"
	req, err := http.NewRequest("POST", "http://localhost/get/b0804ec967", bytes.NewBufferString(oParams))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	rows1 := sqlmock.NewRows([]string{"bucket"}).AddRow("2021-02")
	mock.ExpectQuery("SELECT bucket FROM resource WHERE id=?").
		WithArgs("b0804ec967").WillReturnRows(rows1)

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	file, err := os.Open("/tmp/1.txt")
	all, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, resp.Body, all)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

type testTake string

func (t *testTake) Take(resourceId, bucket string) (reader io.Reader, err error) {
	return nil
}
