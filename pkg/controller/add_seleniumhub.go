package controller

import (
	"github.com/WianVos/selenium_k8s_operator/pkg/controller/seleniumhub"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, seleniumhub.Add)
}
