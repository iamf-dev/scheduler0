package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"scheduler0/pkg/scheduler0time"
	"time"
)

type JobPriorityLevel uint64

type ExecutionTypes string

const (
	ExecutionTypeHTTP ExecutionTypes = "http"
)

// Job job model
type Job struct {
	ID                uint64    `json:"id,omitempty" fake:"{number:1,100}"`
	ProjectID         uint64    `json:"projectId,omitempty" fake:"{number:1,100}"`
	Spec              string    `json:"spec,omitempty"`
	CallbackUrl       string    `json:"callbackUrl,omitempty" fake:"{randomstring:[https://hello.com,https://world.com]}"`
	Data              string    `json:"data,omitempty"`
	ExecutionType     string    `json:"executionType,omitempty"`
	StartDate         time.Time `json:"startDate,omitempty"`
	EndDate           time.Time `json:"endDate,omitempty"`
	LastExecutionDate time.Time `json:"lastExecutionDate,omitempty"`
	Timezone          string    `json:"timezone,omitempty" fake:"{randomstring:[utc, America_NewYork]}"`
	TimezoneOffset    int64     `json:"timezoneOffset,omitempty"`
	ExecutionId       string    `json:"executionId,omitempty"`
	DateCreated       time.Time `json:"dateCreated,omitempty"`
}

// PaginatedJob paginated container of job transformer
type PaginatedJob struct {
	Total  uint64 `json:"total,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
	Data   []Job  `json:"jobs,omitempty"`
}

// ToJSON returns content of transformer as JSON
func (jobModel *Job) ToJSON() ([]byte, error) {
	if data, err := json.Marshal(jobModel); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

// FromJSON extracts content of JSON object into transformer
func (jobModel *Job) FromJSON(body []byte) error {
	if err := json.Unmarshal(body, &jobModel); err != nil {
		return err
	}
	return nil
}

func (jobModel *Job) GetNextExecutionTime() (*time.Time, error) {
	schedule, parseErr := cron.Parse(jobModel.Spec)
	if parseErr != nil {
		return nil, parseErr
	}
	var lastExecutionDateLocal time.Time
	if jobModel.LastExecutionDate.IsZero() {
		dateCreatedInLocal, err := jobModel.ConvertTimeToJobTimezone(jobModel.DateCreated)
		if err != nil {
			return nil, err
		}
		lastExecutionDateLocal = *dateCreatedInLocal
	} else {
		dateCreatedInLocal, err := jobModel.ConvertTimeToJobTimezone(jobModel.LastExecutionDate)
		if err != nil {
			return nil, err
		}
		lastExecutionDateLocal = *dateCreatedInLocal
	}
	schedulerTime := scheduler0time.GetSchedulerTime()
	now := schedulerTime.GetTime(time.Now())
	currentTime := lastExecutionDateLocal

	for now.Sub(currentTime).Round(time.Second*time.Duration(1)) >= time.Duration(0)*time.Second && currentTime.Before(now) {
		currentTime = schedule.Next(currentTime)
	}

	return &currentTime, nil
}

func (jobModel *Job) ConvertTimeToJobTimezone(timeToConvert time.Time) (*time.Time, error) {
	locale, err := time.LoadLocation(jobModel.Timezone)
	if err != nil {
		return nil, err
	}
	timeToConvertInTimezone := timeToConvert.In(locale)
	return &timeToConvertInTimezone, nil
}

func (jobModel *Job) GetNextExecutionId() (string, error) {
	nextExecutionTime, err := jobModel.GetNextExecutionTime()
	if err != nil {
		return "", nil
	}
	uniqueId := fmt.Sprintf(
		"%d-%d-%s-%s",
		jobModel.ProjectID,
		jobModel.ID,
		jobModel.LastExecutionDate.String(),
		nextExecutionTime.String(),
	)
	sha := sha256.New()
	return fmt.Sprintf("%x", sha.Sum([]byte(uniqueId))), err
}
