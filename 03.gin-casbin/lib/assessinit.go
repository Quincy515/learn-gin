package lib

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var E *casbin.Enforcer

func init() {
	initDB()
	adapter, err := gormadapter.NewAdapterByDB(Gorm)
	if err != nil {
		log.Fatal()
	}
	e, err := casbin.NewEnforcer("resources/model_t.conf", adapter)
	if err != nil {
		log.Fatal()
	}
	err = e.LoadPolicy()
	if err != nil {
		log.Fatal()
	}
	E = e
	//initPolicy()
	initPolicyWithDomain()
}

// 从数据库中初始化策略数据 --- 不带租户
func initPolicy() {
	// E.AddPolicy("member", "/depts", "GET")
	// E.AddPolicy("admin", "/depts", "POST")
	// E.AddRoleForUser("zhangsan", "member")
	//return
	// 初始化角色
	m := make([]*RoleRel, 0)
	GetRoles(0, &m, "")
	for _, r := range m {
		_, err := E.AddRoleForUser(r.PRole, r.Role)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 初始化用户角色
	userRoles := GetUserRoles()
	fmt.Println(userRoles)
	for _, user := range userRoles {
		_, err := E.AddRoleForUser(user.UserName, user.RoleName)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 初始化路由角色
	routerRoles := GetRouterRoles()
	for _, rr := range routerRoles {
		_, err := E.AddPolicy(rr.RoleName, rr.RouterUri, rr.RouterMethod)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// initPolicyWithDomain 租户初始化
func initPolicyWithDomain() {
	// 下面这部分是初始化 角色 关系
	// 拼凑出这种格式
	// g, deptadmin, deptupdater, domain1
	// g, deptupdater, deptselecter, domain2
	// 其中 deptselecter 权限最低, 然后是 deptupdater, 最后是 deptadmin
	roles := GetRolesWithDomain() // 获取角色对应
	for _, r := range roles {
		_, err := E.AddRoleForUserInDomain(r.PRole, r.Role, r.Domain)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 初始化用户角色
	userRoles := GetUserRolesWithDomain()
	for _, ur := range userRoles {
		// 增加 domain 参数
		_, err := E.AddRoleForUserInDomain(ur.UserName, ur.RoleName, ur.Domain)
		if err != nil {
			log.Fatal(err)
		}
	}
	// 初始化 路由角色对应关系
	// p deptselecter domain /depts GET
	routerRoles := GetRouterRolesWithDomain()
	for _, rr := range routerRoles {
		_, err := E.AddPolicy(rr.RoleName, rr.Domain, rr.RouterUri, rr.RouterMethod)
		if err != nil {
			log.Fatal(err)
		}
	}
}
