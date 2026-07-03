package actions

type ErrNoResultsFound struct{}

func (e *ErrNoResultsFound) Error() string {
	return "no results were found"
}

type ErrNilProvider struct{}

func (e *ErrNilProvider) Error() string {
	return "nil provider"
}
