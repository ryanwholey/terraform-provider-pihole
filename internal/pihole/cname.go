package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	pihole "github.com/ryanwholey/go-pihole"
)

type CNAMERecordsListResponse struct {
	Data [][]string
}

// ToCNAMERecordList converts a CNAMERecordsListResponse into a CNAMERecordsList object.
func (rr CNAMERecordsListResponse) ToCNAMERecordList() CNAMERecordList {
	list := CNAMERecordList{}

	for _, record := range rr.Data {
		list = append(list, CNAMERecord{
			Domain: record[0],
			Target: record[1],
		})
	}

	return list
}

type CNAMERecord = pihole.CNAMERecord
type CNAMERecordList = pihole.CNAMERecordList

// ListCNAMERecords returns a list of the configured CNAME Pi-hole records
func (c Client) ListCNAMERecords(ctx context.Context) (CNAMERecordList, error) {
	if c.tokenClient != nil {
		return c.tokenClient.LocalCNAME.List(ctx)
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customcname.php", &url.Values{
		"action": []string{"get"},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var cnameRes CNAMERecordsListResponse
	if err = json.NewDecoder(res.Body).Decode(&cnameRes); err != nil {
		return nil, err
	}

	return cnameRes.ToCNAMERecordList(), nil
}

// GetCNAMERecord returns a CNAMERecord for the passed domain if found
func (c Client) GetCNAMERecord(ctx context.Context, domain string) (*CNAMERecord, error) {
	if c.tokenClient != nil {
		record, err := c.tokenClient.LocalCNAME.Get(ctx, domain)
		if err != nil {
			if errors.Is(err, pihole.ErrorLocalCNAMENotFound) {
				return nil, NewNotFoundError(fmt.Sprintf("cname with domain %q not found", domain))
			}

			return nil, err
		}

		return record, nil
	}

	list, err := c.ListCNAMERecords(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range list {
		if strings.EqualFold(r.Domain, domain) {
			return &r, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("cname with domain %q not found", domain))
}

type CreateCNAMERecordResponse struct {
	Success bool
	Message string
}

// CreateCNAMERecord handles CNAME record creation
func (c Client) CreateCNAMERecord(ctx context.Context, record *CNAMERecord) (*CNAMERecord, error) {
	if c.tokenClient != nil {
		return c.tokenClient.LocalCNAME.Create(ctx, record.Domain, record.Target)
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customcname.php", &url.Values{
		"action": []string{"add"},
		"domain": []string{record.Domain},
		"target": []string{record.Target},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var created CreateCNAMERecordResponse
	if err = json.NewDecoder(res.Body).Decode(&created); err != nil {
		return nil, err
	}

	if !created.Success {
		return nil, fmt.Errorf(created.Message)
	}

	return record, err
}

// DeleteCNAMERecord handles CNAME record deletion for the passed domain
func (c Client) DeleteCNAMERecord(ctx context.Context, domain string) error {
	if c.tokenClient != nil {
		return c.tokenClient.LocalCNAME.Delete(ctx, domain)
	}

	record, err := c.GetCNAMERecord(ctx, domain)
	if err != nil {
		return err
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customcname.php", &url.Values{
		"action": []string{"delete"},
		"domain": []string{record.Domain},
		"target": []string{record.Target},
	})
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
