import {
  createSlice,
  createAsyncThunk,
  createEntityAdapter,
  EntityId,
} from "@reduxjs/toolkit";
import { Piped } from "pipecd/web/model/piped_pb";
import type { AppState } from "~/store";
import * as pipedsApi from "~/api/piped";

const MODULE_NAME = "pipeds";

const pipedsAdapter = createEntityAdapter<Piped.AsObject>({});

const { selectById, selectIds, selectAll } = pipedsAdapter.getSelectors();

export const selectPipedById = (id?: EntityId | null) => (
  state: AppState
): Piped.AsObject | undefined =>
  id ? selectById(state.pipeds, id) : undefined;
export const selectPipedIds = (state: AppState): EntityId[] =>
  selectIds(state.pipeds);
export const selectAllPipeds = (state: AppState): Piped.AsObject[] =>
  selectAll(state.pipeds);

export const fetchPipeds = createAsyncThunk<Piped.AsObject[], boolean>(
  `${MODULE_NAME}/fetchList`,
  async (withStatus: boolean) => {
    const { pipedsList } = await pipedsApi.getPipeds({ withStatus });

    return pipedsList;
  }
);

export const pipedsSlice = createSlice({
  name: MODULE_NAME,
  initialState: pipedsAdapter.getInitialState({}),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchPipeds.fulfilled, (state, action) => {
      pipedsAdapter.removeAll(state);
      pipedsAdapter.addMany(state, action.payload);
    });
  },
});

export { Piped, PipedKey } from "pipecd/web/model/piped_pb";
