export type Cancelable = {
  clear(): void;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait = 300
): T & Cancelable => {
  let timeout: ReturnType<typeof setTimeout>;
  const debounced = (...args: Parameters<T>): void => {
    const later = (): void => {
      func.apply(this, args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };

  debounced.clear = () => {
    clearTimeout(timeout);
  };

  return debounced as T & Cancelable;
};

export default debounce;
