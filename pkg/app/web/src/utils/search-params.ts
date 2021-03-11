import { useMemo } from "react";
import { useLocation } from "react-router-dom";
import { parse, stringify, ParsedQuery } from "query-string";

export function useSearchParams(): ParsedQuery<string> {
  const location = useLocation();

  return useMemo<ParsedQuery<string>>((): ParsedQuery<string> => {
    return parse(location.search);
  }, [location.search]);
}

export const stringifySearchParams = stringify;
