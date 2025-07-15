import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as applicationsAPI from "~/api/applications";
import {
  ApplicationGitRepository,
  ApplicationKind,
} from "~/types/applications";

export const useAddApplication = (): UseMutationResult<
  string,
  unknown,
  {
    name: string;
    pipedId: string;
    repo: ApplicationGitRepository.AsObject;
    repoPath: string;
    configFilename?: string;
    kind?: ApplicationKind;
    platformProvider?: string;
    labels: Array<[string, string]>;
    deployTargets?: Array<{ pluginName: string; deployTarget: string }>;
  },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload) => {
      const deployTargetsMap =
        payload.deployTargets?.reduce((all, { pluginName, deployTarget }) => {
          if (!all[pluginName]) all[pluginName] = [];
          all[pluginName].push(deployTarget);
          return all;
        }, {} as Record<string, string[]>) || {};

      const deployTargetsByPluginMap = Object.entries(deployTargetsMap).map(
        ([pluginName, deployTargetsList]) => {
          return [pluginName, { deployTargetsList }] as [
            string,
            { deployTargetsList: string[] }
          ];
        }
      );

      const { applicationId } = await applicationsAPI.addApplication({
        name: payload.name,
        pipedId: payload.pipedId,
        gitPath: {
          repo: payload.repo,
          path: payload.repoPath,
          configFilename: payload.configFilename || "",
          url: "",
        },
        platformProvider: payload.platformProvider,
        kind: payload.kind,
        deployTargetsByPluginMap,
        description: "",
        labelsMap: payload.labels,
      });

      return applicationId;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["applications"] });
    },
  });
};
