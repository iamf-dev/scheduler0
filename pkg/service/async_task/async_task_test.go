package async_task

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"scheduler0/pkg/config"
	"scheduler0/pkg/db"
	"scheduler0/pkg/fsm"
	"scheduler0/pkg/models"
	"scheduler0/pkg/repository/async_task"
	"scheduler0/pkg/shared_repo"
	"testing"
	"time"
)

func Test_AsyncTaskManager_AddTasks(t *testing.T) {
	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	asyncTaskIds, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}

	assert.Equal(t, asyncTaskIds[0], uint64(1))
}

func Test_AsyncTaskManager_UpdateTasksById(t *testing.T) {
	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	output := "output"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	asyncTaskIds, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}
	updateErr := asyncTaskManager.UpdateTasksById(asyncTaskIds[0], models.AsyncTaskSuccess, output)
	if updateErr != nil {
		t.Fatal("failed update an async task", updateErr)
	}
	task, getTaskErr := asyncTaskManager.GetTaskWithRequestIdNonBlocking(requestId)
	if getTaskErr != nil {
		t.Fatal("failed to get an async task", getTaskErr)
	}
	assert.Equal(t, task.RequestId, requestId)
	assert.Equal(t, task.State, models.AsyncTaskSuccess)
}

func Test_AsyncTaskManager_UpdateTasksByRequestId(t *testing.T) {
	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	output := "output"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	_, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}
	updateErr := asyncTaskManager.UpdateTasksByRequestId(requestId, models.AsyncTaskSuccess, output)
	if updateErr != nil {
		t.Fatal("failed update an async task", updateErr)
	}
	task, getTaskErr := asyncTaskManager.GetTaskWithRequestIdNonBlocking(requestId)
	if getTaskErr != nil {
		t.Fatal("failed to get an async task", getTaskErr)
	}
	assert.Equal(t, task.RequestId, requestId)
	assert.Equal(t, task.State, models.AsyncTaskSuccess)
}

func Test_AsyncTaskManager_AddSubscriber(t *testing.T) {
	bctx := context.Background()
	ctx, cancler := context.WithCancel(bctx)
	defer cancler()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	output := "output"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	taskIds, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}

	executed := false
	asyncTaskManager.ListenForNotifications()

	subscriberId, addSubErr := asyncTaskManager.AddSubscriber(taskIds[0], func(task models.AsyncTask) {
		assert.Equal(t, task.RequestId, requestId)
		executed = true
	})
	if addSubErr != nil {
		t.Fatal("failed add a subscriber for an async task", addSubErr)
	}
	assert.Equal(t, subscriberId, uint64(1))
	updateErr := asyncTaskManager.UpdateTasksByRequestId(requestId, models.AsyncTaskSuccess, output)
	if updateErr != nil {
		t.Fatal("failed update an async task", updateErr)
	}
	time.Sleep(time.Second * 1)
	assert.Equal(t, executed, true)
	cancler()
}

func Test_AsyncTaskManager_DeleteSubscriber(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	output := "output"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	taskIds, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}

	executed := false
	asyncTaskManager.ListenForNotifications()

	subscriberId, addSubErr := asyncTaskManager.AddSubscriber(taskIds[0], func(task models.AsyncTask) {
		assert.Equal(t, task.RequestId, requestId)
		executed = true
	})
	if addSubErr != nil {
		t.Fatal("failed add a subscriber for an async task", addSubErr)
	}
	assert.Equal(t, subscriberId, uint64(1))
	deleteSubErr := asyncTaskManager.DeleteSubscriber(taskIds[0], subscriberId)
	if deleteSubErr != nil {
		t.Fatal("failed delete a subscriber for an async task", deleteSubErr)
	}
	updateErr := asyncTaskManager.UpdateTasksByRequestId(requestId, models.AsyncTaskSuccess, output)
	if updateErr != nil {
		t.Fatal("failed update an async task", updateErr)
	}
	time.Sleep(time.Second * 1)
	assert.Equal(t, executed, false)
}

func Test_AsyncTaskManager_GetTaskBlocking(t *testing.T) {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	output := "output"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	taskIds, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}

	asyncTaskManager.ListenForNotifications()
	tashCh, size, getTaskBlErr := asyncTaskManager.GetTaskBlocking(taskIds[0])
	if getTaskBlErr != nil {
		t.Fatal("failed update an async task", getTaskBlErr)
	}
	updateErr := asyncTaskManager.UpdateTasksByRequestId(requestId, models.AsyncTaskSuccess, output)
	if updateErr != nil {
		t.Fatal("failed update an async task", updateErr)
	}
	assert.Equal(t, size, uint64(1))
	task := <-tashCh
	assert.Equal(t, task.RequestId, requestId)
	assert.Equal(t, task.State, models.AsyncTaskSuccess)
	assert.Equal(t, task.Output, output)
}

func Test_AsyncTaskManager_GetTaskWithRequestIdBlocking(t *testing.T) {
	ctx, canceler := context.WithCancel(context.Background())
	defer canceler()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	output := "output"

	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	_, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}

	asyncTaskManager.ListenForNotifications()
	tashCh, _, getTaskBlErr := asyncTaskManager.GetTaskWithRequestIdBlocking(requestId)
	if getTaskBlErr != nil {
		t.Fatal("failed update an async task", getTaskBlErr)
	}
	updateErr := asyncTaskManager.UpdateTasksByRequestId(requestId, models.AsyncTaskSuccess, output)
	if updateErr != nil {
		t.Fatal("failed update an async task", updateErr)
	}
	//assert.Equal(t, size, 1)
	task := <-tashCh
	assert.Equal(t, task.RequestId, requestId)
	assert.Equal(t, task.State, models.AsyncTaskSuccess)
	assert.Equal(t, task.Output, output)
}

func Test_AsyncTaskManager_GetTaskIdWithRequestId(t *testing.T) {
	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	_, getErr := asyncTaskManager.AddTasks(input, requestId, service)
	if getErr != nil {
		t.Fatal("failed add an async task", getErr)
	}

	taskId, getTaskBlErr := asyncTaskManager.GetTaskIdWithRequestId(requestId)
	if getTaskBlErr != nil {
		t.Fatal("failed update an async task", getTaskBlErr)
	}
	assert.Equal(t, taskId, uint64(1))
}

func Test_AsyncTaskManager_GetUnCommittedTasks(t *testing.T) {
	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	input := "{'a':2}"
	requestId := "request-id"
	service := "asyncService"
	// Create a mock raft cluster
	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	_, batchAddErr := asyncTaskManagerRepo.BatchInsert([]models.AsyncTask{
		{
			RequestId: requestId,
			Service:   service,
			Input:     input,
		},
	}, false)
	if batchAddErr != nil {
		t.Fatal("failed batch add an async task", batchAddErr)
	}
	uncommittedTasks, getUTerr := asyncTaskManager.GetUnCommittedTasks()
	if getUTerr != nil {
		t.Fatal("failed batch add an async task", getUTerr)
	}
	assert.Equal(t, uncommittedTasks[0].RequestId, requestId)
	assert.Equal(t, uncommittedTasks[0].Input, input)
	assert.Equal(t, uncommittedTasks[0].Service, service)
}

func TestAsyncTaskService_SetSingleNodeMode(t *testing.T) {
	ctx := context.Background()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "async-task-manager-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
	scheduler0config := config.NewScheduler0Config()
	sharedRepo := shared_repo.NewSharedRepo(logger, scheduler0config)
	scheduler0RaftActions := fsm.NewScheduler0RaftActions(sharedRepo, nil)
	tempFile, err := ioutil.TempFile("", "test-db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	sqliteDb := db.NewSqliteDbConnection(logger, tempFile.Name())
	sqliteDb.RunMigration()
	sqliteDb.OpenConnectionToExistingDB()
	scheduler0Store := fsm.NewFSMStore(logger, scheduler0RaftActions, scheduler0config, sqliteDb, nil, nil, nil, nil, nil)

	cluster := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:          1,
		Bootstrap:      true,
		Conf:           raft.DefaultConfig(),
		ConfigStoreFSM: false,
		MakeFSMFunc: func() raft.FSM {
			return scheduler0Store.GetFSM()
		},
	})
	defer cluster.Close()
	cluster.FullyConnect()
	scheduler0Store.UpdateRaft(cluster.Leader())

	asyncTaskManagerRepo := async_task.NewAsyncTasksRepo(ctx, logger, scheduler0RaftActions, scheduler0Store)
	asyncTaskManager := NewAsyncTaskManager(ctx, logger, scheduler0Store, asyncTaskManagerRepo, scheduler0config)

	assert.Equal(t, asyncTaskManager.GetSingleNodeMode(), false)

	asyncTaskManager.SetSingleNodeMode(true)

	assert.Equal(t, asyncTaskManager.GetSingleNodeMode(), true)
}
