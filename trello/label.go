package trello

import "github.com/fatih/color"

type TrelloLabel struct {
	Id    string
	Name  string
	Color string
}

func parseLabelData(data interface{}) TrelloLabel {
	list_obj := data.(map[string]interface{})

	id, ok := list_obj["id"]
	if !ok {
		id = ""
	}

	name, ok := list_obj["name"]
	if !ok {
		name = ""
	}

	color, ok := list_obj["color"]
	if !ok {
		color = ""
	}

	return TrelloLabel{
		Name:  name.(string),
		Id:    id.(string),
		Color: color.(string),
	}
}

func (label *TrelloLabel) Format() string {
	var c color.Attribute
	switch label.Color {
	case "green":
		c = color.FgGreen
		break
	case "yello":
		c = color.FgYellow
		break
	case "orange":
		c = color.FgRed
		break
	case "red":
		c = color.FgRed
		break
	case "purple":
		c = color.FgMagenta
		break
	case "blue":
		c = color.FgBlue
		break
	case "sky":
		c = color.FgCyan
		break
	case "lime":
		c = color.FgGreen
		break
	case "pink":
		c = color.FgMagenta
		break
	case "black":
		// We can't just print the text in black
		c = color.FgWhite
		break
	default:
		c = color.FgWhite
		break
	}

	return color.New(c).SprintFunc()(label.Name)
}
