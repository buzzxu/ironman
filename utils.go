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
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
//比较密码
func ComparePasswordAndStr(password string,str string)bool  {
	err := bcrypt.CompareHashAndPassword([]byte(str), []byte(password))
	return err == nil
}