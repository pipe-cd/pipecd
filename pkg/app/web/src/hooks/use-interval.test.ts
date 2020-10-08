import { useInterval } from "./use-interval";
import { renderHook } from "@testing-library/react-hooks";

beforeAll(() => {
  jest.useFakeTimers();
});

afterEach(() => {
  jest.clearAllTimers();
});

afterAll(() => {
  jest.useRealTimers();
});

test("should not called callback if passed null to delay", () => {
  const callback = jest.fn();

  renderHook(() => useInterval(callback, null));

  jest.runAllTimers();

  expect(callback).not.toBeCalled();
});

it("should call callback with passed delay", () => {
  const callback = jest.fn();

  renderHook(() => useInterval(callback, 100));

  jest.advanceTimersByTime(99);

  expect(callback).not.toBeCalled();

  jest.advanceTimersByTime(100);

  expect(callback).toHaveBeenCalledTimes(1);
});

it("should call callback with passed delay", () => {
  const callback = jest.fn();

  renderHook(() => useInterval(callback, 100));

  jest.advanceTimersByTime(99);

  expect(callback).not.toBeCalled();

  jest.advanceTimersByTime(100);

  expect(callback).toHaveBeenCalledTimes(1);
});

it("should clear interval on unmount", () => {
  const callback = jest.fn();

  const { unmount } = renderHook(() => useInterval(callback, 100));

  expect(callback).not.toBeCalled();

  unmount();

  jest.runAllTimers();

  expect(callback).not.toBeCalled();
});

it("should update interval if updated delay value", () => {
  const callback = jest.fn();
  let delay = 100;

  const { rerender } = renderHook(() => useInterval(callback, delay));

  jest.advanceTimersByTime(100);
  expect(callback).toHaveBeenCalledTimes(1);

  delay = 300;
  rerender();

  jest.advanceTimersByTime(300);

  expect(callback).toHaveBeenCalledTimes(2);
});
