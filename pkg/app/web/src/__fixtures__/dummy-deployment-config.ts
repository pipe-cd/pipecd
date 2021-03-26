import { ApplicationKind } from "../modules/applications";
import { DeploymentConfigTemplate } from "../modules/deployment-configs";

export const dummyDeploymentConfigTemplates: DeploymentConfigTemplate.AsObject[] = [
  {
    applicationKind: ApplicationKind.KUBERNETES,
    name: "Simple",
    labelsList: [],
    content:
      "# Deploy plain-yaml manifests placing in the application directory without specifying pipeline.\napiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\nspec:\n",
    fileCreationUrl: "",
  },
  {
    applicationKind: ApplicationKind.KUBERNETES,
    name: "Canary",
    labelsList: [],
    content:
      "# Deploy progressively with canary strategy.\napiVersion: pipecd.dev/v1beta1\nkind: KubernetesApp\nspec:\n  pipeline:\n    stages:\n      # Deploy the workloads of CANARY variant. In this case, the number of\n      # workload replicas of CANARY variant is 10% of the replicas number of PRIMARY variant.\n      - name: K8S_CANARY_ROLLOUT\n        with:\n          replicas: 10%\n\n      # Wait a manual approval from a developer on web.\n      - name: WAIT_APPROVAL\n\n      # Update the workload of PRIMARY variant to the new version.\n      - name: K8S_PRIMARY_ROLLOUT\n\n      # Destroy all workloads of CANARY variant.\n      - name: K8S_CANARY_CLEAN\n",
    fileCreationUrl: "",
  },
];

export function deploymentConfigTemplateFromObject(
  o: DeploymentConfigTemplate.AsObject
): DeploymentConfigTemplate {
  const template = new DeploymentConfigTemplate();
  template.setName(o.name);
  template.setApplicationKind(o.applicationKind);
  template.setContent(o.content);
  template.setFileCreationUrl(o.fileCreationUrl);
  template.setLabelsList(o.labelsList);
  return template;
}
