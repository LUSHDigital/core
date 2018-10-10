package auth

// Consumer represents an API consumer.
type Consumer struct {
	ID     int64   `json:"id"`
	Grants []Grant `json:"grants"`
}

// HasAnyGrant checks if a consumer possess any of a given set of grants
func (c *Consumer) HasAnyGrant(grants ...Grant) bool {
	for _, grant := range grants {
		for _, g := range c.Grants {
			if grant == g {
				return true
			}
		}
	}

	return false
}
