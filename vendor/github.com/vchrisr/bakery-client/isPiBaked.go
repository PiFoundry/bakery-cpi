package bakeryclient

func (c *Client) IsPiBaked(piId string) (bool, error) {
	pi, err := c.GetPi(piId)
	if err != nil {
		return false, err
	}

	if pi.Status != INUSE {
		return false, nil
	}

	return true, nil
}
