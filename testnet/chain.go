package main

import "fmt"

type Chain struct {
	DataDir string
	ID string
	Nodes []*Node
}

func (c *Chain) CreateAndInitializeNodes(count uint8) (err error){
	for i := uint8(0); i < count; i++ {
		node := c.createNode(i)

		err = node.init()
		if err != nil {
			return
		}

		c.Nodes = append(c.Nodes, &node)

		if err := node.createKey("val"); err != nil {
			return err
		}
	}

	return
}

func (c *Chain) createNode(index uint8) (node Node) {
	node = Node{
		Chain:   c,
		Index:   index,
		Moniker: "gravity",
	}

	return
}

func (c *Chain) ConfigDir() string {
	return fmt.Sprintf("%s/%s", c.DataDir, c.ID)
}