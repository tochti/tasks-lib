package tasks

import (
	"errors"
	"net/http"
	"strconv"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumrest"
	"github.com/tochti/session-store"
)

func ReadUserID(c *gin.Context) (int64, error) {
	tmp, ok := c.Get("Session")
	if !ok {
		return -1, errors.New("Missing session")
	}

	s := tmp.(s2tore.Session)

	userID, err := strconv.ParseInt(s.UserID(), 10, 64)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func ReadAll(ginCtx *gin.Context, db *gorp.DbMap) {
	userID, err := ReadUserID(ginCtx)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	tasks, err := ReadAllTasksUser(db, userID)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusNotFound, err)
		return
	}

	ginCtx.JSON(200, tasks)
}

func ReadOne(ginCtx *gin.Context, db *gorp.DbMap) {
	tmp := ginCtx.Param("id")
	taskID, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	userID, err := ReadUserID(ginCtx)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	task, err := ReadOneTaskUser(db, taskID, userID)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusNotFound, err)
		return
	}

	ginCtx.JSON(http.StatusOK, task)

}

func Create(ginCtx *gin.Context, db *gorp.DbMap) {
	userID, err := ReadUserID(ginCtx)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	task := &Task{}
	err = ginCtx.BindJSON(task)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	err = CreateTaskUser(db, task, userID)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusNotFound, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, task)
}

func Update(ginCtx *gin.Context, db *gorp.DbMap) {
	tmp := ginCtx.Params.ByName("id")
	taskID, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	newTask := Task{}
	err = ginCtx.BindJSON(&newTask)
	newTask.ID = taskID

	_, err = db.Update(&newTask)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(200, newTask)

}

func Delete(ginCtx *gin.Context, db *gorp.DbMap) {
	userID, err := ReadUserID(ginCtx)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	tmp := ginCtx.Params.ByName("id")
	taskID, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	err = DeleteTaskUser(db, &Task{ID: taskID}, userID)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, nil)
}
