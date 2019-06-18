package main

import (
	"github.com/fatih/color"

	"github.com/flant/kubedog/pkg/utils"
)

func main() {
	t := utils.NewTable(.7, .1, .1, .1)
	t.Header("NAME", "REPLICAS", "UP-TO-DATE", "AVAILABLE")
	t.Raw("deploy/extended-monitoring", "1/1", 1, 1)
	//t.Raw("deploy/extended-monitoring", "1/1", 1, 1, color.RedString("Error: See the server log for details. BUILD FAILED (total time: 1 second)"), color.RedString("Error: An individual language user's deviations from standard language norms in grammar, pronunciation and orthography are sometimes referred to as errors"))
	st := t.SubTable(.4, .1, .3, .1, .1)
	st.Header("NAME", "READY", "STATUS", "RESTARTS", "AGE")
	st.Raws([][]interface{}{
		{"654fc55df-5zs4m", "3/3", "Pulling", "0", "49m", color.RedString("pod/myapp-backend-cbdb856d7-bvplx Failed: Error: ImagePullBackOff"), color.RedString("pod/myapp-backend-cbdb856d7-b6ms8 Failed: Failed to pull image \"ubuntu:kaka\": rpc error: code Unknown desc = Error response from daemon: manifest for ubuntu:kaka not found")},
		{"654fc55df-hsm67", "3/3", color.GreenString("Running") + " -> " + color.RedString("Terminating"), "0", "49m"},
		{"654fc55df-fffff", "3/3", "Ready", "0", "49m"},
	}...)
	t.Raw("deploy/grafana", "1/1", 1, 1)
	t.Raw("deploy/kube-state-metrics", "1/1", 1, 1)
	t.Raw("deploy/madison-proxy-0450d21f50d1e3f3b3131a07bcbcfe85ec02dd9758b7ee12968ee6eaee7057fc", "1/1", 1, 1)
	t.Raw("deploy/madison-proxy-2c5bdd9ba9f80394e478714dc299d007182bc49fed6c319d67b6645e4812b198", "1/1", 1, 1)
	t.Raw("deploy/madison-proxy-9c6b5f859895442cb645c7f3d1ef647e1ed5388c159a9e5f7e1cf50163a878c1", "1/1", 1, "1 (-1)")
	t.Raw("deploy/prometheus-metrics-adapter", "1/1", 1, "1 (-1)")
	t.Raw("sts/mysql", "1/1", 1, "1 (-1)")
	t.Raw("ds/node-exporter", "1/1", 1, "1 (-1)")
	t.Raw("deploy/trickster", "1/1", 1, "1 (-1)")
	t.Render()
}