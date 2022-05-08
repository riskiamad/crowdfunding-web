package transaction

import "crowdfunding-web/user"

type GetCampaignTransactionsInput struct {
	ID   int `json:"id" binding:"required"`
	User user.User
}
