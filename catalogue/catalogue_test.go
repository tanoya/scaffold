package catalogue

import (
	"fmt"
	"testing"
)

func TestReadTreeConfig(t *testing.T) {
	tree, err := readTreeConfig("../template/dic_tree.tpl")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Print(tree)
}

func TestWriteTree(t *testing.T) {
	tree, err := readTreeConfig("../template/dic_tree.tpl")
	if err != nil {
		t.Error(err.Error())
	}
	err = WriteTree(tree)
	if err != nil {
		t.Error(err.Error())
	}
}
