package uams

import (
	"encoding/json"
	"errors"
	"github.com/jsix/gof/web/ui/tree"
	"strconv"
)

// 用户登陆，返回user_id和user_code,real_name,user_state
func UserLogin(user string, pwd string) (map[string]string, error) {
	data, err := Post("user.login", map[string]string{
		"user": user,
		"pwd":  pwd,
	})
	if err == nil {
		var r Result
		json.Unmarshal(data, &r)
		if r.ErrCode != 0 {
			err = errors.New(r.ErrMsg)
		}
		return r.Data, err
	}
	return nil, err
}

func GetDeparts() ([]*UamsDepart, error) {
	var d []*UamsDepart
	r, err := Post("dept.all", nil)
	if err == nil {
		err = json.Unmarshal([]byte(r), &d)
	}
	return d, err
}

func GetDepartTree() (*tree.TreeNode, error) {
	d := tree.TreeNode{}
	r, err := Post("dept.tree", nil)
	if err == nil {
		err = json.Unmarshal([]byte(r), &d)
	}
	return &d, err
}

func GetRoles() ([]*UamsRole, error) {
	var d []*UamsRole
	r, err := Post("role.all", nil)
	if err == nil {
		err = json.Unmarshal([]byte(r), &d)
	}
	return d, err
}

// 是否匹配部门
func MatchDept(outerUid string, dept int64) error {
	r, err := Post("dept.match", map[string]string{
		"outer_uid": outerUid,
		"dept":      strconv.Itoa(int(dept)),
	})
	if err == nil {
		b, _ := strconv.ParseBool(string(r))
		if !b {
			err = errors.New("not match")
		}
	}
	return err
}

// 是否匹配角色
func MatchRole(outerUid string, role int64) error {
	r, err := Post("role.match", map[string]string{
		"outer_uid": outerUid,
		"role":      strconv.Itoa(int(role)),
	})
	if err == nil {
		b, _ := strconv.ParseBool(string(r))
		if !b {
			err = errors.New("not match")
		}
	}
	return err
}
