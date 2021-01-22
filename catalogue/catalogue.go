package catalogue

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type file struct {
	Depth    int
	Type     string
	Name     string
	HasChild bool
	Root     bool
}

var Root = "./out"
var ProjectName = "" // 如果有值的话，就不走模板文件中的名字了

// Build 构建工程
func Build(path string) error {
	tree, err := ReadTreeConfig(path)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	err = WriteTree(tree)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	return nil
}

// ReadTreeConfig 读取定义的目录结构的文件
func ReadTreeConfig(path string) (*[]file, error) {
	strs, err := readLine(path)
	if err != nil {
		return nil, err
	}
	return parse(strs)
}

// 读取每一行
func readLine(path string) ([]string, error) {
	if path == "" {
		path = "./template/dic_tree.tpl"
	}
	f, err := os.Open(path)
	defer f.Close()

	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(f)
	var rt []string
	for {
		d, _, e := br.ReadLine()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		if isEmpty(string(d[:])) {
			continue
		}
		rt = append(rt, string(d[:]))
	}
	return rt, nil
}

func isEmpty(v string) bool {
	v = strings.ReplaceAll(v, " ", "")
	if v == "" {
		return true
	}
	return false
}

// 将字符串解析为结构体
func parse(src []string) (*[]file, error) {
	if src == nil || len(src) == 0 {
		return nil, nil
	}
	var rt []file
	for i, v := range src {
		v = strings.TrimLeft(v, " ")  // 去掉左侧的空格
		v = strings.TrimRight(v, " ") // 去掉右侧的空格
		var depth int
		for _, c := range v {
			if c == '-' {
				depth++
			} else {
				break
			}
		}
		if depth == 0 {
			rt = append(rt, file{
				Depth: depth,
				Type:  "",
				Name:  v,
				Root:  true,
			})
			continue
		}
		name := v[depth : len(v)-2]
		typo := v[len(v)-1:]
		if i > 0 && rt[i-1].Depth < depth { // 给父类标志是否包含子节点
			rt[i-1].HasChild = true
		}
		rt = append(rt, file{
			Depth: depth,
			Type:  typo,
			Name:  name,
		})
	}
	return &rt, nil
}

// WriteTree
func WriteTree(tree *[]file) error {
	if len(*tree) == 0 {
		return errors.New("请先定义目录结构模板")
	}
	if (*tree)[0].Root {
		(*tree)[0].Name = ProjectName
	}

	if Root != "" { // 初始化根目录
		err := os.Mkdir(Root, os.ModePerm)
		if err != nil {
			return err
		}
	}
	for i, _ := range *tree {
		path, err := findPath(tree, i)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		if (*tree)[i].Root || (*tree)[i].Type == "d" {
			err = os.Mkdir(path, os.ModePerm)
			if err != nil {
				fmt.Errorf(err.Error())
			}
		} else if (*tree)[i].Type == "f" {
			_, err = os.Create(path)
			if err != nil {
				fmt.Errorf(err.Error())
			}
		}
	}
	return nil // 目前不做异常退出
}

func findPath(tree *[]file, cur int) (string, error) {
	if len(*tree) == 0 {
		return "", errors.New("不能在空的目录树中找到路径")
	}
	var rs string = (*tree)[cur].Name
	depth := (*tree)[cur].Depth
	for i := cur; i >= 0; i-- {
		if (*tree)[i].Depth < depth {
			rs = (*tree)[i].Name + "/" + rs
			depth = (*tree)[i].Depth
		}
	}
	return Root + "/" + rs, nil
}
