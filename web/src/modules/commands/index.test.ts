import { createStore } from "~~/test-utils";
import { server } from "~/mocks/server";
import {
  dummyCommand,
  dummySyncSucceededCommand,
} from "~/__fixtures__/dummy-command";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { commandsSlice, CommandStatus, fetchCommand } from "./";
import { addToast } from "../toasts";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("fetchCommand", () => {
  it("should dispatch addToast if command type is SYNC_APPLICATION and that is succeeded", async () => {
    const store = createStore({});
    await store.dispatch(fetchCommand(dummySyncSucceededCommand.id));

    expect(store.getActions()).toMatchObject([
      { type: fetchCommand.pending.type },
      {
        type: addToast.type,
        payload: {
          to: `/deployments/${dummyDeployment.id}`,
        },
      },
      { type: fetchCommand.fulfilled.type },
    ]);
  });
});

describe("commandsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      commandsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      entities: {},
      ids: [],
    });
  });

  describe("fetchCommand", () => {
    it(`should handle ${fetchCommand.fulfilled.type}`, () => {
      expect(
        commandsSlice.reducer(
          {
            entities: {},
            ids: [],
          },
          {
            type: fetchCommand.fulfilled.type,
            payload: dummyCommand,
          }
        )
      ).toEqual({
        entities: { [dummyCommand.id]: dummyCommand },
        ids: [dummyCommand.id],
      });

      expect(
        commandsSlice.reducer(
          {
            entities: { [dummyCommand.id]: dummyCommand },
            ids: [dummyCommand.id],
          },
          {
            type: fetchCommand.fulfilled.type,
            payload: {
              ...dummyCommand,
              status: CommandStatus.COMMAND_SUCCEEDED,
            },
          }
        )
      ).toEqual({
        entities: {},
        ids: [],
      });
    });
  });
});
