package database

// Info holds the details for the connection.
type Info struct {
	Driver    string
	Username  string
	Password  string
	Database  string
	Hostname  string
	Port      int
	Parameter string
	Timeout   int
}
