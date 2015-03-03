package tree

import (
	_ "fmt"
)

func List2Tree(nodeList []TreeNode) (rootNode *TreeNode) {
	for i, k := range nodeList {
		if k.Id == 0 {
			rootNode = &k
			nodeList = append(nodeList[:i], nodeList[i+1:]...)
			break
		}
	}

	if rootNode == nil {
		rootNode = &TreeNode{
			Id:     0,
			Pid:    0,
			Text:   "根节点",
			Value:  "",
			Url:    "",
			Icon:   "",
			Open:   true,
			Childs: nil}
	}
	iterTree(rootNode, nodeList)
	return rootNode
}
func iterTree(node *TreeNode, nodeList []TreeNode) {
	node.Childs = []*TreeNode{}
	for _, _cnode := range nodeList {
		cnode := _cnode //必须要新建变量，否则都会引用到最后一个元素
		if cnode.Pid == node.Id {
			node.Childs = append(node.Childs, &cnode)
			iterTree(&cnode, nodeList)
		}
	}
}
