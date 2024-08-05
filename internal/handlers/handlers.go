package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/risetyll/finuslugi/internal/entities/requests"
	"github.com/risetyll/finuslugi/internal/usecase"
)

type RouteInfo struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type Routes struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Routes {
	return &Routes{
		logger: logger,
	}
}

func (routes *Routes) GetRoutesInfoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := []RouteInfo{
			{
				Path:        "/materials",
				Method:      "POST",
				Description: "Создать новый материал",
			},
			{
				Path:        "/materials/{id:[0-9]+}",
				Method:      "GET",
				Description: "Получить материал по ID",
			},
			{
				Path:        "/materials/{id:[0-9]+}",
				Method:      "PUT",
				Description: "Обновить материал по ID",
			},
			{
				Path:        "/materials",
				Method:      "GET",
				Description: "Получить список материалов с фильтрацией",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(routes); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (routes *Routes) CreateMaterialHandler(uc *usecase.DatabaseUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.CreateMaterialRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		routes.logger.Debug("create request", req)

		if err := uc.CreateMaterial(&req); err != nil {
			http.Error(w, "Failed to create material", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (routes *Routes) GetMaterialByIdHandler(uc *usecase.DatabaseUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuidStr := vars["id"]
		uuid, err := strconv.Atoi(uuidStr)
		if err != nil {
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}
		req := requests.GetMaterialByIdRequest{UUID: uuid}

		routes.logger.Debug("get request", req)

		material, err := uc.GetMaterialById(&req)
		if err != nil {
			http.Error(w, "Failed to get material", http.StatusInternalServerError)
			return
		}
		if material == nil {
			http.Error(w, "Material not found", http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(material); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (routes *Routes) UpdateMaterialHandler(uc *usecase.DatabaseUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.UpdateMaterialRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		routes.logger.Debug("update request", req)

		if err := uc.UpdateMaterial(&req); err != nil {
			http.Error(w, "Failed to update material", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (routes *Routes) GetMaterialsHandler(uc *usecase.DatabaseUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &requests.GetMaterialsRequest{
			Type:        r.URL.Query().Get("type"),
			Page:        parseQueryParam(r.URL.Query().Get("page"), 1),
			PageSize:    parseQueryParam(r.URL.Query().Get("page_size"), 10),
			CreatedFrom: parseTimeQueryParam(r.URL.Query().Get("created_from")),
			CreatedTo:   parseTimeQueryParam(r.URL.Query().Get("created_to")),
		}

		routes.logger.Debug("get request", req)

		materials, err := uc.GetMaterials(req)
		if err != nil {
			http.Error(w, "Failed to get materials", http.StatusInternalServerError)
			return
		}
		response := &requests.GetMaterialsResponse{
			Materials: materials,
			Total:     len(materials),
			Page:      req.Page,
			PageSize:  req.PageSize,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func parseQueryParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func parseTimeQueryParam(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}
