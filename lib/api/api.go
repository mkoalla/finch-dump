package api

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

var API_ROOT = url.URL{
	Scheme: "https",
	Host:   "api.applyfinch.com",
	Path:   "/api/v1",
}

type JobEntry struct {
	ApplyUrl    string   `json:"apply_url,omitempty"`    //": "https://careers.se.com/jobs/139399",
	Category    string   `json:"category,omitempty"`     //": "FULL_TIME",
	CompanyLogo string   `json:"company_logo,omitempty"` //": null,
	CompanyName string   `json:"company_name,omitempty"` //": "Schneider Electric",
	DatePosted  string   `json:"date_posted,omitempty"`  //": "2026-11-02T00:00:00+00:00",
	Deadline    string   `json:"deadline,omitempty"`     //": null,
	ID          string   `json:"id,omitempty"`           //": "091ef9b9-4ac9-471d-9d2c-1b81f3309c2c",
	Kind        string   `json:"kind,omitempty"`         //": "newgrad",
	LastSeenAt  string   `json:"last_seen_at,omitempty"` //": "2026-11-02T00:58:05.728777+00:00",
	LinkKind    string   `json:"link_kind,omitempty"`    //": "posting",
	Locations   []string `json:"locations,omitempty"`    //": [
	// "Boston, Massachusetts"
	// ],
	Pay       float64 `json:"pay,omitempty"`       //: null,
	Reposted  bool    `json:"reposted,omitempty"`  //: false,
	Source    string  `json:"source,omitempty"`    //: "ats_icims",
	Terms     []any   `json:"terms,omitempty"`     //: [],
	Title     string  `json:"title,omitempty"`     //: "Energy Tech Accelerator Graduate Program - Boston",
	Url       string  `json:"url,omitempty"`       //: "https://careers.se.com/jobs/139399",
	Workplace string  `json:"workplace,omitempty"` //: null

}
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

func Jobs(major string, limit int) iter.Seq2[*JobEntry, error] {
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
	return func(yield func(*JobEntry, error) bool) {
		for i := 0; limit == -1 || i < limit; {
			query.Set("offset", strconv.Itoa(i))
			resp, err := Get("/board/jobs", &query)
			if err != nil {
				break
			}
			cls := sync.OnceFunc(func() { resp.Body.Close() })
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
