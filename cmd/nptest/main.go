package main

import (
	"os"

	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"

	"github.com/IntelAI/nodus/pkg/client"
	"github.com/IntelAI/nodus/pkg/config"
	"github.com/IntelAI/nodus/pkg/exec"
)

func main() {
	usage := `nptest - Test Kubernetes Scheduling Scenarios.

Usage:
  nptest --scenario=<config> --pods=<config> --nodes=<config> [--namespace=<ns>]
    [--master=<url> | --kubeconfig=<kconfig>] [--verbose]
  nptest -h | --help

Options:
  -h --help              Show this screen.
  --scenario=<config>    Test scenario config file.
  --pods=<config>        Test pod config file.
  --nodes=<config>       Nodes config file.
  --namespace=<ns>       Namespace to use for tests (will be created if
	                       it does not exist) [default: default]
  --master=<url>         Kubernetes API server URL.
  --kubeconfig=<config>  Kubernetes client config file [default: kconfig].
  --verbose              Enable debug logs.`

	args, _ := docopt.ParseDoc(usage)

	verbose, _ := args.Bool("--verbose")
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	scenarioConfigPath, _ := args.String("--scenario")
	scenario, err := config.ScenarioFromFile(scenarioConfigPath)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("failed to read scenario config")
		os.Exit(1)
	}

	podConfigPath, _ := args.String("--pods")
	podConfig, err := config.PodConfigFromFile(podConfigPath)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("failed to read pod config")
		os.Exit(1)
	}

	nodeConfigPath, _ := args.String("--nodes")
	nodeConfig, err := config.NodeConfigFromFile(nodeConfigPath)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("failed to read node config")
		os.Exit(1)
	}

	// Construct apiserver client
	master, _ := args.String("--master")
	kubeconfigPath, _ := args.String("--kubeconfig")
	k8sclient, err := client.NewK8sClient(master, kubeconfigPath)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("failed to construct kubernetes client")
		os.Exit(1)
	}

	// construct scenario runner
	namespace, _ := args.String("--namespace")
	runner := exec.NewScenarioRunner(k8sclient, namespace, nodeConfig, podConfig)
	err = runner.RunScenario(scenario)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("failed to complete scenario")
		os.Exit(1)
	}
}