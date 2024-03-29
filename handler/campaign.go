package handler

import (
	"fmt"
	"github.com/rizkyunm/senabung-api/campaign"
	"github.com/rizkyunm/senabung-api/helper"
	"github.com/rizkyunm/senabung-api/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	// catch parameter from handler
	// handler --> service
	// service choose which repository will call
	// repository : FindAll, FindByUserID
	// db

	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(uint(userID))
	if err != nil {
		response := helper.APIResponse("Can't get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	// handler : mapping ID from URL to struct input, call formatter to formatting data
	// service : need struct input to catch parameter ID
	// repository : get campaign by ID

	var input campaign.GetCampaignDetailInput

	if err := c.ShouldBindUri(&input); err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaign(input)
	if err != nil {
		errorMessage := gin.H{"errors": []string{err.Error()}}

		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if campaignDetail.ID == 0 {
		errorMessage := gin.H{"errors": []string{"campaign not found"}}

		response := helper.APIResponse("Campaign not found", http.StatusNotFound, "error", errorMessage)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	// catch parameter from user to struct input
	// catch current user from JWT
	// call service
	// call repository
	// save new campaign to db

	var input campaign.CreateCampaignInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been created", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	// user entry input
	// handler
	// mapping from input to input struct
	// input from user, and input from URI
	// service
	// repository update data campaign

	var inputID campaign.GetCampaignDetailInput

	if err := c.ShouldBindUri(&inputID); err != nil {
		response := helper.APIResponse("Failed update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	updateCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign has been updated", http.StatusOK, "success", campaign.FormatCampaign(updateCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UploadCampaignImage(c *gin.Context) {
	// handler
	// catch input and change to struct input
	// save image and assign to some folder
	// service (call repo point no. 2) --> (call repo point no.1)
	// repository :
	// 1. Create image data to table campaign_images
	// 2. Ubah isPrimary status from true to false

	var input campaign.CreateCampaignImageInput

	if err := c.ShouldBind(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := h.service.SaveCampaignImage(input, path); err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image has been uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
