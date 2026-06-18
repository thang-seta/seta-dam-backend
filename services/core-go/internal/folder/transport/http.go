package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/user/seta-dam-backend/services/core-go/internal/common"
	"github.com/user/seta-dam-backend/services/core-go/internal/folder/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/folder/usecase"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
)

type FolderHandler struct {
	usecase usecase.FolderUsecase
}

func NewFolderHandler(u usecase.FolderUsecase) *FolderHandler {
	return &FolderHandler{usecase: u}
}

func (h *FolderHandler) RegisterRoutes(r chi.Router) {
	r.Post("/folders", h.CreateFolder)
	r.Get("/folders/tree", h.GetTree)
	r.Get("/folders/{id}", h.GetFolderByID)
	r.Put("/folders/{id}", h.UpdateFolder)
	r.Put("/folders/{id}/move", h.MoveFolder)
	r.Delete("/folders/{id}", h.DeleteFolder)
}

func (h *FolderHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())
	
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		ParentID    *string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" {
		h.respondWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	folder, err := h.usecase.Create(r.Context(), uCtx, req.Name, req.Description, req.ParentID)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, folder)
}

func (h *FolderHandler) GetTree(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())
	tree, err := h.usecase.ListTree(r.Context(), uCtx)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, tree)
}

func (h *FolderHandler) GetFolderByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	folder, err := h.usecase.GetByID(r.Context(), uCtx, id)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, folder)
}

func (h *FolderHandler) UpdateFolder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" {
		h.respondWithError(w, http.StatusBadRequest, "name is required")
		return
	}

	err := h.usecase.Update(r.Context(), uCtx, id, req.Name, req.Description)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "folder updated successfully"})
}

func (h *FolderHandler) MoveFolder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	var req struct {
		ParentID *string `json:"parent_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.usecase.MoveFolder(r.Context(), uCtx, id, req.ParentID)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "folder moved successfully"})
}

func (h *FolderHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	err := h.usecase.Delete(r.Context(), uCtx, id)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "folder deleted successfully"})
}

func (h *FolderHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *FolderHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *FolderHandler) respondWithErrorForCode(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrFolderNotFound) {
		h.respondWithError(w, http.StatusNotFound, err.Error())
	} else if errors.Is(err, domain.ErrCycleDetected) || errors.Is(err, domain.ErrNotEmpty) {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
	} else if errors.Is(err, permDomain.ErrForbidden) {
		h.respondWithError(w, http.StatusForbidden, err.Error())
	} else if errors.Is(err, permDomain.ErrUnauthorized) {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
	} else {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
