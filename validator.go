package main

type Validator struct {
}

//校验包 头数据
func (v *Validator) IsLegal(data []byte) bool {
	return true
}
