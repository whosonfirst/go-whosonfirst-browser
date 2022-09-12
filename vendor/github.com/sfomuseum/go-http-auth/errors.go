package auth

// AccountNotExist defines a well-known error for signaling that a given account does not exist.
type AccountNotExist struct{}

func (e AccountNotExist) Error() string {
	return "Account does not exist"
}

// AccountNotExist defines a well-known error for signaling that there is no account information.
type NotLoggedIn struct{}

func (e NotLoggedIn) Error() string {
	return "Not logged in"
}
