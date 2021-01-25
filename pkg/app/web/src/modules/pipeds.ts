import {
  createSlice,
  createAsyncThunk,
  createEntityAdapter,
  EntityState,
} from "@reduxjs/toolkit";
import { Piped } from "pipe/pkg/app/web/model/piped_pb";
import * as pipedsApi from "../api/piped";

export interface RegisteredPiped {
  id: string;
  key: string;
}

const MODULE_NAME = "pipeds";

const pipedsAdapter = createEntityAdapter<Piped.AsObject>({});

export const {
  selectById,
  selectIds,
  selectAll,
} = pipedsAdapter.getSelectors();

export const fetchPipeds = createAsyncThunk<Piped.AsObject[], boolean>(
  `${MODULE_NAME}/fetchList`,
  async (withStatus: boolean) => {
    const { pipedsList } = await pipedsApi.getPipeds({ withStatus });

    return pipedsList;
  }
);

export const addPiped = createAsyncThunk<
  RegisteredPiped,
  { name: string; desc: string; envIds: string[] }
>(`${MODULE_NAME}/add`, async (props) => {
  const res = await pipedsApi.registerPiped({
    desc: props.desc,
    envIdsList: props.envIds,
    name: props.name,
  });
  return res;
});

export const disablePiped = createAsyncThunk<void, { pipedId: string }>(
  `${MODULE_NAME}/disable`,
  async ({ pipedId }) => {
    await pipedsApi.disablePiped({ pipedId });
  }
);

export const enablePiped = createAsyncThunk<void, { pipedId: string }>(
  `${MODULE_NAME}/enable`,
  async ({ pipedId }) => {
    await pipedsApi.enablePiped({ pipedId });
  }
);

export const recreatePipedKey = createAsyncThunk<string, { pipedId: string }>(
  `${MODULE_NAME}/recreateKey`,
  async ({ pipedId }) => {
    const { key } = await pipedsApi.recreatePipedKey({ id: pipedId });
    return key;
  }
);
export const editPiped = createAsyncThunk<
  void,
  { pipedId: string; name: string; desc: string; envIds: string[] }
>(`${MODULE_NAME}/edit`, async ({ pipedId, name, desc, envIds }) => {
  await pipedsApi.updatePiped({
    pipedId,
    name,
    desc,
    envIdsList: envIds,
  });
});

export const pipedsSlice = createSlice({
  name: MODULE_NAME,
  initialState: pipedsAdapter.getInitialState<{
    registeredPiped: RegisteredPiped | null;
    updating: boolean;
  }>({
    registeredPiped: null,
    updating: false,
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
      })
      .addCase(editPiped.pending, (state) => {
        state.updating = true;
      })
      .addCase(editPiped.rejected, (state) => {
        state.updating = false;
      })
      .addCase(editPiped.fulfilled, (state) => {
        state.updating = false;
      });
  },
});

export const selectPipedsByEnv = (
  state: EntityState<Piped.AsObject>,
  envId: string
): Piped.AsObject[] => {
  return selectAll(state).filter(
    (piped) => piped.envIdsList.includes(envId) && piped.disabled === false
  );
};

export const { clearRegisteredPipedInfo } = pipedsSlice.actions;
export { Piped } from "pipe/pkg/app/web/model/piped_pb";
