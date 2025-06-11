import {
  ProjectRBACPolicy,
  ProjectRBACResource,
} from "pipecd/web/model/project_pb";
import {
  RESOURCE_ACTION_SEPARATOR,
  RESOURCES_KEY,
  KEY_VALUE_SEPARATOR,
  ACTIONS_KEY,
  RESOURCES_NAME_REGEX,
  RESOURCES_LABELS_REGEX,
  TEXT_TO_RBAC_RESOURCE_TYPE,
  VALUES_SEPARATOR,
  TEXT_TO_RBAC_ACTION_TYPE,
} from "~/constants/project";

export const parseRBACPolicies = ({
  policies,
}: {
  policies: string;
}): ProjectRBACPolicy[] => {
  const ps = policies.split("\n\n").filter((p) => p);
  const ret: ProjectRBACPolicy[] = [];
  ps.map((p) => {
    p = p.replace(/\s/g, "");
    const policyResource: ProjectRBACPolicy = new ProjectRBACPolicy();

    // Policy pattern:
    // resources=RESOURCE_NAME{key1:value1,key2:value2};actions=ACTION
    const policy = p.split(RESOURCE_ACTION_SEPARATOR);

    if (
      policy.length !== 2 ||
      policy[0].startsWith(RESOURCES_KEY + KEY_VALUE_SEPARATOR) === false ||
      policy[1].startsWith(ACTIONS_KEY + KEY_VALUE_SEPARATOR) === false
    ) {
      return;
    }

    // Cut the header `resources=`.
    const resources = policy[0].substring(RESOURCES_KEY.length + 1);
    resources.match(RESOURCES_NAME_REGEX)?.map((r) => {
      const resource = new ProjectRBACResource();
      const labels = r.match(RESOURCES_LABELS_REGEX);
      if (labels) {
        resource.clearLabelsMap(); // ensure no labels
        const labelsMap = labels[1].split(",");
        labelsMap.map((l) => {
          const [key, value] = l.split(":");
          resource.getLabelsMap().set(key, value);
        });
      }
      resource.setType(TEXT_TO_RBAC_RESOURCE_TYPE[r.split("{")[0]]);
      policyResource.addResources(resource);
    });

    // Cut the header `actions=`.
    const actions = policy[1].substring(ACTIONS_KEY.length + 1);
    actions.split(VALUES_SEPARATOR).map((v) => {
      policyResource.addActions(TEXT_TO_RBAC_ACTION_TYPE[v]);
    });

    ret.push(policyResource);
  });
  return ret;
};
