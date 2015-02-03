package tree

import ()

type TreeNode struct {
	//子节点编号
	Id int `json:"id"`
	//父节点编号
	Pid    int         `json:"pid"`
	Text   string      `json:"text"`
	Value  string      `json:"value"`
	Url    string      `json:"url"`
	Icon   string      `json:"icon"`
	Open   bool        `json:"open"`
	Childs []*TreeNode `json:"childs"`
}
