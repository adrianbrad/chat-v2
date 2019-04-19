package user

type Permission uint8

const (
	SendMessage Permission = iota
	MuteOthers
	SendMoney
)

var permissions = []string{"send_message", "mute_others", "send_money"}

func (p Permission) String() string {
	return permissions[p]
}
