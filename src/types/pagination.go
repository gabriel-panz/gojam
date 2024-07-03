package types

type Pagination struct {
	PageIndex int `url:"page"`
	PageSize  int `url:"size"`
}

func (p Pagination) GetOffset() int {
	return p.PageSize * p.PageIndex
}
