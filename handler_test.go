package tasks

import (
	"encoding/json"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/gin-gum/gumwrap"
)

func Test_ReadAllHandler(t *testing.T) {
	db := initTestDB(t)
	tasks, _, _ := fillTestDB(t, db)

	r := gin.New()
	r.GET("/", gumtest.MockAuther(gumwrap.Gorp(ReadAll, db), "1"))

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/", "")

	expectResp := gumtest.JSONResponse{200, []Task{*tasks[0], *tasks[1]}}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_ReadOneHandler(t *testing.T) {
	db := initTestDB(t)
	tasks, _, _ := fillTestDB(t, db)

	r := gin.New()
	r.GET("/:id", gumtest.MockAuther(gumwrap.Gorp(ReadOne, db), "1"))

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/2", "")

	expectResp := gumtest.JSONResponse{200, *tasks[1]}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_CreateHandler(t *testing.T) {
	db := initTestDB(t)

	r := gin.New()
	r.POST("/", gumtest.MockAuther(gumwrap.Gorp(Create, db), "1"))

	task := Task{
		Name: "Help to breakout some fucking candyman out of prison",
	}

	body, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))

	task.ID = 1
	expectResp := gumtest.JSONResponse{201, task}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}
