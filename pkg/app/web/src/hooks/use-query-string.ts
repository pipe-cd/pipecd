import { useState, useCallback } from "react";
import { parse, stringify } from "query-string";

const getQueryStringValue = ( 
  key: string, 
  queryString = window.location.search
) => { 
  const values = parse(queryString); 
  return values[key];
};

const setQueryStringWithoutPageReload = (qsValue: string) => { 
  const newurl = window.location.protocol + "//" +
                 window.location.host + 
                 window.location.pathname + 
                 qsValue;

  window.history.replaceState({ path: newurl }, "", newurl);
};

const setQueryStringValue = ( 
  key: string, 
  value: string, 
  queryString = window.location.search
) => { 
   const values = parse(queryString); 
   const newQsValue = stringify({ ...values, [key]: value }); 
   setQueryStringWithoutPageReload(`?${newQsValue}`);
};

function useQueryString(key: string, initialValue: string): [string | string[], (a: string | string[]) => void] {
  const [value, setValue] = useState<string | string[]>(getQueryStringValue(key) || initialValue);
  const onSetValue = useCallback(
    newValue => {
      setValue(newValue);
      setQueryStringValue(key, newValue);
    },
    [key]
  );

  return [value, onSetValue];
}

export default useQueryString;
