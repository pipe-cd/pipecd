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

  const handleQueryError = useCallback(
    (err: unknown, query): void => {
      if (query.meta && query.meta.preventGlobalError) {
        return;
      }
      handleError(err);
    },
    [handleError]
  );

  const handleMutationError = useCallback(
    (err: unknown, _variables, _context, mutation) => {
      if (mutation.meta && mutation.meta.preventGlobalError) {
        return;
      }
      handleError(err);
    },
    [handleError]
  );

  const queryClient = useMemo(() => {
    return new QueryClient({
      queryCache: new QueryCache({
        onError: handleQueryError,
      }),
      mutationCache: new MutationCache({
        onError: handleMutationError,
      }),
    });
  }, [handleMutationError, handleQueryError]);
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
};

export default QueryClientWrap;
