package domain

type SetCommand struct {
	Key         string
	Value       string
	Expiry_time int
	Condition   string
}

type GetCommand struct {
	key string
}
