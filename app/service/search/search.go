package search

import (
	"os"
	"path/filepath"
)

// TreeNode 定义文件树节点结构
type TreeNode struct {
	Name     string      `json:"name"`               // 文件/文件夹名
	Path     string      `json:"path"`               // 绝对路径或相对路径
	IsDir    bool        `json:"isDir"`              // 是否是目录
	Children []*TreeNode `json:"children,omitempty"` // 子节点（仅目录有）
}

// BuildFileTree 递归构建文件树
func BuildFileTree(path string) (*TreeNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Name:  info.Name(),
		Path:  filepath.ToSlash(path), // 统一使用 /，避免前端处理 \\ 问题
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		return node, nil
	}

	// 读取目录下所有条目
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var children []*TreeNode
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childNode, err := BuildFileTree(childPath) // ✅ 递归调用
		if err != nil {
			// 可以选择跳过错误文件，而不是中断整个树
			continue
		}
		children = append(children, childNode)
	}

	node.Children = children
	return node, nil
}
