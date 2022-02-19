package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type ListDomainsOptions struct {
	Type string
}

const (
	// DomainTypeAllow indicates a domain that is added to the allow list
	DomainTypeAllowExact int = 0
	// DomainTypeDeny indicates a domain that is added to the deny list
	DomainTypeDenyExact int = 1
	// DomainTypeAllowWildcard indicates a wildcard domain added to the allow list
	DomainTypeAllowWildcard int = 2
	// DomainTypeDenyWildcard indicates a wildcard domain added to the deny list
	DomainTypeDenyWildcard int = 3
)

var IntToDomainType = map[int]string{
	DomainTypeAllowExact:    DomainOptionsAllow,
	DomainTypeDenyExact:     DomainOptionsDeny,
	DomainTypeAllowWildcard: DomainOptionsAllow,
	DomainTypeDenyWildcard:  DomainOptionsDeny,
}

const (
	// DomainOptionsAllow is a filter option corresponding to domains on the allowed list
	DomainOptionsAllow string = "allow"
	// DomainOptionsDeny is a filter option corresponding to domains on the deny list
	DomainOptionsDeny string = "deny"
)

type DomainResponseList struct {
	Data []*DomainResponse
}

type DomainResponse struct {
	ID           int64   `json:"id"`
	Type         int     `json:"type"`
	Enabled      int     `json:"enabled"`
	Domain       string  `json:"domain"`
	Comment      string  `json:"comment"`
	DateAdded    int64   `json:"date_added"`
	DateModified int64   `json:"date_modified"`
	Groups       []int64 `json:"groups"`
}

func (l DomainResponseList) ToDomainList() DomainList {
	list := make(DomainList, len(l.Data))

	for i, d := range l.Data {
		list[i] = d.ToDomain()
	}

	return list
}

type Domain struct {
	ID           int64
	Type         string
	Enabled      bool
	Domain       string
	Comment      string
	DateAdded    time.Time
	DateModified time.Time
	Wildcard     bool
	GroupIDs     []int64
}

type DomainList []*Domain

func (d DomainResponse) ToDomain() *Domain {
	return &Domain{
		ID:           d.ID,
		Type:         IntToDomainType[d.Type],
		Enabled:      d.Enabled == 1,
		Domain:       d.Domain,
		Comment:      d.Comment,
		DateAdded:    time.Unix(d.DateAdded, 0),
		DateModified: time.Unix(d.DateModified, 0),
		Wildcard:     d.Type == DomainTypeAllowWildcard || d.Type == DomainTypeDenyWildcard,
		GroupIDs:     d.Groups,
	}
}

// ListDomains returns a list of domains
func (c Client) ListDomains(ctx context.Context, opts ListDomainsOptions) (DomainList, error) {
	if c.tokenClient != nil {
		return nil, fmt.Errorf("%w: list domains", ErrNotImplementedTokenClient)
	}

	values := &url.Values{
		"action": []string{"get_domains"},
	}

	if opts.Type != "" {
		if opts.Type == DomainOptionsAllow {
			values.Add("showtype", "white")
		} else if opts.Type == DomainOptionsDeny {
			values.Add("showtype", "black")
		} else {
			return nil, fmt.Errorf("unknown type passed to ListDomains: %s", opts.Type)
		}
	}

	req, err := c.RequestWithSession(ctx, "POST", "/admin/scripts/pi-hole/php/groups.php", values)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var domainRes DomainResponseList
	if err = json.NewDecoder(res.Body).Decode(&domainRes); err != nil {
		return nil, err
	}

	return domainRes.ToDomainList(), nil
}
