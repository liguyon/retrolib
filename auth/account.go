package auth

type LoginErrorID byte

const (
	ErrConnectionNotFinished     LoginErrorID = 'n' // Login server closed while client was signing in?
	ErrAlreadyLoggedIn           LoginErrorID = 'a' // Already logged in login server
	ErrAlreadyLoggedInGameServer LoginErrorID = 'c' // Already logged in game server (not IG yet)
	ErrBadVersion                LoginErrorID = 'v' // Client version not supported
	ErrNotPlayer                 LoginErrorID = 'p' // Account not valid
	ErrBanned                    LoginErrorID = 'b' // Account banned
	ErrDisconnectedAccount       LoginErrorID = 'd' // Disconnected a character using this account
	ErrKicked                    LoginErrorID = 'k' // Kicked vs Banned? Ban temp vs ban def?
	ErrServerFull                LoginErrorID = 'w' // Login server full
	ErrOldAccount                LoginErrorID = 'o' // Account created before AG enforced login though AG account?
	ErrOldAccountUseNew          LoginErrorID = 'e' // Game account linked to AG account, must use AG account?
	ErrMaintainAccount           LoginErrorID = 'm' // Account under maintenance
	ErrChooseNickname            LoginErrorID = 'r' // No nickname chosen yet
	ErrNicknameAlreadyUsed       LoginErrorID = 's' // Nickname already used
	ErrAccessDenied              LoginErrorID = 'f' // Incorrect username or password
)
