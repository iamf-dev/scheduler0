package async_task_manager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"scheduler0/models"
	"scheduler0/repository"
	"scheduler0/utils"
	"sync"
)

type AsyncTaskManager struct {
	task                 sync.Map // map[uint64]models.AsyncTask
	taskIdRequestIdMap   sync.Map // map[string]uint64
	subscribers          sync.Map // map[uint64]map[uint64]func(task models.AsyncTask)
	subscriberIds        sync.Map // map[uint64][]uint64
	asyncTaskManagerRepo repository.AsyncTasksRepo
	context              context.Context
	logger               *log.Logger
	notificationsCh      chan models.AsyncTask
}

func NewAsyncTaskManager(context context.Context, logger *log.Logger, asyncTaskManagerRepo repository.AsyncTasksRepo) *AsyncTaskManager {
	return &AsyncTaskManager{
		context:              context,
		logger:               logger,
		asyncTaskManagerRepo: asyncTaskManagerRepo,
		notificationsCh:      make(chan models.AsyncTask, 1),
	}
}

func (m *AsyncTaskManager) AddTasks(input, requestId string, service string) ([]uint64, *utils.GenericError) {
	tasks := []models.AsyncTask{
		models.AsyncTask{
			Input:     input,
			RequestId: requestId,
			Service:   service,
		},
	}
	ids, err := m.asyncTaskManagerRepo.BatchInsert(tasks)
	if err != nil {
		return nil, err
	}
	tasks[0].Id = ids[0]
	m.task.Store(ids[0], tasks[0])
	m.taskIdRequestIdMap.Store(tasks[0].RequestId, ids[0])
	return ids, nil
}

func (m *AsyncTaskManager) UpdateTasksById(taskId uint64, state models.AsyncTaskState, output string) *utils.GenericError {
	t, ok := m.task.Load(taskId)
	if !ok {
		m.logger.Println("could not find task with id", taskId)
		return nil
	}
	myT := t.(models.AsyncTask)
	myT.State = state
	m.task.Store(taskId, myT)
	err := m.asyncTaskManagerRepo.UpdateTaskState(myT, state, output)
	if err != nil {
		m.logger.Println("could not update task with id", taskId)
		return err
	}
	go func() { m.notificationsCh <- myT }()
	return nil
}

func (m *AsyncTaskManager) UpdateTasksByRequestId(requestId string, state models.AsyncTaskState, output string) *utils.GenericError {
	tId, ok := m.taskIdRequestIdMap.Load(requestId)
	if !ok {
		m.logger.Println("could not find task id for request id", requestId)
		return nil
	}
	t, ok := m.task.Load(tId)
	if !ok {
		m.logger.Println("could not find task with request id task id", requestId)
		return nil
	}
	myT := t.(models.AsyncTask)
	myT.State = state
	m.task.Store(myT.Id, myT)
	err := m.asyncTaskManagerRepo.UpdateTaskState(myT, state, output)
	if err != nil {
		m.logger.Println("could not update task request id", requestId)
		return err
	}
	go func() { m.notificationsCh <- myT }()
	return nil
}

func (m *AsyncTaskManager) AddSubscriber(taskId uint64, subscriber func(task models.AsyncTask)) (uint64, *utils.GenericError) {
	t, ok := m.task.Load(taskId)
	if !ok {
		m.logger.Println("could not find task with id", taskId)
		return 0, utils.HTTPGenericError(http.StatusNotFound, fmt.Sprintf("could not find task with id %d", taskId))
	}
	myt := t.(models.AsyncTask)
	subIds, ok := m.subscriberIds.Load(taskId)
	var maxId int64 = 0
	if ok {
		maxId = subIds.(int64)
	}
	subId := uint64(maxId + 1)

	sb, ok := m.subscribers.Load(myt.Id)
	var subscribers = map[uint64]func(task models.AsyncTask){}
	if ok {
		storedsubs := sb.(map[uint64]func(task models.AsyncTask))
		subscribers = storedsubs
	}
	subscribers[subId] = subscriber
	m.subscribers.Store(taskId, subscribers)
	m.subscriberIds.Store(taskId, subId)

	return subId, nil
}

func (m *AsyncTaskManager) GetTask(taskId uint64) (chan models.AsyncTask, uint64, *utils.GenericError) {
	task, err := m.asyncTaskManagerRepo.GetTask(taskId)
	if err != nil {
		return nil, 0, err
	}
	if task.State != models.AsyncTaskInProgress && task.State != models.AsyncTaskNotStated {
		return nil, 0, nil
	}

	var taskCh = make(chan models.AsyncTask, 1)

	subs, addErr := m.AddSubscriber(taskId, func(task models.AsyncTask) {
		taskCh <- task
	})
	if addErr != nil {
		return nil, 0, addErr
	}

	return taskCh, subs, nil
}

func (m *AsyncTaskManager) GetTaskWithRequestId(requestId string) (chan models.AsyncTask, uint64, *utils.GenericError) {
	taskId, ok := m.taskIdRequestIdMap.Load(requestId)
	if !ok {
		return nil, 0, nil
	}
	return m.GetTask(taskId.(uint64))
}

func (m *AsyncTaskManager) GetTaskIdWithRequestId(requestId string) (uint64, *utils.GenericError) {
	taskId, ok := m.taskIdRequestIdMap.Load(requestId)
	if ok {
		return taskId.(uint64), nil

	}
	return 0, nil
}

func (m *AsyncTaskManager) DeleteSubscriber(taskId, subscriberId uint64) *utils.GenericError {
	t, ok := m.task.Load(taskId)
	if !ok {
		return utils.HTTPGenericError(http.StatusNotFound, fmt.Sprintf("could not find task with id %d", taskId))
	}
	myt := t.(models.AsyncTask)
	sb, ok := m.subscribers.Load(myt.Id)
	if !ok {
		return utils.HTTPGenericError(http.StatusNotFound, fmt.Sprintf("could not find subscribers for task with id %d", taskId))
	}
	subscribers := sb.(map[uint64]func(task models.AsyncTask))
	delete(subscribers, subscriberId)
	m.subscribers.Store(taskId, subscribers)
	return nil
}

func (m *AsyncTaskManager) ListenForNotifications() {
	go func() {
		for {
			select {
			case taskNotification := <-m.notificationsCh:
				t, ok := m.task.Load(taskNotification.Id)
				if !ok {
					m.logger.Println("could not find task with id", taskNotification.Id)
					return
				}
				myt := t.(models.AsyncTask)
				sb, ok := m.subscribers.Load(myt.Id)

				var subscribers = map[uint64]func(task models.AsyncTask){}
				if ok {
					subscribers = sb.(map[uint64]func(task models.AsyncTask))
				}

				for _, subscriber := range subscribers {
					subscriber(taskNotification)
				}
				if taskNotification.State == models.AsyncTaskFail || taskNotification.State == models.AsyncTaskSuccess {
					m.task.Delete(myt.Id)
				}
				m.subscribers.Delete(myt.Id)
			case <-m.context.Done():
				return
			}
		}
	}()
}
