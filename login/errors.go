package login

// LoginErrorID represents single character identifiers for errors occurring during login.
type LoginErrorID byte

const (
	LoginConnectionNotFinished     LoginErrorID = 'n' // Login server closed while client was signing in?
	LoginAlreadyLoggedIn           LoginErrorID = 'a' // Already logged in login server
	LoginAlreadyLoggedInGameServer LoginErrorID = 'c' // Already logged in game server (not IG yet)
	LoginBadVersion                LoginErrorID = 'v' // Client version not supported
	LoginNotPlayer                 LoginErrorID = 'p' // Account not valid
	LoginBanned                    LoginErrorID = 'b' // Account banned
	LoginDisconnectedAccount       LoginErrorID = 'd' // Disconnected a character using this account
	LoginKicked                    LoginErrorID = 'k' // Kicked vs Banned? Ban temp vs ban def?
	LoginServerFull                LoginErrorID = 'w' // Login server full
	LoginOldAccount                LoginErrorID = 'o' // Account created before AG enforced login through AG account?
	LoginOldAccountUseNew          LoginErrorID = 'e' // Game account linked to AG account, must use AG account?
	LoginMaintainAccount           LoginErrorID = 'm' // Account under maintenance
	LoginChooseNickname            LoginErrorID = 'r' // No nickname chosen yet
	LoginNicknameAlreadyUsed       LoginErrorID = 's' // Nickname already used
	LoginAccessDenied              LoginErrorID = 'f' // Incorrect username or password
)

type SelectServerErrorID byte

const (
	SelectServerDown                  SelectServerErrorID = 'd'
	SelectServerFullAtCharacterChoice SelectServerErrorID = 'f'
	SelectServerFull                  SelectServerErrorID = 'F'
	SelectServerCantSelect            SelectServerErrorID = 'r'
	SelectServerShopOtherServer       SelectServerErrorID = 's'
)
