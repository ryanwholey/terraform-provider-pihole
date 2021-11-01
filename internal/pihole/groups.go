package pihole

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
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
	ID           string
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
		ID:           strconv.FormatInt(gr.ID, 10),
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

	for _, g := range groupRes.Data {
		log.Print("Modified", g.DateModified)
		log.Print("Description", g.Description)
	}

	return groupRes.ToGroupList(), nil
}

// // CreateDNSRecord creates a pihole DNS record entry
// func (c Client) CreateDNSRecord(ctx context.Context, record *DNSRecord) (*DNSRecord, error) {
// 	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customdns.php", &url.Values{
// 		"action": []string{"add"},
// 		"ip":     []string{record.IP},
// 		"domain": []string{record.Domain},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	res, err := c.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer res.Body.Close()

// 	var created CreateDNSRecordResponse
// 	if err = json.NewDecoder(res.Body).Decode(&created); err != nil {
// 		return nil, err
// 	}

// 	if !created.Success {
// 		return nil, fmt.Errorf(created.Message)
// 	}

// 	return record, nil
// }

// // GetDNSRecord searches the pihole local DNS records for the passed domain and returns a result if found
// func (c Client) GetDNSRecord(ctx context.Context, domain string) (*DNSRecord, error) {
// 	list, err := c.ListDNSRecords(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, r := range list {
// 		if r.Domain == domain {
// 			return &r, nil
// 		}
// 	}

// 	return nil, NewNotFoundError(fmt.Sprintf("record %q not found", domain))
// }

// // UpdateDNSRecord deletes a pihole local DNS record by domain name
// func (c Client) UpdateDNSRecord(ctx context.Context, record *DNSRecord) (*DNSRecord, error) {
// 	current, err := c.GetDNSRecord(ctx, record.Domain)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := c.DeleteDNSRecord(ctx, record.Domain); err != nil {
// 		return nil, err
// 	}

// 	updated, err := c.CreateDNSRecord(ctx, record)
// 	if err != nil {
// 		_, recreateErr := c.CreateDNSRecord(ctx, current)
// 		if err != nil {
// 			return nil, recreateErr
// 		}
// 		return nil, err
// 	}

// 	return updated, nil
// }

// // DeleteDNSRecord deletes a pihole local DNS record by domain name
// func (c Client) DeleteDNSRecord(ctx context.Context, domain string) error {
// 	record, err := c.GetDNSRecord(ctx, domain)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customdns.php", &url.Values{
// 		"action": []string{"delete"},
// 		"ip":     []string{record.IP},
// 		"domain": []string{record.Domain},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	res, err := c.client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer res.Body.Close()

// 	return nil
// }
