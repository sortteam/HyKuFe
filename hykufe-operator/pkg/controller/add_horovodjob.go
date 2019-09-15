package controller

import (
	"hykufe-operator/pkg/controller/horovodjob"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, horovodjob.Add)
}
