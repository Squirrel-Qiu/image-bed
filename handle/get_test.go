package handle

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"

	"github.com/Squirrel-Qiu/image-bed/dbb"
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
		DB:   dbInstance,
		Tool: &testTool{},
	}

	ts.GET("/get/:resourceId", api.Get)

	oParams := "b0804ec967"
	req, err := http.NewRequest("POST", "https://localhost/get/b0804ec967", bytes.NewBufferString(oParams))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	//rows1 := sqlmock.NewRows([]string{"bucket"}).AddRow(nil)
	/*mock.ExpectQuery("SELECT bucket FROM resource WHERE id=?").
		WithArgs(oParams).WillReturnRows(rows1)*/

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	//assert.Equal(t, resp.Body, bytes.NewBuffer(all))

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func (t *testTool) Take(resourceId, bucket string) (reader io.Reader, err error) {
	_, err = os.Open("/tmp/1.txt")
	if err != nil {
		panic(err)
	}
	return nil, nil
}
