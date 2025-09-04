import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as applicationAPI from "~/api/applications";
import {
  ApplicationGitRepository,
  ApplicationKind,
} from "~/types/applications";

export const useUpdateApplication = (): UseMutationResult<
  void,
  unknown,
  {
    applicationId: string;
    name: string;
    pipedId: string;
    repo: ApplicationGitRepository.AsObject;
    repoPath: string;
    configFilename?: string;
    kind?: ApplicationKind;
    platformProvider?: string;
    deployTargets?: Array<{ pluginName: string; deployTarget: string }>;
  }
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

      await applicationAPI.updateApplication({
        applicationId: payload.applicationId,
        name: payload.name,
        pipedId: payload.pipedId,
        platformProvider: payload.platformProvider,
        kind: payload.kind,
        deployTargetsByPluginMap,
        configFilename: payload.configFilename || "",
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["applications"] });
    },
  });
};
