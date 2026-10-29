import { useState, useCallback } from "react";
import queryString from "query-string";

export const getQueryStringValue = (
  key: string,
  search = window.location.search
): string | string[] | null => {
  const value = queryString.parse(search)[key];
  // query-string v9 may return null entries for array values such as "?key".
  if (Array.isArray(value)) {
    return value.filter((v): v is string => v !== null);
  }
  return value ?? null;
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
  search = window.location.search
): void => {
  const values = queryString.parse(search);
  const newQsValue = queryString.stringify({ ...values, [key]: value });
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
