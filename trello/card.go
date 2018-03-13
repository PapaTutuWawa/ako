package trello

type TrelloCard struct {
	Id     string
	IdList string
	Name   string
	Desc   string
	Labels []string
	Users  []string
}

func ParseCardData(data interface{}) TrelloCard {
	card_cast := data.(map[string]interface{})

	name, ok := card_cast["name"]
	if !ok {
		name = ""
	}
	id, ok := card_cast["id"]
	if !ok {
		id = ""
	}
	desc, ok := card_cast["desc"]
	if !ok {
		desc = ""
	}
	idlist, ok := card_cast["idList"]
	if !ok {
		idlist = ""
	}

	users_arr, ok := card_cast["idMembers"]
	users := make([]string, 0)
	if ok {
		for _, user_obj := range users_arr.([]interface{}) {
			users = append(users, user_obj.(string))
		}
	}

	labels := make([]string, 0)
	label_raw := card_cast["labels"].([]interface{})
	if len(label_raw) > 0 {
		for _, label_obj := range label_raw {
			labelId, ok := label_obj.(map[string]interface{})["id"]
			if !ok {
				labelId = ""
			}

			labels = append(labels, labelId.(string))
		}
	}

	return TrelloCard{
		Id:     id.(string),
		IdList: idlist.(string),
		Name:   name.(string),
		Desc:   desc.(string),
		Labels: labels,
		Users:  users,
	}
}
