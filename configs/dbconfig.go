package configs

type DBconfig struct {
	Host             string
	Port             string
	User             string
	Pass             string
	Name             string
	MigrationVersion uint
}
