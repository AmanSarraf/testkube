package commands

import (
	"fmt"
	"time"

	"github.com/kubeshop/kubtest/pkg/k8sclient"
	"github.com/kubeshop/kubtest/pkg/process"
	"github.com/kubeshop/kubtest/pkg/ui"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

func init() {

}

func NewInstallCmd() *cobra.Command {
	var chart, name, namespace string
	installIngress := false
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Helm chart registry in current kubectl context",
		Long:  `Install can be configured with use of particular `,
		Run: func(cmd *cobra.Command, args []string) {

			ui.Verbose = true

			ui.Logo()

			chart := cmd.Flag("chart").Value.String()
			name := cmd.Flag("name").Value.String()
			namespace := cmd.Flag("namespace").Value.String()
			var err error

			k8sClient, err := k8sclient.ConnectToK8s()
			if err != nil {
				ui.Info("Cannot connect to cluster", err.Error())
				return
			}

			if installIngress {
				err = installIngressController(k8sClient, namespace)
				ui.PrintOnError("installing ingress controller", err)
			}

			_, err = process.Execute("helm", "repo", "add", "kubeshop", "https://kubeshop.github.io/helm-charts")
			ui.ExitOnError("adding kubtest repo", err)

			_, err = process.Execute("helm", "repo", "update")
			ui.ExitOnError("updating helm repositories", err)

			out, err := process.Execute("helm", "install", "--namespace", namespace, name, chart)
			ui.ExitOnError("executing helm install", err)

			ui.Info("Helm install output", string(out))
			err = printDashboardAddress(k8sClient, namespace)
			ui.PrintOnError("getting dashboard address", err)
		},
	}

	cmd.Flags().StringVar(&chart, "chart", "kubeshop/kubtest", "chart name")
	cmd.Flags().StringVar(&name, "name", "kubtest", "installation name")
	cmd.Flags().StringVar(&namespace, "namespace", "default", "namespace where to install")
	cmd.Flags().BoolVarP(&installIngress, "ingress", "i", false, "install ingress if not present in the cluster to expose the endpoint for the dashboard")
	return cmd
}

func installIngressController(k8sClient *kubernetes.Clientset, namespace string) error {
	_, err := process.Execute("helm", "repo", "add", "ingress-nginx", "https://kubernetes.github.io/ingress-nginx")
	if err != nil {
		return err
	}

	_, err = process.Execute("helm", "repo", "update")
	if err != nil {
		return err
	}

	_, err = process.Execute("helm", "install", "--namespace", namespace, IngressControllerName, "ingress-nginx/ingress-nginx")
	if err != nil {
		return err
	}

	err = k8sclient.WaitForPodsReadyForInstance(k8sClient, namespace, IngressControllerName, 30*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func printDashboardAddress(k8sClient *kubernetes.Clientset, namespace string) error {
	//TODO: get automatically the name of the api-server
	ingressAddress, err := k8sclient.GetIngressAddress(k8sClient, IngressApiServerName, namespace)
	if err != nil {
		ui.Info("cannot get the ingress address", err.Error())
		return err
	}

	// TODO: Add version from constant
	completeServerApiAddress := fmt.Sprintf("%s%s/%s/v1/executions/", DashboardURI, ingressAddress, DashboardPrefix)

	ui.Info("kubtest dashboard can be accessed at the address ", completeServerApiAddress)
	ui.Info("a certificate should be added to the ingress to make connection secure")
	return nil
}
