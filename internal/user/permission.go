package user

type Permission uint8

const (
	SendMessage Permission = iota
)

var permissions = []string{"send_message"}

func (p Permission) String() string {
	return permissions[p]
}
