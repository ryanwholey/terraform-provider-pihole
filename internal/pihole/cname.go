package pihole

import (
	"context"
	"encoding/json"
	"net/url"
)

type CNAMERecordsListResponse struct {
	Data [][]string
}

// ToDNSRecordList converts a CNAMERecordsListResponse into a CNAMERecordsList object.
func (rr CNAMERecordsListResponse) ToDNSRecordList() CNAMERecordList {
	list := CNAMERecordList{}

	for _, record := range rr.Data {
		list = append(list, CNAMERecord{
			Domain: record[0],
			Target: record[1],
		})
	}

	return list
}

type CNAMERecordList []CNAMERecord

type CNAMERecord struct {
	Domain string
	Target string
}

// ListCNAMERecords returns a list of the configured CNAME Pi-hole records
func (c Client) ListCNAMERecords(ctx context.Context) (CNAMERecordList, error) {
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

	return cnameRes.ToDNSRecordList(), nil
}
