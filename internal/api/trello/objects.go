package trello

type TrelloSearchBoards struct {
	Boards []struct {
		Id   string
		Name string
	}
}

type TrelloBoardLists struct {
	Id   string
	Name string
}

type TrelloBasicCard struct {
	Id     string
	Name   string
	Badges struct {
		Comments    int
		Attachments int
	}
	Closed      bool
	Description string `json:"desc"`
	LastUpdate  string `json:"dateLastActivity"`
	Labels      []struct {
		Name  string
		Color string
	}
}

type TrelloAttachment struct {
	Id       string
	IsUpload bool
	Name     string
	Url      string
	Filename string
}

type TrelloAction struct {
	Id   string
	Type string
	Date string
	Data struct {
		Text string
	}
}

type TrelloFullCard struct {
	Id          string
	Name        string
	Attachments []TrelloAttachment
	Actions     []TrelloAction
}
