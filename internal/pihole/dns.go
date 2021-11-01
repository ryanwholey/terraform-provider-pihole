package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type DNSRecordsListResponse struct {
	Data [][]string
}

// ToDNSRecordList converts a DNSRecordsListResponse into a DNSRecordList object.
func (rr DNSRecordsListResponse) ToDNSRecordList() DNSRecordList {
	list := DNSRecordList{}

	for _, record := range rr.Data {
		list = append(list, DNSRecord{
			Domain: record[0],
			IP:     record[1],
		})
	}

	return list
}

type DNSRecordList []DNSRecord

type DNSRecord struct {
	Domain string
	IP     string
}

// ListDNSRecords Returns the list of custom DNS records configured in pihole
func (c Client) ListDNSRecords(ctx context.Context) (DNSRecordList, error) {
	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customdns.php", &url.Values{
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

	var dnsRes DNSRecordsListResponse
	if err = json.NewDecoder(res.Body).Decode(&dnsRes); err != nil {
		return nil, err
	}

	return dnsRes.ToDNSRecordList(), nil
}

type CreateDNSRecordResponse struct {
	Success bool
	Message string
}

// CreateDNSRecord creates a pihole DNS record entry
func (c Client) CreateDNSRecord(ctx context.Context, record *DNSRecord) (*DNSRecord, error) {
	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customdns.php", &url.Values{
		"action": []string{"add"},
		"ip":     []string{record.IP},
		"domain": []string{record.Domain},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var created CreateDNSRecordResponse
	if err = json.NewDecoder(res.Body).Decode(&created); err != nil {
		return nil, err
	}

	if !created.Success {
		return nil, fmt.Errorf(created.Message)
	}

	return record, nil
}

// GetDNSRecord searches the pihole local DNS records for the passed domain and returns a result if found
func (c Client) GetDNSRecord(ctx context.Context, domain string) (*DNSRecord, error) {
	list, err := c.ListDNSRecords(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range list {
		if r.Domain == domain {
			return &r, nil
		}
	}

	return nil, NewNotFoundError(fmt.Sprintf("record %q not found", domain))
}

// UpdateDNSRecord deletes a pihole local DNS record by domain name
func (c Client) UpdateDNSRecord(ctx context.Context, record *DNSRecord) (*DNSRecord, error) {
	current, err := c.GetDNSRecord(ctx, record.Domain)
	if err != nil {
		return nil, err
	}

	if err := c.DeleteDNSRecord(ctx, record.Domain); err != nil {
		return nil, err
	}

	updated, err := c.CreateDNSRecord(ctx, record)
	if err != nil {
		_, recreateErr := c.CreateDNSRecord(ctx, current)
		if err != nil {
			return nil, recreateErr
		}
		return nil, err
	}

	return updated, nil
}

// DeleteDNSRecord deletes a pihole local DNS record by domain name
func (c Client) DeleteDNSRecord(ctx context.Context, domain string) error {
	record, err := c.GetDNSRecord(ctx, domain)
	if err != nil {
		return err
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/customdns.php", &url.Values{
		"action": []string{"delete"},
		"ip":     []string{record.IP},
		"domain": []string{record.Domain},
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
