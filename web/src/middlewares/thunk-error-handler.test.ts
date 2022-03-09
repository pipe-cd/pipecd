import configureMockStore from "redux-mock-store";
import thunk from "redux-thunk";
import type { AppState, AppDispatch } from "~/store";
import { thunkErrorHandler } from "./thunk-error-handler";

const middlewares = [thunkErrorHandler, thunk];
const mockStore = configureMockStore<AppState, AppDispatch>(middlewares);

it("passes through exception that not has code and message", async () => {
  const store = mockStore();
  store
    .dispatch(
      (): Promise<void> => {
        throw {};
      }
    )
    .catch(() => null);
  store
    .dispatch(
      (): Promise<void> => {
        throw { code: 1 };
      }
    )
    .catch(() => null);
  store
    .dispatch(
      (): Promise<void> => {
        throw { message: "hello" };
      }
    )
    .catch(() => null);
  expect(store.getActions()).toEqual([]);
});

it("should handle the exception that has code and message", async () => {
  const store = mockStore();
  jest.spyOn(Date, "now").mockImplementation(() => 1);
  store
    .dispatch(
      (): Promise<void> => {
        throw { code: 10, message: "error" };
      }
    )
    .catch(() => null);
  expect(store.getActions()).toEqual([
    {
      type: "toasts/addToast",
      payload: { message: "error", severity: "error" },
    },
  ]);
});
