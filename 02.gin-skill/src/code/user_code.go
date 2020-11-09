package code

func IsA(userid int) bool { return userid > 10 } // 是否是付费用户
func IsB(userid int) bool { return userid > 15 } // 是否领取优惠券
func IsC(userid int) bool { return userid > 20 } // 是否连续登陆次数过多

type BoolFunc func(int) bool // 第1步抽取公共函数体

func And(id int, f1 BoolFunc, f2 BoolFunc) bool {
	return f1(id) && f2(id)
}

func Or(id int, f1 BoolFunc, f2 BoolFunc) bool {
	return f1(id) || f2(id)
}

// 根据用户 ID 获取用户等级
func GetUserLevel(userid int) int {
	if And(userid, IsA, IsB) { // 同时满足
		return 1
	}
	if Or(userid, IsB, IsC) { // 满足一个
		return 2
	}
	return 0
}
