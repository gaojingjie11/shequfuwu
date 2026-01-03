package model

import "time"

// SysUser 对应数据库 sys_user 表 [cite: 39]
type SysUser struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // JSON返回时隐藏密码
	RealName  string    `json:"real_name"`
	Mobile    string    `gorm:"unique" json:"mobile"` // 手机号唯一
	Age       int       `json:"age"`
	Gender    int       `json:"gender"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Balance   float64   `json:"balance"`
	Role      string    `json:"role"`   // user, admin
	Status    int       `json:"status"` // 1正常 0冻结
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SysUser) TableName() string {
	return "sys_user"
}

// SysRole 角色表
type SysRole struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`               // 角色名称
	Code      string    `gorm:"unique" json:"code"` // 角色标识 admin, resident
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysRole) TableName() string {
	return "sys_role"
}

// SysMenu 菜单表 (权限)
type SysMenu struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ParentID  int64     `json:"parent_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Component string    `json:"component"`
	Sort      int       `json:"sort"`
	Type      int       `json:"type"` // 1目录 2菜单 3按钮
	CreatedAt time.Time `json:"created_at"`
}

func (SysMenu) TableName() string {
	return "sys_menu"
}

// SysUserRole 用户角色关联
type SysUserRole struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}

// SysRoleMenu 角色菜单关联
type SysRoleMenu struct {
	ID     int64 `gorm:"primaryKey" json:"id"`
	RoleID int64 `json:"role_id"`
	MenuID int64 `json:"menu_id"`
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
