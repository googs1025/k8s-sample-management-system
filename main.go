package main

import (
	"github.com/shenyisyn/goft-gin/goft"
	"k8s-Management-System/src/configs"
	"k8s-Management-System/src/controllers"
)

func main() {

	goft.Ignite().Config(
		configs.NewK8sHandler(),  //1
		configs.NewK8sConfig(), //2
		configs.NewK8sMaps(), //3
		configs.NewServiceConfig(), //4
	).
		Mount("v1",controllers.NewDeploymentCtl()).
		Launch()

}


