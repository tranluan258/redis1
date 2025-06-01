package internal

const (
	SET = "SET"
	GET = "GET"
)

var dataStore = make(map[string]string)

type Command struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
}

func NewCommand(action, key, value string) *Command {
	return &Command{
		Action: action,
		Key:    key,
		Value:  value,
	}
}

func (c *Command) IsValid() bool {
	if c.Action == "" || c.Key == "" {
		return false
	}
	if c.Action == SET && c.Value == "" {
		return false
	}
	return true
}

func (c *Command) Execute() string {
	if !c.IsValid() {
		return "Invalid command"
	}

	switch c.Action {
	case SET:
		return c.executeSet()
	case GET:
		return c.executeGet()
	default:
		return "Unknown command"
	}
}

func (c *Command) executeSet() string {
	dataStore[c.Key] = c.Value
	return "Value set successfully"
}

func (c *Command) executeGet() string {
	value, exists := dataStore[c.Key]
	if !exists {
		return "Key not found"
	}
	return value
}
