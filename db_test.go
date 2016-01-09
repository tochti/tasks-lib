package tasks

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/tochti/gin-gum/gumauth"
	"github.com/tochti/gin-gum/gumspecs"
	"github.com/tochti/gin-gum/gumtest"

	"gopkg.in/gorp.v1"
)

const (
	TestDatabase = "testing"
)

func setenvMySQL() {
	os.Clearenv()
	os.Setenv("MYSQL_USER", "tochti")
	os.Setenv("MYSQL_PASSWORD", "123")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB_NAME", TestDatabase)
}

func initTestDB(t *testing.T) *gorp.DbMap {
	setenvMySQL()

	mysql := gumspecs.ReadMySQL()

	sqlDB, err := mysql.DB()
	if err != nil {
		t.Fatal(err)
	}

	db := GorpInit(sqlDB)

	err = db.DropTablesIfExists()
	if err != nil {
		t.Fatal(err)
	}

	err = db.CreateTables()
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func fillTestDB(t *testing.T, db *gorp.DbMap) ([]*Task, []*gumauth.User, []*TasksUsers) {
	users := []*gumauth.User{
		{
			ID:        1,
			Username:  "robot",
			FirstName: "Mr.",
			LastName:  "Robot",
		},
		{
			ID:        2,
			Username:  "elliot",
			FirstName: "Elliot",
			LastName:  "Alderson",
		},
	}

	tasks := []*Task{
		{
			1,
			"Infect evel corp",
			"Buy raspberry pi 2",
			gumtest.SimpleNow().Add(1 * time.Hour),
			false,
		},
		{
			2,
			"Make fsociety video",
			"Don't forget the crazy masks",
			gumtest.SimpleNow().Add(1 * time.Hour),
			false,
		},
		{
			3,
			"Take some morphine",
			"",
			gumtest.SimpleNow().Add(-1 * time.Hour),
			true,
		},
		{
			4,
			"Kiss my own sister",
			"not!!",
			gumtest.SimpleNow().Add(1 * time.Hour),
			false,
		},
	}

	tasksUsers := []*TasksUsers{
		{1, 1},
		{1, 2},
		{2, 3},
		{2, 4},
	}

	err := db.Insert(gumtest.IfaceSlice(tasks)...)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Insert(gumtest.IfaceSlice(users)...)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Insert(gumtest.IfaceSlice(tasksUsers)...)
	if err != nil {
		t.Fatal(err)
	}

	return tasks, users, tasksUsers
}

func Test_ReadAllTasksUser(t *testing.T) {
	db := initTestDB(t)
	tasks, _, _ := fillTestDB(t, db)

	tasksResult, err := ReadAllTasksUser(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	expectTasks := []Task{
		*tasks[0],
		*tasks[1],
	}

	for i, task := range expectTasks {
		if reflect.DeepEqual(task, tasksResult[i]) {
			t.Fatalf("Expect %v was %v", task, tasksResult[i])
		}
	}

}

func Test_ReadOneTaskUser(t *testing.T) {
	db := initTestDB(t)
	tasks, _, _ := fillTestDB(t, db)

	task, err := ReadOneTaskUser(db, 1, 1)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(*tasks[0], task) {
		t.Fatalf("Expect %v was %v", tasks[0], task)
	}
}

func Test_CreateTaskUser(t *testing.T) {
	db := initTestDB(t)

	task := &Task{
		Name: "Build botnetwork and bring down the world",
		Desc: "Search for help at the Chinese",
		Done: false,
	}

	err := CreateTaskUser(db, task, 1)
	if err != nil {
		t.Fatal(err)
	}

	if task.ID == 0 {
		t.Fatal("Expect id gt then 0")
	}
}

func Test_DeleteTaskUser(t *testing.T) {
	db := initTestDB(t)
	tasks, _, _ := fillTestDB(t, db)

	userID := int64(1)
	err := DeleteTaskUser(db, tasks[0], userID)
	if err != nil {
		t.Fatal(err)
	}

	tasksResult, err := ReadAllTasksUser(db, userID)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasksResult) != 1 {
		t.Fatal("Expect 1 task in list found", len(tasks))
	}

	if reflect.DeepEqual(*tasks[1], tasksResult[0]) {
		t.Fatalf("Expect %v was %v", *tasks[1], tasksResult[0])
	}

}
