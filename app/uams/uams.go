package uams

import (
	"encoding/json"
	"errors"
	"github.com/ixre/gof/api"
	"github.com/ixre/gof/web/ui/tree"
	"strconv"
)

var (
	RInternalError = &api.Response{
		Code:    api.RInternalError.Code,
		Message: "内部服务器出错",
	}
	RAccessDenied = &api.Response{
		Code:    api.RAccessDenied.Code,
		Message: "没有权限访问该接口",
	}
	RIncorrectApiParams = &api.Response{
		Code:    api.RIncorrectApiParams.Code,
		Message: "缺少接口参数，请联系技术人员解决",
	}
	RUndefinedApi = &api.Response{
		Code:    api.RUndefinedApi.Code,
		Message: "调用的API名称不正确",
	}
	RNoSuchApp = &api.Response{
		Code:    10096,
		Message: "no such app",
	}
)

type Client struct {
	*api.Client
}

func NewCli(server, key, secret, signType string) *Client {
	return &Client{
		Client: api.NewClient(server, key, secret, signType, checkApiRespErr),
	}
}

// 请求接口
func (c *Client) Post(apiName string, data map[string]string) ([]byte, error) {
	data["version"] = "1.0.0"
	return c.Client.Post(apiName, data)
}

// 用户登陆，返回user_id和user_code,real_name,user_state
func (c *Client) UserLogin(user string, pwd string) (map[string]string, error) {
	data, err := c.Post("user.check_credential", map[string]string{
		"user": user,
		"pwd":  pwd,
		"cred": "user",
	})
	if err == nil {
		var r Result
		err = c.parseResult(data, &r)
		if err == nil && r.ErrCode != 0 {
			err = errors.New(r.ErrMsg)
		}
		return r.Data, err
	}
	return nil, err
}

// 获取所有部门
func (c *Client) GetDeparts(app string) ([]*UamsDepart, error) {
	var d []*UamsDepart
	r, err := c.Post("dept.all", map[string]string{
		"app": app,
	})
	if err == nil {
		err = json.Unmarshal([]byte(r), &d)
	}
	return d, err
}

// 获取部门树
func (c *Client) GetDepartTree(app string) (*tree.TreeNode, error) {
	d := tree.TreeNode{}
	r, err := c.Post("dept.tree", map[string]string{
		"app": app,
	})
	if err == nil {
		err = json.Unmarshal([]byte(r), &d)
	}
	return &d, err
}

// 获取角色
func (c *Client) GetRoles(app string) ([]*Role, error) {
	var d []*Role
	r, err := c.Post("role.all", map[string]string{
		"app": app,
	})
	if err == nil {
		err = json.Unmarshal([]byte(r), &d)
	}
	return d, err
}

// 是否匹配部门
func (c *Client) MatchDept(app string, outerUid string, dept int64) error {
	r, err := c.Post("dept.match", map[string]string{
		"outer_uid": outerUid,
		"dept":      strconv.Itoa(int(dept)),
		"app":       app,
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
func (c *Client) MatchRole(app string, outerUid string, role int64) error {
	r, err := c.Post("role.match", map[string]string{
		"outer_uid": outerUid,
		"role":      strconv.Itoa(int(role)),
		"app":       app,
	})
	if err == nil {
		b, _ := strconv.ParseBool(string(r))
		if !b {
			err = errors.New("not match")
		}
	}
	return err
}

// 插件是否能权限访问资源
func (c *Client) CheckPrivilege(app string, userCode string, resCode string, resUri string) error {
	bytes, err := c.Post("user.privilege", map[string]string{
		"app":       app,
		"user_code": userCode,
		"res_code":  resCode,
		"res_uri":   resUri,
	})
	if err == nil {
		var r Result
		err = c.parseResult(bytes, &r)
	}
	return err
}
func (c *Client) parseResult(bytes []byte, r *Result) error {
	err := json.Unmarshal(bytes, r)
	if err == nil && r.ErrCode != 0 {
		err = errors.New(r.ErrMsg)
	}
	return err
}

// 如果返回接口请求错误, 响应状态码以-10开头
func checkApiRespErr(code int, text string) error {
	switch code {
	case api.RAccessDenied.Code:
		text = RAccessDenied.Message
	case api.RIncorrectApiParams.Code:
		text = RIncorrectApiParams.Message
	case api.RUndefinedApi.Code:
		text = RUndefinedApi.Message
	case RNoSuchApp.Code:
		text = RNoSuchApp.Message
	}
	return errors.New(text)
}

type (
	Result struct {
		ErrCode int               `thrift:"Code,1" db:"Code" json:"ErrCode"`
		ErrMsg  string            `thrift:"Message,2" db:"Message" json:"ErrMsg"`
		Data    map[string]string `thrift:"Data,3" db:"Data" json:"Data"`
	}

	// 菜单
	Res struct {
		ID       int64
		Title    string
		Code     string
		URI      string
		Data     map[string]string
		Children []*Res
	}

	// 资源
	UamsResource struct {
		// ID
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 上级资源
		ParentId int `db:"parent_id"`
		// 应用编号
		AppId int `db:"app_id"`
		// 资源标题
		ResTitle string `db:"res_title"`
		// 资源类型，如menu/button
		ResType string `db:"res_type"`
		// 资源编码
		ResCode string `db:"res_code"`
		// 资源权限可选值，用"|"分割
		ResOption string `db:"res_option"`
		// 资源是否单独开放给用户
		ResOpen int32 `db:"res_open"`
		// 资源图标
		ResIcon string `db:"-"`
		// 资源URL
		ResUri string `db:"res_uri"`
		// 资源数据
		ResData string `db:"res_data"`
		// 资源状态
		ResState int32 `db:"res_state"`
	}

	// 资源
	ResourceTreeNode struct {
		// ID
		ID int `json:"id"`
		// 资源标题
		Title string `json:"title"`
		// 资源类型，如menu/button
		Type string `json:"type"`
		// 资源编码
		Code string `json:"code"`
		// 图标
		Icon string `json:"icon"`
		// 资源URL
		Uri string `json:"uri"`
		// 子资源
		Children []*ResourceTreeNode `json:"children"`
	}

	// 部门
	UamsDepart struct {
		// 部门ID
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 应用编号
		AppId int64 `db:"app_id"`
		// 上级部门编号
		DeptId int64 `db:"dept_id"`
		// 部门名称
		DeptName string `db:"dept_name"`
		// 部门编码
		DeptCode string `db:"dept_code"`
		// 部门描述
		DeptDesc string `db:"dept_desc"`
		// 子部门
		Children []*UamsDepart `db:"-"`
	}

	// 用户字段
	UamsUserField struct {
		// 字段ID
		ID int64 `db:"id" pk:"yes" auto:"yes"`
		// 应用编号
		AppId int64 `db:"app_id"`
		// 字段名称
		FieldName string `db:"field_name"`
		// 字段编码
		FieldCode string `db:"field_code"`
		// 字段类型
		FieldType string `db:"field_type"`
		// 字段可选值
		FieldOption string `db:"field_option"`
		// 字段默认值
		FieldDefault string `db:"field_default"`
		// 字段长度
		FieldLen int64 `db:"field_len"`
		// 字段是否必填
		FieldRequire int32 `db:"field_require"`
		// 排序序号
		SortNum int64 `db:"sort_num"`
		// 字段状态
		FieldState int32 `db:"field_state"`
	}

	// 角色
	Role struct {
		// 角色编号
		ID int `db:"id" pk:"yes" auto:"yes"`
		// 应用编号
		AppId int `db:"app_id"`
		// 角色名称
		RoleName string `db:"role_name"`
		// 角色位值
		RoleFlag int `db:"role_flag"`
		// 角色描述
		RoleDesc string `db:"role_desc"`
		// 默认角色
		IsDefault int `db:"is_default"`
		// 是否启用
		Enabled int `db:"enabled"`
	}
)
