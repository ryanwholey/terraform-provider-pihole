package pihole

import (
	"context"
	"encoding/json"
	"fmt"
)

type EnableAdBlockResponse struct {
	Status string
}

// ToEnableAdBlock turns am EnableAdBlockResponse into an EnableAdBlock object
func (eb EnableAdBlockResponse) ToEnableAdBlock() *EnableAdBlock {
	return &EnableAdBlock{
		Enabled: eb.Status == "enabled",
	}
}

type EnableAdBlock struct {
	Enabled bool
}

// GetAdBlockerStatus returns whether pihole ad blocking is enabled or not
func (c Client) GetAdBlockerStatus(ctx context.Context) (*EnableAdBlock, error) {
	req, err := c.Request(ctx, "GET", "/admin/api.php?status", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var blocked EnableAdBlockResponse
	if err := json.NewDecoder(res.Body).Decode(&blocked); err != nil {
		return nil, err
	}

	return blocked.ToEnableAdBlock(), nil
}

// SetAdBlockEnabled sets whether pihole ad blocking is enabled or not
func (c Client) SetAdBlockEnabled(ctx context.Context, enable bool) (*EnableAdBlock, error) {
	enabledParam := "enable"
	if !enable {
		enabledParam = "disable"
	}

	req, err := c.RequestWithAuth(ctx, "GET", fmt.Sprintf("/admin/api.php?%s", enabledParam), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var blocked EnableAdBlockResponse
	if err = json.NewDecoder(res.Body).Decode(&blocked); err != nil {
		return nil, err
	}

	if blocked.Status != fmt.Sprintf("%sd", enabledParam) {
		return nil, fmt.Errorf("ad blocking could not be turned to %q", enabledParam)
	}

	return blocked.ToEnableAdBlock(), nil
}
