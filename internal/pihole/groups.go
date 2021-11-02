package pihole

import (
	"context"
	"encoding/json"
	"net/url"
	"time"
)

type GroupResponseList struct {
	Data []GroupResponse
}

type GroupResponse struct {
	ID           int64  `json:"id"`
	Enabled      int    `json:"enabled"`
	Name         string `json:"name"`
	DateAdded    int64  `json:"date_added"`
	DateModified int64  `json:"date_modified"`
	Description  string `json:"description"`
}

type GroupList []*Group

type Group struct {
	ID           int64
	Enabled      bool
	Name         string
	DateAdded    time.Time
	DateModified time.Time
	Description  string
}

// ToGroup converts a GroupResponseList to a GroupList
func (grl GroupResponseList) ToGroupList() GroupList {
	list := make(GroupList, len(grl.Data))

	for i, g := range grl.Data {
		list[i] = g.ToGroup()
	}

	return list
}

// ToGroup converts a GroupResponse to a Group
func (gr GroupResponse) ToGroup() *Group {
	return &Group{
		ID:           gr.ID,
		Enabled:      gr.Enabled == 1,
		Name:         gr.Name,
		DateAdded:    time.Unix(gr.DateAdded, 0),
		DateModified: time.Unix(gr.DateModified, 0),
		Description:  gr.Description,
	}
}

// ListGroups returns the list of gravity DB groups
func (c Client) ListGroups(ctx context.Context) (GroupList, error) {
	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/groups.php", &url.Values{
		"action": []string{"get_groups"},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var groupRes GroupResponseList
	if err = json.NewDecoder(res.Body).Decode(&groupRes); err != nil {
		return nil, err
	}

	return groupRes.ToGroupList(), nil
}
