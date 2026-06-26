package adapters

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sosvzla/sosvzla.lat/app"
	"github.com/sosvzla/sosvzla.lat/domain"
)

// HTTPHandler handles HTTP requests for the API.
type HTTPHandler struct {
	searchService *app.SearchService
}

// NewHTTPHandler creates a new HTTPHandler.
func NewHTTPHandler(searchService *app.SearchService) *HTTPHandler {
	return &HTTPHandler{searchService: searchService}
}

// RegisterRoutes registers the API routes to the given multiplexer.
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	// Persons
	mux.HandleFunc("POST /persons", h.RegisterPerson)
	mux.HandleFunc("GET /persons/{id}", h.GetPersonByID)
	mux.HandleFunc("PUT /persons/{id}", h.UpdatePerson)
	mux.HandleFunc("DELETE /persons/{id}", h.DeletePerson)
	mux.HandleFunc("GET /persons", h.ListPersons)

	// Reports
	mux.HandleFunc("POST /reports", h.FileReport)
	mux.HandleFunc("GET /persons/{id}/reports", h.ListReportsByPersonID)

	// Hospitals
	mux.HandleFunc("POST /hospitals", h.RegisterHospital)
	mux.HandleFunc("GET /hospitals", h.ListHospitals)
	mux.HandleFunc("POST /hospitals/{id}/admit", h.AdmitPersonToHospital)
	mux.HandleFunc("GET /hospitals/{id}/persons", h.ListPersonsInHospital)
}

// JSON helper for responding with data.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Error helper.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// RegisterPerson registers a new person.
// @Summary Register a new person
// @Description Register a missing, found, or in_hospital person in the centralized database.
// @Tags persons
// @Accept json
// @Produce json
// @Param person body domain.Person true "Person Object"
// @Success 201 {object} domain.Person
// @Failure 400 {object} map[string]string "Invalid payload"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons [post]
func (h *HTTPHandler) RegisterPerson(w http.ResponseWriter, r *http.Request) {
	var person domain.Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := h.searchService.RegisterPerson(r.Context(), &person)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, person)
}

// GetPersonByID retrieves a person by their ID.
// @Summary Get a person by ID
// @Description Get details of a registered person using their database ID.
// @Tags persons
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {object} domain.Person
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons/{id} [get]
func (h *HTTPHandler) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid person ID")
		return
	}

	person, err := h.searchService.GetPersonByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if person == nil {
		respondWithError(w, http.StatusNotFound, "Person not found")
		return
	}

	respondWithJSON(w, http.StatusOK, person)
}

// GetPersonByNationalID retrieves a person by their national ID.
// @Summary Get a person by National ID
// @Description Get details of a registered person using their national identification number (e.g., Cédula).
// @Tags persons
// @Produce json
// @Param national_id path string true "National ID"
// @Success 200 {object} domain.Person
// @Failure 404 {object} map[string]string "Person not found"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons/national/{national_id} [get]
func (h *HTTPHandler) GetPersonByNationalID(w http.ResponseWriter, r *http.Request) {
	nationalID := r.PathValue("national_id")
	if nationalID == "" {
		respondWithError(w, http.StatusBadRequest, "National ID is required")
		return
	}

	person, err := h.searchService.GetPersonByNationalID(r.Context(), nationalID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if person == nil {
		respondWithError(w, http.StatusNotFound, "Person not found")
		return
	}

	respondWithJSON(w, http.StatusOK, person)
}

// UpdatePerson updates an existing person's details.
// @Summary Update a person
// @Description Update the personal details or status of a person.
// @Tags persons
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body domain.Person true "Updated Person Object"
// @Success 200 {object} domain.Person
// @Failure 400 {object} map[string]string "Invalid ID or payload"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons/{id} [put]
func (h *HTTPHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid person ID")
		return
	}

	var person domain.Person
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	person.ID = id
	err = h.searchService.UpdatePerson(r.Context(), &person)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, person)
}

// DeletePerson deletes a person by ID.
// @Summary Delete a person
// @Description Remove a person's entry from the database.
// @Tags persons
// @Produce json
// @Param id path int true "Person ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons/{id} [delete]
func (h *HTTPHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid person ID")
		return
	}

	err = h.searchService.DeletePerson(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPersons lists persons with filtering and pagination.
// @Summary List or search persons
// @Description Retrieve a list of persons, optionally filtered by status (missing, found, in_hospital) or searched by national_id.
// @Tags persons
// @Produce json
// @Param status query string false "Filter by status"
// @Param national_id query string false "Search by National ID (Cédula)"
// @Param limit query int false "Pagination limit (default: 10)"
// @Param offset query int false "Pagination offset (default: 0)"
// @Success 200 {array} domain.Person
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons [get]
func (h *HTTPHandler) ListPersons(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	nationalID := q.Get("national_id")
	if nationalID != "" {
		person, err := h.searchService.GetPersonByNationalID(r.Context(), nationalID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if person == nil {
			respondWithJSON(w, http.StatusOK, []*domain.Person{})
			return
		}
		respondWithJSON(w, http.StatusOK, []*domain.Person{person})
		return
	}

	status := q.Get("status")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	persons, err := h.searchService.ListPersons(r.Context(), status, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if persons == nil {
		persons = []*domain.Person{}
	}

	respondWithJSON(w, http.StatusOK, persons)
}

// FileReport files a report about a person.
// @Summary File a search/found report
// @Description Create a report regarding a missing or found person.
// @Tags reports
// @Accept json
// @Produce json
// @Param report body domain.Report true "Report Object"
// @Success 201 {object} domain.Report
// @Failure 400 {object} map[string]string "Invalid payload"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /reports [post]
func (h *HTTPHandler) FileReport(w http.ResponseWriter, r *http.Request) {
	var report domain.Report
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&report); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := h.searchService.FileReport(r.Context(), &report)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, report)
}

// ListReportsByPersonID lists reports for a specific person.
// @Summary List reports for a person
// @Description Retrieve all filed reports associated with a specific person ID.
// @Tags reports
// @Produce json
// @Param id path int true "Person ID"
// @Success 200 {array} domain.Report
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /persons/{id}/reports [get]
func (h *HTTPHandler) ListReportsByPersonID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid person ID")
		return
	}

	reports, err := h.searchService.ListReportsByPersonID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, reports)
}

// RegisterHospital registers a new hospital.
// @Summary Register a hospital
// @Description Register a health center to track admitted persons.
// @Tags hospitals
// @Accept json
// @Produce json
// @Param hospital body domain.Hospital true "Hospital Object"
// @Success 201 {object} domain.Hospital
// @Failure 400 {object} map[string]string "Invalid payload"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /hospitals [post]
func (h *HTTPHandler) RegisterHospital(w http.ResponseWriter, r *http.Request) {
	var hospital domain.Hospital
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&hospital); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := h.searchService.RegisterHospital(r.Context(), &hospital)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, hospital)
}

// ListHospitals lists hospitals.
// @Summary List hospitals
// @Description Retrieve a list of registered hospitals with pagination.
// @Tags hospitals
// @Produce json
// @Param limit query int false "Pagination limit (default: 10)"
// @Param offset query int false "Pagination offset (default: 0)"
// @Success 200 {array} domain.Hospital
// @Failure 500 {object} map[string]string "Internal error"
// @Router /hospitals [get]
func (h *HTTPHandler) ListHospitals(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	hospitals, err := h.searchService.ListHospitals(r.Context(), limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, hospitals)
}

// AdmitPersonToHospital admits a person to a hospital.
// @Summary Admit a person to a hospital
// @Description Associate a person with a hospital admission.
// @Tags hospitals
// @Accept json
// @Produce json
// @Param id path int true "Hospital ID"
// @Param admission body domain.HospitalPerson true "Hospital Person association details"
// @Success 200 {object} domain.HospitalPerson
// @Failure 400 {object} map[string]string "Invalid ID or payload"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /hospitals/{id}/admit [post]
func (h *HTTPHandler) AdmitPersonToHospital(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	hospitalID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid hospital ID")
		return
	}

	var hp domain.HospitalPerson
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&hp); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	hp.HospitalID = hospitalID
	err = h.searchService.AdmitPersonToHospital(r.Context(), &hp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, hp)
}

// ListPersonsInHospital lists persons admitted to a specific hospital.
// @Summary List persons admitted in a hospital
// @Description Retrieve a list of persons who are currently admitted or listed in a specific hospital.
// @Tags hospitals
// @Produce json
// @Param id path int true "Hospital ID"
// @Success 200 {array} domain.Person
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /hospitals/{id}/persons [get]
func (h *HTTPHandler) ListPersonsInHospital(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	hospitalID, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid hospital ID")
		return
	}

	persons, err := h.searchService.ListPersonsInHospital(r.Context(), hospitalID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// This step is specifically for fixing the output where [] is expected instead of null if empty.
	if persons == nil {
		persons = []*domain.Person{}
	}

	// Support CORS or general responses cleanly
	respondWithJSON(w, http.StatusOK, persons)
}

// Simple CORS middleware for demo/dev purposes
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
