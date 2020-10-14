import { dummyCommand } from "../__fixtures__/dummy-command";
import { commandsSlice, CommandStatus, fetchCommand } from "./commands";

describe("commandsSlice reducer", () => {
  it("should handle initial state", () => {
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
