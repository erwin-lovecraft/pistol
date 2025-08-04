package handler

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page, Size int
}

func (p *Pagination) FromRequest(r *http.Request) error {
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return err
		}
		p.Page = page
	}

	if sizeStr := r.URL.Query().Get("size"); sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return err
		}
		p.Size = size
	}

	return nil
}
