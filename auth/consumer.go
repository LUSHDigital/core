package auth

// Consumer represents an API user
type Consumer struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Language  string   `json:"language"`
	Grants    []string `json:"grants"`
	Roles     []string `json:"roles"`
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

// HasAnyGrant checks if a consumer possess any of a given set of grants
func (c *Consumer) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		for _, r := range c.Roles {
			if role == r {
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
