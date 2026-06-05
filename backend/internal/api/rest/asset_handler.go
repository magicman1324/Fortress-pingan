package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/service"
)

type AssetHandler struct {
	svc *service.AssetService
}

func NewAssetHandler(svc *service.AssetService) *AssetHandler {
	return &AssetHandler{svc: svc}
}

type createAssetReq struct {
	Name       string `json:"name" binding:"required"`
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port"`
	Username   string `json:"username" binding:"required"`
	AuthType   string `json:"auth_type"`
	Credential string `json:"credential" binding:"required"`
}

type updateAssetReq struct {
	Name       string `json:"name" binding:"required"`
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port"`
	Username   string `json:"username" binding:"required"`
	AuthType   string `json:"auth_type"`
	Credential string `json:"credential"`
}

func (h *AssetHandler) List(c *gin.Context) {
	search := c.Query("search")
	assets, err := h.svc.List(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if assets == nil {
		assets = []model.Asset{}
	}
	c.JSON(http.StatusOK, assets)
}

func (h *AssetHandler) Create(c *gin.Context) {
	var req createAssetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	if req.Port == 0 {
		req.Port = 22
	}
	if req.AuthType == "" {
		req.AuthType = "password"
	}
	asset, err := h.svc.Create(req.Name, req.Host, req.Port, req.Username, req.AuthType, req.Credential, "offline")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, asset)
}

func (h *AssetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateAssetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	if req.Port == 0 {
		req.Port = 22
	}
	if req.AuthType == "" {
		req.AuthType = "password"
	}
	asset, err := h.svc.Update(id, req.Name, req.Host, req.Port, req.Username, req.AuthType, req.Credential, "offline")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, asset)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *AssetHandler) TestConnect(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.TestConnect(id); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "connected"})
}
