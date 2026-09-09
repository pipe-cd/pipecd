import debounce from "./debounce";

beforeEach(() => {
  jest.useFakeTimers();
});

afterEach(() => {
  jest.useRealTimers();
});

describe("debounce", () => {
  it("does not call the function before the wait period elapses", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    debounced();
    expect(fn).not.toHaveBeenCalled();

    jest.advanceTimersByTime(299);
    expect(fn).not.toHaveBeenCalled();
  });

  it("calls the function after the wait period", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    debounced();
    jest.advanceTimersByTime(300);
    expect(fn).toHaveBeenCalledTimes(1);
  });

  it("uses the default wait of 300ms when no wait is specified", () => {
    const fn = jest.fn();
    const debounced = debounce(fn);

    debounced();
    jest.advanceTimersByTime(299);
    expect(fn).not.toHaveBeenCalled();

    jest.advanceTimersByTime(1);
    expect(fn).toHaveBeenCalledTimes(1);
  });

  it("resets the timer on repeated calls and fires only once", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    debounced();
    jest.advanceTimersByTime(200);
    debounced();
    jest.advanceTimersByTime(200);
    debounced();
    jest.advanceTimersByTime(300);

    expect(fn).toHaveBeenCalledTimes(1);
  });

  it("passes the most recent arguments to the underlying function", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    debounced("first");
    jest.advanceTimersByTime(100);
    debounced("second");
    jest.advanceTimersByTime(300);

    expect(fn).toHaveBeenCalledWith("second");
  });

  it("can be called again after the wait period fires", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    debounced();
    jest.advanceTimersByTime(300);
    debounced();
    jest.advanceTimersByTime(300);

    expect(fn).toHaveBeenCalledTimes(2);
  });

  it("clear() prevents the pending call from executing", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    debounced();
    jest.advanceTimersByTime(100);
    debounced.clear();
    jest.advanceTimersByTime(300);

    expect(fn).not.toHaveBeenCalled();
  });

  it("clear() is a no-op when called with no pending invocation", () => {
    const fn = jest.fn();
    const debounced = debounce(fn, 300);

    expect(() => debounced.clear()).not.toThrow();
    expect(fn).not.toHaveBeenCalled();
  });
});
