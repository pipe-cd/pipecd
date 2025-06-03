import {
  MutationCache,
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { FC, PropsWithChildren, useCallback, useMemo } from "react";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { useToast } from "./toast-context";

const QueryClientWrap: FC<PropsWithChildren<unknown>> = ({ children }) => {
  const { addToast } = useToast();

  const handleError = useCallback(
    (err: unknown): void => {
      if (process.env.NODE_ENV === "development") {
        console.error(err);
      }

      if (
        err &&
        typeof err === "object" &&
        "code" in err &&
        err.code &&
        "message" in err &&
        err.message
      ) {
        addToast({ message: err.message as string, severity: "error" });
      } else {
        throw err;
      }
    },
    [addToast]
  );

  const queryClient = useMemo(() => {
    return new QueryClient({
      queryCache: new QueryCache({
        onError: handleError,
      }),
      mutationCache: new MutationCache({
        onError: handleError,
      }),
    });
  }, [handleError]);
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
};

export default QueryClientWrap;
