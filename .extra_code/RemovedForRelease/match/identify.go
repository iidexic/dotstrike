package match

type Identifiable interface {
	Identify() string
}

func MatchesID(id string, items ...Identifiable) bool {
	for _, item := range items {
		if item.Identify() == id {
			return true
		}
	}
	return false
}
