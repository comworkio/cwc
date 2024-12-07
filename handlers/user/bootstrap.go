package user

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func HandleBootstrap(cmd *cobra.Command, releaseName, nameSpace string, otherValues []string, flagVerbose bool, keepDir bool) {
	repoURL := "https://gitlab.comwork.io/oss/cwcloud/cwcloud-helm.git"
	directory := "./cwcloud-helm-cwc"
	branch := "main"

	// Clone the helm repository
	if err := CloneRepo(repoURL, directory, branch, keepDir); err != nil {
		log.Printf("Error cloning repository: %v", err)
		return
	}

	log.Println("Starting Helm chart installation...")

	patchString := buildPatchString(otherValues)
	log.Printf("Constructed patch string: %s", patchString)

	if err := runHelmInstall(releaseName, directory, nameSpace, patchString); err != nil {
		log.Fatalf("Error running helm command: %v", err)
	}

	log.Println("Helm chart installation completed successfully.")
}

func buildPatchString(otherValues []string) string {
	clusterIP := GetClusterIP()

	var builder strings.Builder
	builder.WriteString("clusterIP=" + clusterIP)

	for _, opt := range otherValues {
		builder.WriteString("," + opt)
	}

	return builder.String()
}

func runHelmInstall(releaseName, directory, nameSpace, patchString string) error {
	helmCommand := "helm"
	helmArgs := []string{
		"install",
		releaseName,
		directory,
		"--create-namespace",
		"--namespace", nameSpace,
		"--set",
		patchString,
	}

	log.Printf("Executing helm command: %s %s", helmCommand, strings.Join(helmArgs, " "))

	helmInstallation := exec.Command(helmCommand, helmArgs...)
	helmInstallation.Stdout = os.Stdout
	helmInstallation.Stderr = os.Stderr

	return helmInstallation.Run()
}

func CloneRepo(repoURL, directory, branch string, keepDir bool) error {
	if !keepDir {
		if _, err := os.Stat(directory); !os.IsNotExist(err) {
			fmt.Printf("Deleting existing directory: %s\n", directory)
			if err := os.RemoveAll(directory); err != nil {
				return fmt.Errorf("failed to delete existing directory: %v", err)
			}
		}
	}

	_, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Progress:      os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	fmt.Println("Repository cloned successfully.")
	return nil
}

func HandleUninstall(cmd *cobra.Command, releaseName string, nameSpace string) {
	log.Println("Starting Helm chart uninstallation...")

	if err := runHelmUninstall(releaseName, nameSpace); err != nil {
		log.Fatalf("Error running helm uninstall command: %v", err)
	}

	log.Println("Helm chart uninstallation completed successfully.")
}

func runHelmUninstall(releaseName, nameSpace string) error {
	helmCommand := "helm"
	helmArgs := []string{
		"uninstall",
		releaseName,
		"--namespace", nameSpace,
	}

	log.Printf("Executing helm command: %s %s", helmCommand, strings.Join(helmArgs, " "))

	helmUninstallation := exec.Command(helmCommand, helmArgs...)
	helmUninstallation.Stdout = os.Stdout
	helmUninstallation.Stderr = os.Stderr

	return helmUninstallation.Run()
}

func HandlePortForward(cmd *cobra.Command, nameSpace string) {
	log.Println("Starting tunnel on CWCloud...")

	if err := runPortForward(nameSpace, "api", 8000); err != nil {
		log.Fatalf("Error running kubectl: %v", err)
	}

	if err := runPortForward(nameSpace, "ui", 3000); err != nil {
		log.Fatalf("Error running kubectl: %v", err)
	}

	log.Println("Now you can go here: http://localhost:3000")
}

func runPortForward(nameSpace string, service string, port int) error {
	kubectlCommand := "kubectl"
	kubectlArgs := []string{
		"-n",
		nameSpace,
		"port-forward",
		"svc/cwcloud-" + service,
		"" + strconv.Itoa(port) + ":" + strconv.Itoa(port),
	}

	log.Printf("Executing kubectl command: %s %s", kubectlCommand, strings.Join(kubectlArgs, " "))

	kubectlPortForward := exec.Command(kubectlCommand, kubectlArgs...)

	return kubectlPortForward.Start()
}
