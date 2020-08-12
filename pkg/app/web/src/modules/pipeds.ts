import {
  createSlice,
  createAsyncThunk,
  createEntityAdapter,
} from "@reduxjs/toolkit";
import { Piped as PipedModel } from "pipe/pkg/app/web/model/piped_pb";
import * as pipedsApi from "../api/piped";

export type Piped = Required<PipedModel.AsObject>;

export interface RegisteredPiped {
  id: string;
  key: string;
}

const pipedsAdapter = createEntityAdapter<Piped>({});

export const {
  selectById,
  selectIds,
  selectAll,
} = pipedsAdapter.getSelectors();

export const fetchPipeds = createAsyncThunk<Piped[], boolean>(
  "pipeds/fetchList",
  async (withStatus: boolean) => {
    const { pipedsList } = await pipedsApi.getPipeds({ withStatus });

    return pipedsList;
  }
);

export const addPiped = createAsyncThunk<
  RegisteredPiped,
  { name: string; desc: string }
>("pipeds/add", async (props) => {
  const res = await pipedsApi.registerPiped(props);
  return res;
});

export const disablePiped = createAsyncThunk<void, { pipedId: string }>(
  "pipeds/disable",
  async ({ pipedId }) => {
    await pipedsApi.disablePiped({ pipedId });
  }
);

export const enablePiped = createAsyncThunk<void, { pipedId: string }>(
  "pipeds/enable",
  async ({ pipedId }) => {
    await pipedsApi.enablePiped({ pipedId });
  }
);

export const recreatePipedKey = createAsyncThunk<string, { pipedId: string }>(
  "pipeds/recreateKey",
  async ({ pipedId }) => {
    const { key } = await pipedsApi.recreatePipedKey({ id: pipedId });
    return key;
  }
);

export const pipedsSlice = createSlice({
  name: "pipeds",
  initialState: pipedsAdapter.getInitialState<{
    registeredPiped: RegisteredPiped | null;
  }>({
    registeredPiped: null,
  }),
  reducers: {
    clearRegisteredPipedInfo(state) {
      state.registeredPiped = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(addPiped.fulfilled, (state, action) => {
        state.registeredPiped = action.payload;
      })
      .addCase(fetchPipeds.fulfilled, (state, action) => {
        pipedsAdapter.removeAll(state);
        pipedsAdapter.addMany(state, action.payload);
      })
      .addCase(recreatePipedKey.fulfilled, (state, action) => {
        state.registeredPiped = {
          id: action.meta.arg.pipedId,
          key: action.payload,
        };
      });
  },
});

export const { clearRegisteredPipedInfo } = pipedsSlice.actions;
export { Piped as PipedModel } from "pipe/pkg/app/web/model/piped_pb";
