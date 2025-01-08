package controllers

import (
	"api-tarefas/models"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	DB *sql.DB
}

type ResponseErr struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

/*
Description: Lista todas as tarefas cadastradas no banco de dados

	@Tags Tarefas
	@Accept json
	@Produce json
	@Success 200 {object} []Tasks
	@Failure 400 {object} ResponseErr
	@Failure 500 {object} ResponseErr
	@Router /tarefas [get]
	@Param titulo query string false "Título da tarefa"
*/
func (tc *TaskController) GetTasks(ctx *gin.Context) {
	query := "select * from Tasks"
	rows, err := tc.DB.Query(query)

	if err != nil {
		response := ResponseErr{
			Error:   "Invalid request data.",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	defer rows.Close()

	var tasks []models.Tasks
	for rows.Next() {
		var task models.Tasks
		if err := rows.Scan(&task.ID, &task.Title, &task.Descricao, &task.Status); err != nil {
			response := ResponseErr{
				Error:   "Error when retrieving tasks from the database",
				Details: err.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
		tasks = append(tasks, task)
	}

	ctx.JSON(http.StatusOK, tasks)
}

/*
	 Busca uma tarefa cadastrada no banco de dados pelo ID

		@Tags Tarefas
		@Accept json
		@Produce json
		@Success 200 {object} Tasks
		@Failure 400 {object} ResponseErr
		@Failure 404 {object} ResponseErr
		@Failure 500 {object} ResponseErr
		@Router /tarefas/{id} [get]
		@Param id path int true "ID da tarefa"
*/
func (tc *TaskController) GetTask(ctx *gin.Context) {
	var task models.Tasks

	id := ctx.Param("id")
	query := "select * from Tasks where id = $1"

	rows, err := tc.DB.Query(query, id)
	if err != nil {
		response := ResponseErr{
			Error:   "Invalid request data.",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&task.ID, &task.Title, &task.Descricao, &task.Status); err != nil {
			response := ResponseErr{
				Error:   "Error when retrieving task from the database",
				Details: err.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	if task.ID == 0 {
		response := ResponseErr{
			Error: "Task not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, task)
}

/*
*
Description Busca tarefas cadastradas no banco de dados pelo título

	@Tags Tarefas
	@Accept json
	@Produce json
	@Success 200 {object} []Tasks
	@Failure 400 {object} ResponseErr
	@Failure 404 {object} ResponseErr
	@Failure 500 {object} ResponseErr
	@Router /tarefas/buscar [get]
	@Param titulo query string true "Título da tarefa"
*/
func (tc *TaskController) FindTasks(ctx *gin.Context) {
	title := ctx.Query("titulo")
	var tasks []models.Tasks

	if strings.TrimSpace(title) == "" {
		response := ResponseErr{
			Error:   "Invalid data",
			Details: "The 'titulo' query parameter is empty or was not filled.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	query := "SELECT * FROM Tasks WHERE title ILIKE $1"
	rows, err := tc.DB.Query(query, "%"+title+"%")

	if err != nil {
		response := ResponseErr{
			Error:   "Error executing query",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Tasks
		if err := rows.Scan(&task.ID, &task.Title, &task.Descricao, &task.Status); err != nil {
			response := ResponseErr{
				Error:   "Error when scanning tasks from the database",
				Details: err.Error(),
			}
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		response := ResponseErr{
			Error: "No tasks found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

/*
*
@Summary Cria uma nova tarefa
@Description Cria uma nova tarefa no banco de dados
@Tags Tarefas
@Accept json
@Produce json
@Success 201 {object} gin.H
@Failure 400 {object} ResponseErr
@Failure 500 {object} ResponseErr
@Router /tarefas [post]
@Param task body models.Tasks true "Dados da tarefa"
*/
func (tc *TaskController) CreateTask(ctx *gin.Context) {
	var task models.Tasks

	if err := ctx.ShouldBindJSON(&task); err != nil {
		response := ResponseErr{
			Error:   "Invalid data",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	query := "INSERT INTO Tasks (title, descricao, status) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := tc.DB.QueryRow(query, task.Title, task.Descricao, task.Status).Scan(&id)
	if err != nil {
		response := ResponseErr{
			Error:   "Error creating task",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	task.ID = id
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    task,
	})
}

/*
*
@Summary Atualiza uma tarefa existente
@Description Atualiza os dados de uma tarefa no banco de dados
@Tags Tarefas
@Accept json
@Produce json
@Success 200 {object} gin.H
@Failure 400 {object} ResponseErr
@Failure 404 {object} ResponseErr
@Failure 500 {object} ResponseErr
@Router /tarefas/{id} [put]
@Param id path int true "ID da tarefa"
@Param task body models.Tasks true "Dados da tarefa"
*/
func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var task models.Tasks

	if err := ctx.ShouldBindJSON(&task); err != nil {
		response := ResponseErr{
			Error:   "Invalid data",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	query := "UPDATE Tasks SET title = $1, descricao = $2, status = $3 WHERE id = $4"
	result, err := tc.DB.Exec(query, task.Title, task.Descricao, task.Status, id)

	if err != nil {
		response := ResponseErr{
			Error:   "Error updating task",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := ResponseErr{
			Error:   "Error retrieving affected rows",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if rowsAffected == 0 {
		response := ResponseErr{
			Error: "Task not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

/*
*
@Summary Exclui uma tarefa existente
@Description Exclui uma tarefa do banco de dados
@Tags Tarefas
@Accept json
@Produce json
@Success 200 {object} gin.H
@Failure 400 {object} ResponseErr
@Failure 404 {object} ResponseErr
@Failure 500 {object} ResponseErr
@Router /tarefas/{id} [delete]
@Param id path int true "ID da tarefa"
*/
func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	query := "DELETE FROM Tasks WHERE id = $1"
	result, err := tc.DB.Exec(query, id)

	if err != nil {
		response := ResponseErr{
			Error:   "Error deleting task",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := ResponseErr{
			Error:   "Error retrieving affected rows",
			Details: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if rowsAffected == 0 {
		response := ResponseErr{
			Error: "Task not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
