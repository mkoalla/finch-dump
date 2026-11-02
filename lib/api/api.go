package api

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var API_ROOT = url.URL{
	Scheme: "https",
	Host:   "api.applyfinch.com",
	Path:   "/api/v1",
}

// gormable (?)
type JobEntry struct {
	ID          string     `json:"id,omitempty" gorm:"primaryKey"`
	ApplyUrl    string     `json:"apply_url,omitempty"`
	Category    string     `json:"category,omitempty" gorm:"index"`
	CompanyLogo *string    `json:"company_logo,omitempty"`
	CompanyName string     `json:"company_name,omitempty" gorm:"index"`
	DatePosted  *time.Time `json:"date_posted,omitempty" gorm:"index"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Kind        string     `json:"kind,omitempty" gorm:"index"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	LinkKind    string     `json:"link_kind,omitempty"`
	Locations   []string   `json:"locations,omitempty" gorm:"serializer:json"`
	Pay         *float64   `json:"pay,omitempty"`
	Reposted    bool       `json:"reposted,omitempty"`
	Source      string     `json:"source,omitempty"`
	Terms       []any      `json:"terms,omitempty" gorm:"serializer:json"`
	Title       string     `json:"title,omitempty"`
	Url         string     `json:"url,omitempty"`
	Workplace   *string    `json:"workplace,omitempty"`
}

// not gorm obj
type JobsBody struct {
	Jobs            []JobEntry `json:"jobs,omitempty"`
	LocationApplied bool       `json:"location_applied,omitempty"`
	Major           string     `json:"major,omitempty"`
	Offset          int        `json:"offset,omitempty"`
	PageSize        int        `json:"page_size,omitempty"`
	QCorrected      []any      `json:"q_corrected,omitempty"`
	Total           int        `json:"total,omitempty"`
	Type            any        `json:"type,omitempty"`
	USOnly          bool       `json:"us_only,omitempty"`
	WindowDays      int        `json:"window_days,omitempty"`
}

func Get(path string, query *url.Values) (*http.Response, error) {
	url := API_ROOT.Clone()
	url.Path += path
	url.RawQuery = query.Encode()
	return http.Get(url.String())
}

func Jobs(major string, limit int, us bool) iter.Seq2[*JobEntry, error] {
	// aggregate := make([]JobEntry, 0, max(limit, 256))
	var err error

	query := url.Values{}
	pageLim := min(50, limit)
	if pageLim == -1 {
		pageLim = 50
	}
	query.Set("limit", strconv.Itoa(pageLim)) //enforced
	if major != "" {
		query.Set("major", major)
	}
	if us {
		query.Set("us_only", "1")
	}
	return func(yield func(*JobEntry, error) bool) {
		for i := 0; limit == -1 || i < limit; {
			query.Set("offset", strconv.Itoa(i))
			resp, err := Get("/board/jobs", &query)
			if err != nil {
				break
			}
			cls := sync.OnceFunc(func() { resp.Body.Close() }) //nolint
			if resp.StatusCode != http.StatusOK {
				cls()
				yield(nil, fmt.Errorf("HTTP %d", resp.StatusCode))
				return
			}
			var body JobsBody
			if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
				break
			}
			cls()
			if len(body.Jobs) == 0 {
				break
			}
			for k := range body.Jobs {
				if (i >= limit && limit != -1) || !yield(&body.Jobs[k], nil) {
					return
				}
				i++
			}
			// aggregate = append(aggregate, body.Jobs...)
		}
		if err != nil {
			yield(nil, fmt.Errorf("GET Jobs: %w", err))
		}
	}
}
