package main

import (
	demo "github.com/saschagrunert/demo"
	"github.com/urfave/cli/v2"
	"os/exec"
)

const DEMO_NS = "test"

func main() {
	// cleanupNamespace()
	// if err := setupNamespace(); err != nil {
	// 	panic(err)
	// }

	d := demo.New()

	d.HideVersion = true
	d.Setup(cleanup)
	d.Add(sigstoreDemo(), "sigstore-demo", "sigstore demo")
	d.Add(pspDemo(), "psp-demo", "psps demo")
	d.Run()

}

func sigstoreDemo() *demo.Run {

	exec.Command("kubectl", "delete", "podsecuritypolicy", "restricted").Run()

	r := demo.NewRun(
		"Kubewarden 💖 Sigstore",
	)

	r.Step(demo.S(
		"A simple k8s cluster with Kubewarden",
	), demo.S("kubectl get nodes"))

	r.Step(demo.S("Deploying a policy to verify signatures of container images"),
		demo.S("## https://artifacthub.io/packages/kubewarden/verify-image-signatures/verify-image-signatures"))
	r.Step(nil,
		demo.S(`kwctl scaffold manifest \
  --type ClusterAdmissionPolicy \
  --settings-path policy-settings.yml \
  --title verify-image-signatures \
  registry://ghcr.io/kubewarden/policies/verify-image-signatures:v0.1.5`))
	r.Step(demo.S("Check the policy, make it mutating, include UPDATE operation"),
		demo.S("bat verify-image-signatures-policy.yml"))

	r.Step(demo.S("Apply the policy"),
		demo.S("kubectl apply -f verify-image-signatures-policy.yml"))
	r.Step(nil,
		demo.S("kubectl get deployments -n kubewarden"))
	r.Step(nil,
		demo.S("kubectl rollout status deployment/policy-server-default -n kubewarden"))
	r.Step(nil,
		demo.S("kubectl get clusteradmissionpolicies"))
	r.Step(nil,
		demo.S("kubectl wait --timeout=2m --for=condition=PolicyActive clusteradmissionpolicies --all"))
	r.Step(nil,
		demo.S("kubectl get clusteradmissionpolicies"))

	r.Step(nil,
		demo.S("bat verify-image-signatures-policy.yml"))

	r.Step(demo.S("Deploy a pod with untrusted images"),
		demo.S("bat test_data/goreleaser.yml"))
	r.StepCanFail(nil,
		demo.S("kubectl apply -f test_data/goreleaser.yml"))

	r.Step(demo.S("Deploy a pod with a trusted image"),
		demo.S("bat test_data/jitesoft-alpine.yml"))
	r.StepCanFail(nil,
		demo.S("kubectl apply -f test_data/jitesoft-alpine.yml"))

	return r
}

func pspDemo() *demo.Run {
	r := demo.NewRun(
		"Kubewarden ✨ PSPs",
	)

	r.Step(demo.S(
		"We have a PSP deployed",
	), demo.S("kubectl get podsecuritypolicies"))

	r.Step(demo.S("Transform the PSPs into Kubewarden policies"),
		demo.S("## https://docs.kubewarden.io/tasksDir/psp-migration"))
	r.Step(nil,
		demo.S("./psp-to-kubewarden > psps-migration.yml"))

	r.Step(demo.S("We review them, make changes.. (e.g: `mode: monitor`)"),
		demo.S("bat psps-migration.yml"))

	r.Step(nil,
		demo.S("kubectl apply -f psps-migration.yml"))

	return r
}

func cleanup(ctx *cli.Context) error {
	exec.Command("kubectl", "delete", "clusteradmissionpolicy", "--all").Run()
	exec.Command("kubectl", "delete", "pod", "jitesoft-alpine").Run()
	exec.Command("kubectl", "apply", "-f", "psp-restricted.yml").Run()
	return nil
}
