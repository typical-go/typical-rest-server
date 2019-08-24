package typirelease

type changelog struct {
	Hash    string
	Message string
}

func (l changelog) String() string {
	return l.Hash + " " + l.Message
}
