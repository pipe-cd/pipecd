import { useMutation, UseMutationResult } from "@tanstack/react-query";
import { generateApplicationSealedSecret } from "~/api/piped";

export const useGenerateSealedSecret = (): UseMutationResult<
  string,
  unknown,
  { pipedId: string; data: string; base64Encoding: boolean }
> => {
  return useMutation({
    mutationFn: async (params: {
      pipedId: string;
      data: string;
      base64Encoding: boolean;
    }) => {
      const res = await generateApplicationSealedSecret(params);
      return res.data;
    },
  });
};
