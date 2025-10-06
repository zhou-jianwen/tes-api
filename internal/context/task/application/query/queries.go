package query

// Queries ...
type Queries struct {
	List         ListHandler
	Get          GetHandler
	Gather       GatherHandler
	ListAccounts ListAccountsHandler
}

// NewQueries ...
func NewQueries(readModel ReadModel) *Queries {
	return &Queries{
		List:         NewListHandler(readModel),
		Get:          NewGetHandler(readModel),
		Gather:       NewGatherHandler(readModel),
		ListAccounts: NewListAccountsHandler(readModel),
	}
}
