package auth

// Consumer represents an API consumer.
type Consumer struct {
	ID     int64    `json:"id"`
	Grants []string `json:"grants"`
}

// HasAnyGrant checks if a consumer possess any of a given set of grants
func (c *Consumer) HasAnyGrant(grants ...string) bool {
	for _, grant := range grants {
		for _, g := range c.Grants {
			if grant == g {
				return true
			}
		}
	}

	return false
}

// IsUser checks if a consumer has the same ID as a user
func (c *Consumer) IsUser(userID int64) bool {
	return c.ID == userID
}
