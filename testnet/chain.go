package main

import "fmt"

type Chain struct {
	DataDir    string
	ID         string
	Validators []*Validator
	Orchestrators []*Orchestrator
}

func (c *Chain) CreateAndInitializeValidators(count uint8) (err error){
	for i := uint8(0); i < count; i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		err = node.init()
		if err != nil {
			return
		}

		c.Validators = append(c.Validators, &node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
	}

	//for _, n := range c.Validators {
	//	if err = addGenesisAccount(n.ConfigDir(), n.Moniker, n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken"); err != nil {
	//		return
	//	}
	//}

	return
}

func (c *Chain) CreateAndInitializeOrchestrators(count uint8) (err error){
	for i := uint8(0); i < count; i++ {
		// create orchestrator
		orchestrator := c.createOrchestrator(i)

		// create keys
		mnemonic, info, err := createMemoryKey();
		if err != nil {
			return err
		}
		orchestrator.KeyInfo = *info
		orchestrator.Mnemonic = mnemonic

		c.Orchestrators = append(c.Orchestrators, &orchestrator)
	}
	return
}

//func (c *Chain) RotateKeys() (err error) {
//	return
//}

func (c *Chain) createValidator(index uint8) (validator Validator) {
	validator = Validator{
		Chain:   c,
		Index:   index,
		Moniker: "gravity",
	}

	return
}

func (c *Chain) createOrchestrator(index uint8) (orchestrator Orchestrator) {
	orchestrator = Orchestrator{
		Chain:   c,
		Index:   index,
	}

	return
}

func (c *Chain) ConfigDir() string {
	return fmt.Sprintf("%s/%s", c.DataDir, c.ID)
}