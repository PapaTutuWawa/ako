package trello

import "errors"

type TrelloBoard struct {
	Name           string
	Id             string
	Desc           string
	IdOrganization string
	Url            string
	//	Lists          map[string]string
	//      Labels         map[string]string
}

// Try to turn the data of a board into a TrelloBoard
func parseBoardData(data interface{}) (TrelloBoard, error) {
	board := data.(map[string]interface{})

	// Extract all needed fields
	name, ok := board["name"]
	if !ok {
		return TrelloBoard{}, errors.New("Field 'name' not found")
	}
	id, ok := board["id"]
	if !ok {
		return TrelloBoard{}, errors.New("Field 'id' not found")
	}
	desc, ok := board["desc"]
	if !ok {
		return TrelloBoard{}, errors.New("Field 'desc' not found")
	}
	idOrg, ok := board["idOrganization"]
	// Well, the idOrganization field may be empty, so it gets a special treatment
	if !ok || idOrg == nil {
		idOrg = ""
		// Set a hasOrg-Flag to false
	} else {
		idOrg = idOrg.(string)
	}

	url, ok := board["url"]
	if !ok {
		return TrelloBoard{}, errors.New("Field 'url' not found")
	}

	return TrelloBoard{
		Name:           name.(string),
		Id:             id.(string),
		Desc:           desc.(string),
		IdOrganization: idOrg.(string),
		Url:            url.(string),
	}, nil
}

// Retrieve all cards of a board
func (board *TrelloBoard) GetCards(user TrelloUser) ([]TrelloCard, error) {
	url := BuildRequestUrl("/boards/"+board.Id+"/cards", user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return []TrelloCard{}, err
	}

	cards := make([]TrelloCard, 0)
	for _, card_raw := range data.([]interface{}) {
		card := ParseCardData(card_raw)

		cards = append(cards, card)
	}

	return cards, nil
}

func parseListData(data interface{}) (string, string) {
	list_obj := data.(map[string]interface{})

	id, ok := list_obj["id"]
	if !ok {
		id = ""
	}

	name, ok := list_obj["name"]
	if !ok {
		name = ""
	}

	return id.(string), name.(string)
}

// Retrieve the lists of lists of a board
func (board *TrelloBoard) GetLists(user TrelloUser) (map[string]string, error) {
	url := BuildRequestUrl("/boards/"+board.Id+"/lists", user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return map[string]string{}, err
	}

	listmap := make(map[string]string, 0)
	for _, list_obj := range data.([]interface{}) {
		id, name := parseListData(list_obj)
		listmap[id] = name
	}

	return listmap, nil
}

//
func (board *TrelloBoard) GetLabels(user TrelloUser) (map[string]TrelloLabel, error) {
	url := BuildRequestUrl("/boards/"+board.Id+"/labels", user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return map[string]TrelloLabel{}, err
	}

	labelmap := make(map[string]TrelloLabel, 0)
	for _, label_obj := range data.([]interface{}) {
		label := parseLabelData(label_obj)
		labelmap[label.Id] = label
	}

	return labelmap, nil
}
