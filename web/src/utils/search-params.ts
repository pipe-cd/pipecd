import { useMemo } from "react";
import { useLocation } from "react-router-dom";
import queryString, { ParsedQuery } from "query-string";

const { parse, stringify } = queryString;

export function useSearchParams(): ParsedQuery<string> {
  const location = useLocation();

  return useMemo<ParsedQuery<string>>((): ParsedQuery<string> => {
    // NOTE: Without specifying arrayFormat, the value will be considered as a string when the length is 1.
    return parse(location.search, { arrayFormat: arrayFormat });
  }, [location.search]);
}

export const stringifySearchParams = stringify;
export const arrayFormat = "index";
