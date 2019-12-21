package auth

type StatefulGuard interface {
	Guard
	// 尝试使用给定凭据对用户进行身份验证
	Attempt(credentials map[string]string, login bool) (res interface{}, ok bool)
	// 登录
	Login(user Authenticatable) (data interface{}, err error)
	// 登出
	Logout(token string) error
	//name
	Name() string
}
