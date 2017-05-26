package ironman

import "golang.org/x/crypto/bcrypt"

type Optional struct {
	val interface{}
}

func OptionalOf(val interface{}) *Optional {
	return &Optional{val:val}
}
func OptionalOfNil() *Optional {
	return OptionalOf(nil)
}
func (optinal *Optional) Get() interface{} {
	return optinal.val
}
func (optinal *Optional) IsPresent() bool {
	return optinal.val != nil
}

//生成密码
func GeneratePassword(password string) string {
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),15)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
//比较密码
func ComparePasswordAndStr(password ,hash string)bool  {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}