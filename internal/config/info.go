package config

// Info structures the application settings.
type Info struct {
	Database database
	Casbin   casbin
	Authorizer authorizer
	AppVersion string
	DataDirectory string
	HTTPAddr string
}


// Info holds the details for the connection.
type database struct {
	Driver    string
	Username  string
	Password  string
	Database  string
	Hostname  string
	Port      int
	Parameter string
	Timeout   int64
	URL       string
}

type casbin struct {
	Model  string
	Policy string
	Table  string
}

type authorizer struct {
	CacheExpiration int
}

