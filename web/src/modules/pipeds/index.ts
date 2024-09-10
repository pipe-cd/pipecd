import {
  createSlice,
  createAsyncThunk,
  createEntityAdapter,
  EntityId,
} from "@reduxjs/toolkit";
import { Piped } from "pipecd/web/model/piped_pb";
import type { AppState } from "~/store";
import * as pipedsApi from "~/api/piped";

export interface RegisteredPiped {
  id: string;
  key: string;
  isNewKey: boolean;
}

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

export const fetchReleasedVersions = createAsyncThunk<string[]>(
  `${MODULE_NAME}/fetchReleasedVersions`,
  async () => {
    const { versionsList } = await pipedsApi.listReleasedVersions();
    return versionsList;
  }
);

export const fetchBreakingChanges = createAsyncThunk<
  string,
  { projectId: string }
>(`${MODULE_NAME}/fetchBreakingChanges`, async (props) => {
  const { notes } = await pipedsApi.listBreakingChanges({
    projectId: props.projectId,
  });
  return notes;
});

export const fetchPipeds = createAsyncThunk<Piped.AsObject[], boolean>(
  `${MODULE_NAME}/fetchList`,
  async (withStatus: boolean) => {
    const { pipedsList } = await pipedsApi.getPipeds({ withStatus });

    return pipedsList;
  }
);

export const addPiped = createAsyncThunk<
  RegisteredPiped,
  { name: string; desc: string }
>(`${MODULE_NAME}/add`, async (props) => {
  const res = await pipedsApi.registerPiped({
    desc: props.desc,
    name: props.name,
  });
  return { ...res, isNewKey: false };
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

export const restartPiped = createAsyncThunk<void, { pipedId: string }>(
  `${MODULE_NAME}/restart`,
  async ({ pipedId }) => {
    await pipedsApi.restartPiped({ pipedId });
  }
);

export const addNewPipedKey = createAsyncThunk<string, { pipedId: string }>(
  `${MODULE_NAME}/addNewKey`,
  async ({ pipedId }) => {
    const { key } = await pipedsApi.recreatePipedKey({ id: pipedId });
    return key;
  }
);

export const deleteOldKey = createAsyncThunk<void, { pipedId: string }>(
  `${MODULE_NAME}/deleteOldKey`,
  async ({ pipedId }) => {
    await pipedsApi.deleteOldPipedKey({ pipedId });
  }
);

export const editPiped = createAsyncThunk<
  void,
  { pipedId: string; name: string; desc: string }
>(`${MODULE_NAME}/edit`, async ({ pipedId, name, desc }) => {
  await pipedsApi.updatePiped({
    pipedId,
    name,
    desc,
  });
});

export const updatePipedDesiredVersion = createAsyncThunk<
  void,
  { version: string; pipedIds: string[] }
>(`${MODULE_NAME}/updatePipedDesiredVersion`, async ({ version, pipedIds }) => {
  await pipedsApi.updatePipedDesiredVersion({
    version,
    pipedIdsList: pipedIds,
  });
});

export const pipedsSlice = createSlice({
  name: MODULE_NAME,
  initialState: pipedsAdapter.getInitialState<{
    registeredPiped: RegisteredPiped | null;
    updating: boolean;
    releasedVersions: string[];
    breakingChangesNote: string;
  }>({
    registeredPiped: null,
    updating: false,
    releasedVersions: [],
    breakingChangesNote: "",
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
      .addCase(addNewPipedKey.fulfilled, (state, action) => {
        state.registeredPiped = {
          id: action.meta.arg.pipedId,
          key: action.payload,
          isNewKey: true,
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
      })
      .addCase(fetchReleasedVersions.fulfilled, (state, action) => {
        state.releasedVersions = action.payload;
      })
      .addCase(fetchBreakingChanges.fulfilled, (state, action) => {
        state.breakingChangesNote = action.payload;
      });
  },
});

export const { clearRegisteredPipedInfo } = pipedsSlice.actions;
export { Piped, PipedKey } from "pipecd/web/model/piped_pb";
