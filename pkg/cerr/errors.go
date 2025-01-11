package cerr

const (
	ErrInvalidBody       = "INVALID_BODY"
	ErrUserLogined       = "USER_ALREADY_LOGGED"
	ErrEmailNotVerified  = "EMAIL_NOT_VERIFIED"
	ErrEmailExists       = "EMAIL_EXISTS"
	ErrUsernameExists    = "USERNAME_EXISTS"
	UserVerified         = "USER_ALREADY_VERIFIED"
	InvalidVerCode       = "INVALID_VERIFICATION_CODE"
	InvalidPass          = "INVALID_PASSWORD"
	ExpToken             = "EXIRED_TOKEN"
	BelongsToAnotherUser = "BELONGS_TO_ANOTHER_USER"
	ErrMissingCookie     = "MISSING_SESSION"
	ErrLimitOfBookmarks  = "BOOKMARKS_LIMIT"
)
