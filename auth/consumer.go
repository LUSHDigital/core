package auth

// Consumer represents an API user
type Consumer struct {
	ID        int64    `json:"id"`
	UUID      string   `json:"uuid"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Language  string   `json:"language"`
	Grants    []string `json:"grants"`
	Roles     []string `json:"roles"`
	Needs     []string `json:"needs"`
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

// MissesAnyGrant checks if a consumer is missing any of a given set of grants
func (c Consumer) MissesAnyGrant(grants ...string) bool {
	for _, grant := range grants {
		for _, g := range c.Grants {
			if grant == g {
				return false
			}
		}
	}

	return true
}

// HasAnyRole checks if a consumer possess any of a given set of roles
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

// MissesAnyRole checks if a consumer is missing any of a given set of roles
func (c *Consumer) MissesAnyRole(roles ...string) bool {
	for _, role := range roles {
		for _, r := range c.Roles {
			if role == r {
				return false
			}
		}
	}

	return true
}

// HasAnyNeed checks if a consumer has any of the given needs
func (c *Consumer) HasAnyNeed(needs ...string) bool {
	for _, role := range needs {
		for _, r := range c.Needs {
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

// HasUUID checks if a consumer has the same uuid as a user
func (c *Consumer) HasUUID(id string) bool {
	return c.UUID == id
}
