package response

type Pagination struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	PerPage     int64 `json:"per_page,omitempty"`
	HasNextPage bool  `json:"has_next_page"`
	Total       int64 `json:"total"`
}
