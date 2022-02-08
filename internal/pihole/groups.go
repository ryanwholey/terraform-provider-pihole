package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
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

type GroupUpdateRequest struct {
	Name        string
	Enabled     *bool
	Description string
}

type GroupCreateRequest struct {
	Name        string
	Description string
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

// GetGroup returns a Pi-hole group by name
func (c Client) GetGroup(ctx context.Context, name string) (*Group, error) {
	groups, err := c.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.Name == name {
			return g, nil
		}
	}
	return nil, NewNotFoundError(fmt.Sprintf("Group with name %q not found", name))
}

// GetGroupByID returns a Pi-hole group by ID
func (c Client) GetGroupByID(ctx context.Context, id int64) (*Group, error) {
	groups, err := c.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.ID == id {
			return g, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("Group with ID %q not found", id))
}

// validName indicates whether the name given to the group is valid
func validGroupName(name string) bool {
	validName := regexp.MustCompile(`^\S*$`)

	return validName.MatchString(name)
}

type GroupBasicResponse struct {
	Success bool
	Message string
}

// CreateGroup creates a group with the passed attributes
func (c Client) CreateGroup(ctx context.Context, gr *GroupCreateRequest) (*Group, error) {
	name := strings.TrimSpace(gr.Name)

	if !validGroupName(name) {
		return nil, fmt.Errorf("group names must not contain spaces")
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/groups.php", &url.Values{
		"action": []string{"add_group"},
		"name":   []string{name},
		"desc":   []string{gr.Description},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var created GroupBasicResponse
	if err = json.NewDecoder(res.Body).Decode(&created); err != nil {
		return nil, err
	}

	if !created.Success {
		return nil, fmt.Errorf(created.Message)
	}

	return c.GetGroup(ctx, name)
}

// UpdateGroup updates a group resource with the passed attribute
func (c Client) UpdateGroup(ctx context.Context, gr *GroupUpdateRequest) (*Group, error) {
	original, err := c.GetGroup(ctx, gr.Name)
	if err != nil {
		return nil, err
	}

	enabled := "1"
	if gr.Enabled != nil && !*gr.Enabled {
		enabled = "0"
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/groups.php", &url.Values{
		"action": []string{"edit_group"},
		"name":   []string{gr.Name},
		"desc":   []string{gr.Description},
		"status": []string{enabled},
		"id":     []string{strconv.FormatInt(original.ID, 10)},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var updated GroupBasicResponse
	if err = json.NewDecoder(res.Body).Decode(&updated); err != nil {
		return nil, err
	}

	if !updated.Success {
		return nil, fmt.Errorf(updated.Message)
	}

	return c.GetGroup(ctx, gr.Name)
}

// DeleteGroup deletes a group
func (c Client) DeleteGroup(ctx context.Context, name string) error {
	toDelete, err := c.GetGroup(ctx, name)
	if err != nil {
		return err
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/groups.php", &url.Values{
		"action": []string{"delete_group"},
		"id":     []string{strconv.FormatInt(toDelete.ID, 10)},
	})
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var deleted GroupBasicResponse
	if err = json.NewDecoder(res.Body).Decode(&deleted); err != nil {
		return err
	}

	if !deleted.Success {
		return fmt.Errorf(deleted.Message)
	}

	return nil
}
