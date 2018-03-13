# Ako
An incredibly bad CLI tool for retrieving your cards from Trello

## Config
The configuration file is expected to be at ```$HOME/.config/ako/config.yaml```

For using ako, two keys will need to be set:

- ```key```: This is your Trello API key
- ```token```: This is your Trello API token

Under the key ```aliases```, you can define an alias for any given ID. For example, you can
map the name ```cool_stuff``` to the ID ```ABCDEFGHI``` of a board or card.

## Usage
- ```ako boards```: Lists all your boards
- ```ako cards -id <id>```: Lists all cards of a given board
- ```ako card -id <id>```: Gives information about a given card
- ```ako self -id <id>```: Lista all cards of a given board that are assigned to you
