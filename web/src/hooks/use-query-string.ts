import { useState, useCallback } from "react";
import queryString from "query-string";

const { parse, stringify } = queryString;

export const getQueryStringValue = (
  key: string,
  queryString = window.location.search
): string | string[] | null => {
  const value = parse(queryString)[key];
  if (Array.isArray(value)) {
    return value.filter((v): v is string => v !== null);
  }
  return value;
};

const setQueryStringWithoutPageReload = (qsValue: string): void => {
  const newurl =
    window.location.protocol +
    "//" +
    window.location.host +
    window.location.pathname +
    qsValue;

  window.history.replaceState({ path: newurl }, "", newurl);
};

const setQueryStringValue = (
  key: string,
  value: string,
  queryString = window.location.search
): void => {
  const values = parse(queryString);
  const newQsValue = stringify({ ...values, [key]: value });
  setQueryStringWithoutPageReload(`?${newQsValue}`);
};

function useQueryString(
  key: string,
  initialValue: string
): [string | string[], (a: string | string[]) => void] {
  const [value, setValue] = useState<string | string[]>(
    getQueryStringValue(key) || initialValue
  );
  const onSetValue = useCallback(
    (newValue) => {
      setValue(newValue);
      setQueryStringValue(key, newValue);
    },
    [key]
  );

  return [value, onSetValue];
}

export default useQueryString;
