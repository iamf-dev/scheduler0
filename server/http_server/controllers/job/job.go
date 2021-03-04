package job

import (
	"github.com/gorilla/mux"
	"net/http"
	"scheduler0/server/http_server/controllers"
	"scheduler0/server/service"
	"scheduler0/server/transformers"
	"scheduler0/utils"
	"strconv"
)

type Controller controllers.Controller

func (jobController *Controller) List(w http.ResponseWriter, r *http.Request) {
	projectUUID, err := utils.ValidateQueryString("projectUUID", r)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	limitParam, err := utils.ValidateQueryString("limit", r)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	offsetParam, err := utils.ValidateQueryString("offset", r)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	jobService := service.JobService{
		Pool: jobController.Pool,
		Ctx:  r.Context(),
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	jobs, getJobsByProjectIDError := jobService.GetJobsByProjectUUID(projectUUID, offset, limit, "date_created")
	if getJobsByProjectIDError != nil {
		utils.SendJSON(w, getJobsByProjectIDError.Message, false, getJobsByProjectIDError.Type, nil)
		return
	}

	utils.SendJSON(w, jobs, true, http.StatusOK, nil)
}

func (jobController *Controller) CreateOne(w http.ResponseWriter, r *http.Request) {
	body := utils.ExtractBody(w, r)
	jobBody := transformers.Job{}
	err := jobBody.FromJSON(body)

	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	jobService := service.JobService{
		Pool: jobController.Pool,
		Ctx:  r.Context(),
	}

	job, createJobError := jobService.CreateJob(jobBody)
	if createJobError != nil {
		utils.SendJSON(w, createJobError.Message, false, createJobError.Type, nil)
		return
	}

	utils.SendJSON(w, job, true, http.StatusCreated, nil)
}

func (jobController *Controller) GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	jobService := service.JobService{
		Pool: jobController.Pool,
		Ctx:  r.Context(),
	}

	job := transformers.Job{
		UUID: params["uuid"],
	}

	jobT, getOneJobError := jobService.GetJob(job)
	if getOneJobError != nil {
		utils.SendJSON(w, getOneJobError.Message, false, getOneJobError.Type, nil)
		return
	}

	utils.SendJSON(w, jobT, true, http.StatusOK, nil)
}

func (jobController *Controller) UpdateOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	body := utils.ExtractBody(w, r)
	jobBody := transformers.Job{}
	err := jobBody.FromJSON(body)

	jobBody.UUID = params["uuid"]

	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	jobService := service.JobService{
		Pool: jobController.Pool,
		Ctx:  r.Context(),
	}

	jobT, updateOneJobError := jobService.UpdateJob(jobBody)
	if updateOneJobError != nil {
		utils.SendJSON(w, updateOneJobError.Message, false, updateOneJobError.Type, nil)
		return
	}

	utils.SendJSON(w, jobT, true, http.StatusOK, nil)
}

func (jobController *Controller) DeleteOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	jobService := service.JobService{
		Pool: jobController.Pool,
		Ctx:  r.Context(),
	}

	job := transformers.Job{
		UUID: params["uuid"],
	}

	deleteOneJobError := jobService.DeleteJob(job)
	if deleteOneJobError != nil {
		utils.SendJSON(w, deleteOneJobError.Message, false, deleteOneJobError.Type, nil)
		return
	}

	utils.SendJSON(w, nil, true, http.StatusNoContent, nil)
}