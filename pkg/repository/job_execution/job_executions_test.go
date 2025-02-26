package job_execution

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"scheduler0/pkg/config"
	"scheduler0/pkg/db"
	"scheduler0/pkg/fsm"
	"scheduler0/pkg/models"
	"scheduler0/pkg/repository/job"
	project_repo "scheduler0/pkg/repository/project"
	"scheduler0/pkg/shared_repo"
	"testing"
	"time"
)

func Test_JobExecutionsRepo_BatchInsert(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Define test data
	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         1,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         1,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	// Create a project to delete
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogScheduleState
	jobQueueVersion := uint64(1)
	jobExecutionVersions := map[uint64]uint64{
		1: 1,
		2: 2,
	}

	// Call the BatchInsert method
	jobExecutionsRepo.BatchInsert(jobs, nodeID, state, jobQueueVersion, jobExecutionVersions)

	// Retrieve the inserted job execution logs from the database
	query := fmt.Sprintf("SELECT * FROM %s;", ExecutionsUnCommittedTableName)
	rows, err := scheduler0Store.GetDataStore().GetOpenConnection().Query(query)
	if err != nil {
		t.Fatalf("Failed to query job execution logs: %v", err)
	}
	defer rows.Close()

	// Assert the number of retrieved job execution logs
	expectedCount := len(jobs)
	count := 0

	for rows.Next() {
		count++
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("Failed to iterate over job execution logs: %v", err)
	}

	assert.Equal(t, expectedCount, count)
}

func Test_JobExecutionsRepo_GetLastExecutionLogForJobIds(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Define test data
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogScheduleState
	jobQueueVersion := uint64(1)
	jobExecutionVersions := map[uint64]uint64{
		1: 1,
		2: 2,
	}

	// Call the BatchInsert method
	jobExecutionsRepo.BatchInsert(jobs, nodeID, state, jobQueueVersion, jobExecutionVersions)

	// Call the GetLastExecutionLogForJobIds method
	jobExecutionLogs := jobExecutionsRepo.GetLastExecutionLogForJobIds([]uint64{1, 2})

	// Assert the number of retrieved job execution logs
	expectedCount := len(jobs)
	count := len(jobExecutionLogs)
	assert.Equal(t, expectedCount, count)

	// Assert the correctness of the retrieved job execution logs
	for _, log := range jobExecutionLogs {
		assert.Contains(t, []uint64{1, 2}, log.JobId)
	}
}

func Test_JobExecutionsRepo_CountLastFailedExecutionLogs(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Define test data
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogFailedState
	jobQueueVersion := uint64(1)
	jobExecutionVersions := map[uint64]uint64{
		1: 1,
		2: 2,
	}

	// Call the BatchInsert method
	jobExecutionsRepo.BatchInsert(jobs, nodeID, state, jobQueueVersion, jobExecutionVersions)

	// Call the CountLastFailedExecutionLogs method
	count := jobExecutionsRepo.CountLastFailedExecutionLogs(jobs[0].ID, nodeID, 1)

	// Assert the count of last failed execution logs
	expectedCount := 1
	assert.Equal(t, uint64(expectedCount), count)
}

func Test_JobExecutionsRepo_CountExecutionLogs(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Define test data
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogScheduleState
	jobQueueVersion := uint64(1)
	jobExecutionVersions := map[uint64]uint64{
		1: 1,
		2: 2,
	}

	// Call the BatchInsert method
	jobExecutionsRepo.BatchInsert(jobs, nodeID, state, jobQueueVersion, jobExecutionVersions)

	// Call the CountExecutionLogs method for uncommitted execution logs
	uncommittedCount := jobExecutionsRepo.CountExecutionLogs(false)

	assert.Equal(t, uint64(len(jobs)), uncommittedCount)
}

func Test_JobExecutionsRepo_GetUncommittedExecutionsLogForNode(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Define test data
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogScheduleState
	jobQueueVersion := uint64(1)
	jobExecutionVersions := map[uint64]uint64{
		1: 1,
		2: 2,
	}

	// Call the BatchInsert method
	jobExecutionsRepo.BatchInsert(jobs, nodeID, state, jobQueueVersion, jobExecutionVersions)

	// Call the GetUncommittedExecutionsLogForNode method
	executionLogs := jobExecutionsRepo.GetUncommittedExecutionsLogForNode(nodeID)

	// Assert the number of retrieved execution logs
	expectedCount := len(jobs)
	count := len(executionLogs)
	assert.Equal(t, expectedCount, count)

	// Assert the correctness of the retrieved execution logs
	for _, log := range executionLogs {
		assert.Equal(t, nodeID, log.NodeId)
	}
}

func Test_LogJobExecutionStateInRaft(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Define test data
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogScheduleState
	jobQueueVersion := uint64(1)
	jobExecutionVersions := map[uint64]uint64{
		1: 1,
		2: 2,
	}

	jobExecutionsRepo.LogJobExecutionStateInRaft(jobs, state, jobExecutionVersions, jobQueueVersion, nodeID)

	// Call the GetLastExecutionLogForJobIds method
	jobExecutionLogs := jobExecutionsRepo.GetLastExecutionLogForJobIds([]uint64{1, 2})

	// Assert the number of retrieved job execution logs
	expectedCount := len(jobs)
	count := len(jobExecutionLogs)
	assert.Equal(t, expectedCount, count)

	// Assert the correctness of the retrieved job execution logs
	for _, log := range jobExecutionLogs {
		assert.Contains(t, []uint64{1, 2}, log.JobId)
	}
}

func Test_RaftInsertExecutionLogs(t *testing.T) {
	scheduler0config := config.NewScheduler0Config()
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "job-executions-repo-test",
		Level: hclog.LevelFromString("DEBUG"),
	})
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

	// Create a new JobExecutionsRepo instance
	jobExecutionsRepo := NewExecutionsRepo(logger, scheduler0RaftActions, scheduler0Store)

	// Create a JobRepo instance
	jobRepo := job.NewJobRepo(logger, scheduler0RaftActions, scheduler0Store)

	projectRepo := project_repo.NewProjectRepo(logger, scheduler0RaftActions, scheduler0Store, jobRepo)

	// Define test data
	project := models.Project{
		ID:          1,
		Name:        "Test Project",
		Description: "Test project description",
	}

	// Insert the project into the database
	_, pcreateErr := projectRepo.CreateOne(&project)
	if pcreateErr != nil {
		t.Fatal("failed to create project:", pcreateErr)
	}

	jobs := []models.Job{
		{
			ID:                1,
			ExecutionId:       "1",
			Spec:              "*/5 * * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-time.Hour),
			DateCreated:       time.Now(),
		},
		{
			ID:                2,
			ExecutionId:       "2",
			Spec:              "0 0 * * *",
			ProjectID:         project.ID,
			LastExecutionDate: time.Now().Add(-2 * time.Hour),
			DateCreated:       time.Now(),
		},
	}

	_, insertErr := jobRepo.BatchInsertJobs(jobs)
	if insertErr != nil {
		t.Fatal("failed to insert job", insertErr)
	}

	nodeID := uint64(1)
	state := models.ExecutionLogScheduleState
	jobQueueVersion := uint64(1)

	executionLogs := []models.JobExecutionLog{
		{
			UniqueId:              "one",
			State:                 state,
			NodeId:                nodeID,
			LastExecutionDatetime: time.Now(),
			NextExecutionDatetime: time.Now(),
			JobId:                 1,
			JobQueueVersion:       jobQueueVersion,
			ExecutionVersion:      1,
			DataCreated:           time.Now(),
		},
		{
			UniqueId:              "two",
			State:                 state,
			NodeId:                nodeID,
			LastExecutionDatetime: time.Now(),
			NextExecutionDatetime: time.Now(),
			JobId:                 2,
			JobQueueVersion:       jobQueueVersion,
			ExecutionVersion:      1,
			DataCreated:           time.Now(),
		},
	}

	jobExecutionsRepo.RaftInsertExecutionLogs(executionLogs, 1)

	// Call the GetLastExecutionLogForJobIds method
	jobExecutionLogs := jobExecutionsRepo.GetLastExecutionLogForJobIds([]uint64{1, 2})

	// Assert the number of retrieved job execution logs
	expectedCount := len(jobs)
	count := len(jobExecutionLogs)
	assert.Equal(t, expectedCount, count)

	// Assert the correctness of the retrieved job execution logs
	for _, log := range jobExecutionLogs {
		assert.Contains(t, []uint64{1, 2}, log.JobId)
	}
}
