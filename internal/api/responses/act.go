package responses

// ActCreated is an act indicates that create action was finished
type ActCreated struct {
	Created bool `json:"created"`
}


// ActDeleted is an act indicates that delete action was finished
type ActDeleted struct {
	Deleted bool `json:"deleted"`
}


// ActUpdated is an act indicates that update action was finished
type ActUpdated struct {
	Updated bool `json:"updated"`
}

// ActGranted is an act indicates that granting action was finished
type ActGranted struct {
	Granted bool `json:"granted"`
}
// ActRevoked is an act indicates that revoking action was finished
type ActRevoked struct {
	Revoked bool `json:"revoked"`
}