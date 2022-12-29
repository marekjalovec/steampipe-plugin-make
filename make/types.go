package make

type SortDir string

const (
	SortDirAsc  SortDir = "asc"
	SortDirDesc SortDir = "desc"
)

type Pagination struct {
	SortBy  string  `json:"sortBy"`
	Limit   int     `json:"limit"`
	SortDir SortDir `json:"sortDir"`
	Offset  int     `json:"offset"`
}

type PaginatedResponse struct {
	Pg Pagination `json:"pg"`
}
