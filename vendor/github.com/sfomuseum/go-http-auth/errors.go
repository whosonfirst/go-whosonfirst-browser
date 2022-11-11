package auth

// AccountNotExist defines a well-known error for signaling that a given account does not exist.
type AccountNotExist struct{}

func (e AccountNotExist) Error() string {
	return "Account does not exist"
}

// NotLoggedIn defines a well-known error for signaling that the account is not logged in.
type NotLoggedIn struct{}

func (e NotLoggedIn) Error() string {
	return "Not logged in"
}

// NotAuthorized defines a well-known error for signaling that the request is not authorized.
type NotAuthorized struct{}

func (e NotAuthorized) Error() string {
	return "Not authorized"
}
