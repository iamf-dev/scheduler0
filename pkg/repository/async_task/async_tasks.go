package async_task

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"scheduler0/pkg/constants"
	"scheduler0/pkg/fsm"
	"scheduler0/pkg/models"
	"scheduler0/pkg/scheduler0time"
	"scheduler0/pkg/utils"
	"time"
)

type asyncTasksRepo struct {
	context               context.Context
	fsmStore              fsm.Scheduler0RaftStore
	logger                hclog.Logger
	scheduler0RaftActions fsm.Scheduler0RaftActions
}

//go:generate mockery --name AsyncTasksRepo --output ../mocks
type AsyncTasksRepo interface {
	BatchInsert(tasks []models.AsyncTask, committed bool) ([]uint64, *utils.GenericError)
	RaftBatchInsert(tasks []models.AsyncTask, fromNodeId uint64) ([]uint64, *utils.GenericError)
	RaftUpdateTaskState(task models.AsyncTask, state models.AsyncTaskState, output string) *utils.GenericError
	UpdateTaskState(task models.AsyncTask, state models.AsyncTaskState, output string) *utils.GenericError
	GetTask(taskId uint64) (*models.AsyncTask, *utils.GenericError)
	GetAllTasks(committed bool) ([]models.AsyncTask, *utils.GenericError)
}

func NewAsyncTasksRepo(context context.Context, logger hclog.Logger, scheduler0RaftActions fsm.Scheduler0RaftActions, fsmStore fsm.Scheduler0RaftStore) AsyncTasksRepo {
	return &asyncTasksRepo{
		context:               context,
		logger:                logger.Named("async-task-repo"),
		fsmStore:              fsmStore,
		scheduler0RaftActions: scheduler0RaftActions,
	}
}

func (repo *asyncTasksRepo) BatchInsert(tasks []models.AsyncTask, committed bool) ([]uint64, *utils.GenericError) {
	repo.fsmStore.GetDataStore().ConnectionLock()
	defer repo.fsmStore.GetDataStore().ConnectionUnlock()

	batches := utils.Batch[models.AsyncTask](tasks, 5)
	results := make([]uint64, 0, len(tasks))

	schedulerTime := scheduler0time.GetSchedulerTime()
	now := schedulerTime.GetTime(time.Now())

	table := constants.CommittedAsyncTableName
	if !committed {
		table = constants.UnCommittedAsyncTableName
	}

	for _, batch := range batches {
		query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?)",
			table,
			constants.AsyncTasksRequestIdColumn,
			constants.AsyncTasksInputColumn,
			constants.AsyncTasksOutputColumn,
			constants.AsyncTasksStateColumn,
			constants.AsyncTasksServiceColumn,
			constants.AsyncTasksDateCreatedColumn,
		)
		params := []interface{}{
			batch[0].RequestId,
			batch[0].Input,
			batch[0].Output,
			0,
			batch[0].Service,
			now,
		}

		for _, row := range batch[1:] {
			query += ",(?, ?, ?, ?, ?, ?)"
			params = append(params, row.RequestId, row.Input, row.Output, 0, row.Service, now)
		}

		ids := make([]uint64, 0, len(batch))

		query += ";"

		res, err := repo.fsmStore.GetDataStore().GetOpenConnection().Exec(query, params...)
		if err != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
		}

		if res == nil {
			return nil, utils.HTTPGenericError(http.StatusServiceUnavailable, "service is unavailable")
		}

		lastInsertedId, err := res.LastInsertId()
		if err != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
		}
		for i := lastInsertedId - int64(len(batch)) + 1; i <= lastInsertedId; i++ {
			ids = append(ids, uint64(i))
		}

		results = append(results, ids...)
	}

	return results, nil
}

func (repo *asyncTasksRepo) RaftBatchInsert(tasks []models.AsyncTask, fromNodeId uint64) ([]uint64, *utils.GenericError) {
	batches := utils.Batch[models.AsyncTask](tasks, 6)
	results := make([]uint64, 0, len(tasks))
	schedulerTime := scheduler0time.GetSchedulerTime()
	now := schedulerTime.GetTime(time.Now())

	table := constants.CommittedAsyncTableName

	for _, batch := range batches {

		query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES (?, ?, ?, ?, ?, ?)",
			table,
			constants.AsyncTasksRequestIdColumn,
			constants.AsyncTasksInputColumn,
			constants.AsyncTasksOutputColumn,
			constants.AsyncTasksStateColumn,
			constants.AsyncTasksServiceColumn,
			constants.AsyncTasksDateCreatedColumn,
		)
		params := []interface{}{
			batch[0].RequestId,
			batch[0].Input,
			batch[0].Output,
			0,
			batch[0].Service,
			now,
		}
		for _, row := range batch[1:] {
			query += ",(?, ?, ?, ?, ?, ?)"
			params = append(params, row.RequestId, row.Input, row.Output, 0, row.Service, now)
		}

		ids := make([]uint64, 0, len(batch))

		query += ";"

		res, applyErr := repo.scheduler0RaftActions.WriteCommandToRaftLog(
			repo.fsmStore.GetRaft(),
			constants.CommandTypeDbExecute,
			query,
			params,
			[]uint64{fromNodeId},
			constants.CommandActionQueueJob)
		if applyErr != nil {
			return nil, applyErr
		}

		if res == nil {
			return nil, utils.HTTPGenericError(http.StatusServiceUnavailable, "service is unavailable")
		}

		lastInsertedId := uint64(res.Data.LastInsertedId)
		for i := lastInsertedId - uint64(len(batch)) + 1; i <= lastInsertedId; i++ {
			ids = append(ids, i)
		}

		results = append(results, ids...)
	}

	return results, nil
}

func (repo *asyncTasksRepo) RaftUpdateTaskState(task models.AsyncTask, state models.AsyncTaskState, output string) *utils.GenericError {
	updateQuery := sq.Update(constants.CommittedAsyncTableName).
		Set(constants.AsyncTasksStateColumn, state).
		Set(constants.AsyncTasksOutputColumn, output).
		Where(fmt.Sprintf("%s = ?", constants.AsyncTasksIdColumn), task.Id)

	query, params, err := updateQuery.ToSql()
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	_, applyErr := repo.scheduler0RaftActions.WriteCommandToRaftLog(repo.fsmStore.GetRaft(), constants.CommandTypeDbExecute, query, params, []uint64{}, 0)
	if err != nil {
		return applyErr
	}

	return nil
}

func (repo *asyncTasksRepo) UpdateTaskState(task models.AsyncTask, state models.AsyncTaskState, output string) *utils.GenericError {
	repo.fsmStore.GetDataStore().ConnectionLock()
	defer repo.fsmStore.GetDataStore().ConnectionUnlock()

	updateQuery := sq.Update(constants.UnCommittedAsyncTableName).
		Set(constants.AsyncTasksStateColumn, state).
		Set(constants.AsyncTasksOutputColumn, output).
		Where(fmt.Sprintf("%s = ?", constants.AsyncTasksIdColumn), task.Id)

	query, params, err := updateQuery.ToSql()
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}

	_, applyErr := repo.fsmStore.GetDataStore().GetOpenConnection().Exec(query, params...)
	if err != nil {
		return utils.HTTPGenericError(http.StatusInternalServerError, applyErr.Error())
	}

	return nil
}

func (repo *asyncTasksRepo) GetTask(taskId uint64) (*models.AsyncTask, *utils.GenericError) {
	repo.fsmStore.GetDataStore().ConnectionLock()
	defer repo.fsmStore.GetDataStore().ConnectionUnlock()

	query := fmt.Sprintf(
		"select %s, %s, %s, %s, %s, %s, %s from %s where %s = ? union select %s, %s, %s, %s, %s, %s, %s from %s where %s = ?",
		constants.AsyncTasksIdColumn,
		constants.AsyncTasksRequestIdColumn,
		constants.AsyncTasksInputColumn,
		constants.AsyncTasksOutputColumn,
		constants.AsyncTasksStateColumn,
		constants.AsyncTasksServiceColumn,
		constants.AsyncTasksDateCreatedColumn,
		constants.CommittedAsyncTableName,
		constants.AsyncTasksIdColumn,
		constants.AsyncTasksIdColumn,
		constants.AsyncTasksRequestIdColumn,
		constants.AsyncTasksInputColumn,
		constants.AsyncTasksOutputColumn,
		constants.AsyncTasksStateColumn,
		constants.AsyncTasksServiceColumn,
		constants.AsyncTasksDateCreatedColumn,
		constants.UnCommittedAsyncTableName,
		constants.AsyncTasksIdColumn,
	)

	rows, err := repo.fsmStore.GetDataStore().GetOpenConnection().Query(query, taskId, taskId)
	if err != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()
	var asyncTask models.AsyncTask
	for rows.Next() {
		scanErr := rows.Scan(
			&asyncTask.Id,
			&asyncTask.RequestId,
			&asyncTask.Input,
			&asyncTask.Output,
			&asyncTask.State,
			&asyncTask.Service,
			&asyncTask.DateCreated,
		)
		if scanErr != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, scanErr.Error())
		}
	}
	if rows.Err() != nil {
		return nil, utils.HTTPGenericError(http.StatusInternalServerError, rows.Err().Error())
	}
	return &asyncTask, nil
}

func (repo *asyncTasksRepo) countAsyncTasks(committed bool) uint64 {
	tableName := constants.UnCommittedAsyncTableName

	if committed {
		tableName = constants.CommittedAsyncTableName
	}

	selectBuilder := sq.Select("count(*)").
		From(tableName).
		RunWith(repo.fsmStore.GetDataStore().GetOpenConnection())

	rows, err := selectBuilder.Query()
	if err != nil {
		repo.logger.Error("failed to count async tasks rows", err)
		return 0
	}
	var count uint64 = 0
	for rows.Next() {
		scanErr := rows.Scan(&count)
		if err != nil {
			repo.logger.Error("failed to scan rows ", scanErr)
			return 0
		}
	}
	if rows.Err() != nil {
		repo.logger.Error("failed to count async tasks rows error", rows.Err())
		return 0
	}
	return count
}

func (repo *asyncTasksRepo) getAsyncTasksMinMaxIds(committed bool) (uint64, uint64) {
	tableName := constants.UnCommittedAsyncTableName

	if committed {
		tableName = constants.CommittedAsyncTableName
	}

	selectBuilder := sq.Select("min(id)", "max(id)").
		From(tableName).
		RunWith(repo.fsmStore.GetDataStore().GetOpenConnection())

	rows, err := selectBuilder.Query()
	if err != nil {
		repo.logger.Error("failed to count async tasks rows", err)
		return 0, 0
	}
	var minId uint64 = 0
	var maxId uint64 = 0
	for rows.Next() {
		scanErr := rows.Scan(&minId, &maxId)
		if err != nil {
			repo.logger.Error("failed to scan rows ", scanErr)
			return 0, 0
		}
	}
	if rows.Err() != nil {
		repo.logger.Error("failed to count async tasks rows error", rows.Err())
		return 0, 0
	}
	return minId, maxId
}

func (repo *asyncTasksRepo) GetAllTasks(committed bool) ([]models.AsyncTask, *utils.GenericError) {
	repo.fsmStore.GetDataStore().ConnectionLock()
	defer repo.fsmStore.GetDataStore().ConnectionUnlock()

	table := constants.CommittedAsyncTableName
	if !committed {
		table = constants.UnCommittedAsyncTableName
	}

	min, max := repo.getAsyncTasksMinMaxIds(committed)
	count := repo.countAsyncTasks(committed)
	results := make([]models.AsyncTask, 0, count)
	expandedIds := utils.ExpandIdsRange(min, max)

	batches := utils.Batch(expandedIds, 7)

	for _, batch := range batches {
		var params = []interface{}{batch[0]}
		var paramPlaceholders = "?"

		for _, b := range batch[1:] {
			paramPlaceholders += ",?"
			params = append(params, b)
		}

		query := fmt.Sprintf(
			"select %s, %s, %s, %s, %s, %s, %s from %s where id in (%s)",
			constants.AsyncTasksIdColumn,
			constants.AsyncTasksRequestIdColumn,
			constants.AsyncTasksInputColumn,
			constants.AsyncTasksOutputColumn,
			constants.AsyncTasksStateColumn,
			constants.AsyncTasksServiceColumn,
			constants.AsyncTasksDateCreatedColumn,
			table,
			paramPlaceholders,
		)
		rows, err := repo.fsmStore.GetDataStore().GetOpenConnection().Query(query, params...)
		if err != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, err.Error())
		}
		for rows.Next() {
			var asyncTask models.AsyncTask
			scanErr := rows.Scan(
				&asyncTask.Id,
				&asyncTask.RequestId,
				&asyncTask.Input,
				&asyncTask.Output,
				&asyncTask.State,
				&asyncTask.Service,
				&asyncTask.DateCreated,
			)
			if scanErr != nil {
				return nil, utils.HTTPGenericError(http.StatusInternalServerError, scanErr.Error())
			}
			results = append(results, asyncTask)
		}
		if rows.Err() != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, rows.Err().Error())
		}
		closeErr := rows.Close()
		if closeErr != nil {
			return nil, utils.HTTPGenericError(http.StatusInternalServerError, closeErr.Error())
		}
	}

	return results, nil
}
