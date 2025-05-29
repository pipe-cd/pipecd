import { useEffect, useRef } from "react";

export function useInterval(callback: () => void, delay: number | null): void {
  const savedCallback = useRef<() => void>();

  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  useEffect(() => {
    if (delay !== null) {
      const id = setInterval(() => {
        if (savedCallback.current) savedCallback.current();
      }, delay);
      return (): void => {
        clearInterval(id);
      };
    }
    return;
  }, [delay]);
}
