package main

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/configs"
	"k8s-Management-System/src/controllers"
)

func main() {

	goft.Ignite().Config(configs.NewK8sConfig()).
		Mount("v1", controllers.NewDeploymentCtl()).
		Launch()


}


