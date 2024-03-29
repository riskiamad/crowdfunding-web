package campaign

import "github.com/rizkyunm/senabung-api/user"

type GetCampaignDetailInput struct {
	ID uint `uri:"id" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string  `json:"name" binding:"required"`
	ShortDescription string  `json:"short_description" binding:"required"`
	Description      string  `json:"description" binding:"required"`
	GoalAmount       float64 `json:"goal_amount" binding:"required"`
	Perks            string  `json:"perks" binding:"required"`
	User             user.User
}

type CreateCampaignImageInput struct {
	CampaignID uint `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary"`
	User       user.User
}
