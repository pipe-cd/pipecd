import {
  createSlice,
  createAsyncThunk,
  createEntityAdapter,
} from "@reduxjs/toolkit";
import { Piped as PipedModel } from "pipe/pkg/app/web/api_client/service_pb";
import { registerPiped, getPipeds } from "../api/piped";

export type Piped = Required<PipedModel.AsObject>;

export interface RegisteredPiped {
  id: string;
  key: string;
}

const pipedsAdapter = createEntityAdapter<Piped>({});

export const { selectById, selectIds } = pipedsAdapter.getSelectors();

export const fetchPipeds = createAsyncThunk<Piped[], boolean>(
  "pipeds/fetchList",
  async (withStatus: boolean) => {
    const { pipedsList } = await getPipeds({ withStatus });
    return pipedsList;
  }
);

export const addPiped = createAsyncThunk<RegisteredPiped, string>(
  "pipeds/add",
  async (desc) => {
    const res = await registerPiped({ desc });
    return res;
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
        pipedsAdapter.addMany(state, action.payload);
      });
  },
});

export const { clearRegisteredPipedInfo } = pipedsSlice.actions;
