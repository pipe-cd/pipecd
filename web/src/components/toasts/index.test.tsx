import userEvent from "@testing-library/user-event";
import { render, screen, createStore } from "~~/test-utils";
import { IToast, removeToast } from "~/modules/toasts";
import { Toasts } from "./";

beforeAll(() => {
  jest.useFakeTimers();
});

afterAll(() => {
  jest.useRealTimers();
});

describe("common toast", () => {
  const toast: IToast = {
    id: "toast-1",
    message: "Toast message",
  };

  it("shows toast components if state has toast", () => {
    render(<Toasts />, {
      initialState: {
        toasts: {
          entities: { [toast.id]: toast },
          ids: [toast.id],
        },
      },
    });

    expect(screen.getByText("Toast message")).toBeInTheDocument();
  });

  it("should dispatch remove toast action after 5 sec", () => {
    const store = createStore({
      toasts: {
        entities: { [toast.id]: toast },
        ids: [toast.id],
      },
    });
    render(<Toasts />, {
      store,
    });

    expect(screen.getByText("Toast message")).toBeInTheDocument();

    jest.advanceTimersByTime(4999);

    expect(store.getActions()).toEqual([]);

    jest.advanceTimersByTime(1);

    expect(store.getActions()).toEqual([
      {
        payload: {
          id: "toast-1",
        },
        type: removeToast.type,
      },
    ]);
  });
});

describe("error toast", () => {
  const errorToast: IToast = {
    id: "toast-1",
    message: "Error message",
    severity: "error",
  };

  it("should dispatch remove toast action if click close button", () => {
    const store = createStore({
      toasts: {
        entities: { [errorToast.id]: errorToast },
        ids: [errorToast.id],
      },
    });
    render(<Toasts />, {
      store,
    });

    userEvent.click(screen.getByRole("button", { name: "Close" }));

    expect(store.getActions()).toEqual([
      {
        payload: {
          id: "toast-1",
        },
        type: removeToast.type,
      },
    ]);
  });

  it("should not dispatch remove toast action after 5 sec if toast severity is error", () => {
    const store = createStore({
      toasts: {
        entities: { [errorToast.id]: errorToast },
        ids: [errorToast.id],
      },
    });
    render(<Toasts />, {
      store,
    });

    expect(screen.getByText("Error message")).toBeInTheDocument();

    jest.advanceTimersByTime(4999);

    expect(store.getActions()).toEqual([]);

    jest.advanceTimersByTime(1);

    expect(store.getActions()).toEqual([]);
  });
});
