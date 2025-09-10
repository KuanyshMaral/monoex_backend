package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"monoex_backend/internal/models"
	"monoex_backend/internal/services"
)

type LegislationHandler struct {
	service *services.LegislationService
}

func NewLegislationHandler(service *services.LegislationService) *LegislationHandler {
	return &LegislationHandler{service: service}
}

// Create new legislation
func (h *LegislationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var legislation models.Legislation
	if err := json.NewDecoder(r.Body).Decode(&legislation); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &legislation); err != nil {
		// логируем ошибку на сервере
		fmt.Printf("❌ Failed to create legislation: %v\n", err)

		// возвращаем 500, а не 400
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(legislation)
}

// Get legislation by ID
func (h *LegislationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	legislation, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if legislation == nil {
		http.Error(w, "Legislation not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(legislation)
}

// Get all legislation with pagination
func (h *LegislationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit <= 0 {
		limit = 10
	}

	legislations, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total, err := h.service.GetTotalCount(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":   legislations,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update legislation
func (h *LegislationHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var legislation models.Legislation
	if err := json.NewDecoder(r.Body).Decode(&legislation); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	legislation.ID = id
	if err := h.service.Update(r.Context(), &legislation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(legislation)
}

// Delete legislation
func (h *LegislationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LegislationHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Ограничим размер файла
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("❌ FormFile error:", err)
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if len(handler.Filename) < 4 || handler.Filename[len(handler.Filename)-4:] != ".pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	// Можно добавить уникальный префикс, чтобы не перезаписывать файлы
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + handler.Filename
	savePath := "./uploads/legislation/" + fileName

	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Cannot save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Возвращаем путь для фронта
	response := map[string]string{
		"url": "/uploads/legislation/" + fileName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
