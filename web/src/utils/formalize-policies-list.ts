import { ProjectRBACPolicy } from "pipecd/web/model/project_pb";
import {
  ACTIONS_KEY,
  KEY_VALUE_SEPARATOR,
  RBAC_ACTION_TYPE_TEXT,
  RBAC_RESOURCE_TYPE_TEXT,
  RESOURCE_ACTION_SEPARATOR,
  RESOURCES_KEY,
  VALUES_SEPARATOR,
} from "~/constants/project";

export const formalizePoliciesList = ({
  policiesList,
}: {
  policiesList: ProjectRBACPolicy.AsObject[];
}): string => {
  const policies: string[] = [];
  policiesList.map((policy) => {
    const resources: string[] = [];
    policy.resourcesList.map((resource) => {
      let rsc = RBAC_RESOURCE_TYPE_TEXT[resource.type];
      if (resource.labelsMap.length > 0) {
        rsc += "{";
        resource.labelsMap.map((label) => {
          rsc += label[0] + ":" + label[1] + ",";
        });
        rsc = rsc.slice(0, -1); // remove last comma
        rsc += "}";
      }
      resources.push(rsc);
    });

    const actions: string[] = [];
    policy.actionsList.map((action) => {
      actions.push(RBAC_ACTION_TYPE_TEXT[action]);
    });

    const resource =
      RESOURCES_KEY + KEY_VALUE_SEPARATOR + resources.join(VALUES_SEPARATOR);
    const action =
      ACTIONS_KEY + KEY_VALUE_SEPARATOR + actions.join(VALUES_SEPARATOR);
    policies.push(resource + RESOURCE_ACTION_SEPARATOR + action);
  });

  return policies.join("\n\n");
};
