package pagination

import (
	"context"
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func SetPaginationContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit <= 0 {
			limit = -1
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset <= 0 {
			offset = -1
		}

		pagination := &Pagination{
			Limit:  limit,
			Offset: offset,
		}

		ctx := context.WithValue(r.Context(), "pagination", pagination)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Get(ctx context.Context) *Pagination {
	val := ctx.Value("pagination")
	if pagination, ok := val.(*Pagination); ok {
		return pagination
	}
	return nil
}
