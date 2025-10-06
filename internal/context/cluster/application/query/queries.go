package query

// Queries ...
type Queries struct {
	List ListHandler
}

// NewQueries ...
func NewQueries(readModel ReadModel) *Queries {
	return &Queries{List: NewListHandler(readModel)}
}
