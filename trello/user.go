package trello

import "errors"

type TrelloUser struct {
	ApiKey   string
	ApiToken string
}

// Get the ID of the key's user
func (user *TrelloUser) GetUserId() (string, error) {
	url := BuildRequestUrl("/members/me", user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return "", err
	}

	member_obj := data.(map[string]interface{})
	id, ok := member_obj["id"]
	if !ok {
		return "", errors.New("Field 'id' not found")
	}

	return id.(string), nil
}

// Get a list of all boards the user has access to
func (user *TrelloUser) GetBoards() ([]TrelloBoard, error) {
	// Query the API
	url := BuildRequestUrl("/members/me/boards", user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return []TrelloBoard{}, err
	}

	boards := make([]TrelloBoard, 0)
	for _, board_raw := range data.([]interface{}) {
		board, err := parseBoardData(board_raw)
		if err != nil {
			// TODO: Maybe append to an error list
			continue
		}

		boards = append(boards, board)
	}

	return boards, nil
}

func (user *TrelloUser) GetBoard(id string) (TrelloBoard, error) {
	// Query the API
	url := BuildRequestUrl("/boards/"+id, user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return TrelloBoard{}, err
	}

	board, err := parseBoardData(data)
	if err != nil {
		return TrelloBoard{}, err
	}

	return board, nil
}

func (user *TrelloUser) GetCard(id string) (TrelloCard, error) {
	// Query the API
	url := BuildRequestUrl("/cards/"+id, user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return TrelloCard{}, err
	}

	card := ParseCardData(data)

	return card, nil
}

func (user *TrelloUser) GetUsernameFromId(id string) (string, error) {
	// Query the API
	url := BuildRequestUrl("/members/"+id+"/fullName", user.ApiKey, user.ApiToken)
	data, err := GetUnmarshalledData(url)
	if err != nil {
		return "", err
	}

	name_raw, ok := data.(map[string]interface{})["_value"]
	if !ok {
		return "", errors.New("Field '_value' missing")
	}

	return name_raw.(string), nil
}
