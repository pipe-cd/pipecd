import { waitFor } from "@testing-library/react";
import { setupServer } from "msw/node";
import { GetMeResponse } from "pipecd/web/api_client/service_pb";
import { createHandler } from "~/mocks/create-handler";
import { MemoryRouter, render, screen } from "~~/test-utils";
import { Routes } from "./routes";

const server = setupServer(
  createHandler<GetMeResponse>("/GetMe", () => {
    const response = new GetMeResponse();
    return response;
  })
);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("Pages", () => {
  test("redirect to login page", async () => {
    render(
      <MemoryRouter initialEntries={["/"]} initialIndex={0}>
        <Routes />
      </MemoryRouter>,
      { initialState: { me: { isLogin: false } } }
    );
    await waitFor(() =>
      expect(
        screen.getByRole("textbox", { name: /project name/i })
      ).toBeInTheDocument()
    );
  });
});
