package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/user/seta-dam-backend/services/core-go/internal/common"
	"github.com/user/seta-dam-backend/services/core-go/internal/permission/domain"
	"github.com/user/seta-dam-backend/services/core-go/internal/permission/usecase"
)

type PermissionHandler struct {
	evaluator usecase.PermissionEvaluator
}

func NewPermissionHandler(eval usecase.PermissionEvaluator) *PermissionHandler {
	return &PermissionHandler{evaluator: eval}
}

func (h *PermissionHandler) RegisterRoutes(r chi.Router) {
	r.Post("/permissions", h.AddPermission)
	r.Delete("/permissions", h.DeletePermission)
	r.Get("/permissions/effective", h.GetEffectivePermissions)
}

func (h *PermissionHandler) AddPermission(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())

	var req struct {
		UserID     string `json:"user_id"`
		ObjectType string `json:"object_type"`
		ObjectID   string `json:"object_id"`
		Action     string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.UserID == "" || req.ObjectType == "" || req.ObjectID == "" || req.Action == "" {
		h.respondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	perm := &domain.ObjectPermission{
		UserID:     req.UserID,
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectID,
		Action:     req.Action,
	}

	err := h.evaluator.AddPermission(r.Context(), uCtx, perm)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "permission granted successfully"})
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())

	var req struct {
		UserID     string `json:"user_id"`
		ObjectType string `json:"object_type"`
		ObjectID   string `json:"object_id"`
		Action     string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.UserID == "" || req.ObjectType == "" || req.ObjectID == "" || req.Action == "" {
		h.respondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	err := h.evaluator.DeletePermission(r.Context(), uCtx, req.UserID, req.ObjectType, req.ObjectID, req.Action)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "permission revoked successfully"})
}

func (h *PermissionHandler) GetEffectivePermissions(w http.ResponseWriter, r *http.Request) {
	uCtx := common.GetUserContext(r.Context())

	userID := r.URL.Query().Get("user_id")
	objectType := r.URL.Query().Get("object_type")
	objectID := r.URL.Query().Get("object_id")

	if userID == "" || objectType == "" || objectID == "" {
		h.respondWithError(w, http.StatusBadRequest, "missing required query parameters: user_id, object_type, object_id")
		return
	}

	actions, err := h.evaluator.GetEffectivePermissions(r.Context(), uCtx, userID, objectType, objectID)
	if err != nil {
		h.respondWithErrorForCode(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":     userID,
		"object_type": objectType,
		"object_id":   objectID,
		"actions":     actions,
	})
}

func (h *PermissionHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (h *PermissionHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *PermissionHandler) respondWithErrorForCode(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrForbidden) {
		h.respondWithError(w, http.StatusForbidden, err.Error())
	} else if errors.Is(err, domain.ErrUnauthorized) {
		h.respondWithError(w, http.StatusUnauthorized, err.Error())
	} else {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
