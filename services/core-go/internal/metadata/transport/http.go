package transport

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/user/seta-dam-backend/services/core-go/internal/common"
	folderDomain "github.com/user/seta-dam-backend/services/core-go/internal/folder/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/metadata/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/metadata/usecase"
	permDomain "github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
	"net/http"
)

type MetadataHandler struct {
	usecase usecase.MetadataUsecase
}

func NewMetadataHandler(u usecase.MetadataUsecase) *MetadataHandler {
	return &MetadataHandler{usecase: u}
}

func (h *MetadataHandler) RegisterRoutes(r chi.Router) {
	r.Post("/metadata", h.CreateMetadata)
	r.Get("/metadata", h.ListMetadata)
	r.Get("/metadata/{id}", h.GetMetadataByID)
	r.Put("/metadata/{id}", h.UpdateMetadata)
	r.Delete("/metadata/{id}", h.DeleteMetadata)
	r.Get("/folders/{folderId}/metadata", h.ListByFolder)
}

func (h *MetadataHandler) CreateMetadata(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())

	var req domain.Metadata
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	meta, err := h.usecase.Create(r.Context(), uCtx, &req)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, meta)
}

func (h *MetadataHandler) GetMetadataByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	meta, err := h.usecase.GetByID(r.Context(), uCtx, id)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}
	h.respondWithJSON(w, http.StatusOK, meta)
}

func (h *MetadataHandler) UpdateMetadata(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	var req domain.Metadata
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id

	meta, err := h.usecase.Update(r.Context(), uCtx, &req)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, meta)
}

func (h *MetadataHandler) DeleteMetadata(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uCtx := common.GetUserContext(r.Context())

	err := h.usecase.Delete(r.Context(), uCtx, id)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "metadata deleted successfully"})
}

func (h *MetadataHandler) ListMetadata(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())

	list, err := h.usecase.List(r.Context(), uCtx)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, list)
}

func (h *MetadataHandler) ListByFolder(w http.ResponseWriter, r *http.Request) {
	folderID := chi.URLParam(r, "folderId")
	uCtx := common.GetUserContext(r.Context())

	list, err := h.usecase.ListByFolder(r.Context(), uCtx, folderID)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, list)
}

func (h *MetadataHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *MetadataHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *MetadataHandler) respondWithErrorForCode(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrMetadataNotFound) || errors.Is(err, folderDomain.ErrFolderNotFound) {
		h.respondWithError(w, http.StatusNotFound, err.Error())
	} else if errors.Is(err, domain.ErrTitleRequired) || errors.Is(err, domain.ErrFolderRequired) || errors.Is(err, domain.ErrInvalidMetadataJSON) {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
	} else if errors.Is(err, permDomain.ErrForbidden) {
		h.respondWithError(w, http.StatusForbidden, err.Error())
	} else if errors.Is(err, permDomain.ErrUnauthorized) {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
	} else {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
