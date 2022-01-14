import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
  Dictionary,
  EntityId,
} from "@reduxjs/toolkit";
import { Environment } from "pipe/pkg/app/web/model/environment_pb";
import type { AppState } from "~/store";
import * as envsApi from "~/api/environments";

const MODULE_NAME = "environments";

export const environmentsAdapter = createEntityAdapter<Environment.AsObject>(
  {}
);

const {
  selectById,
  selectAll,
  selectEntities,
  selectIds,
} = environmentsAdapter.getSelectors();

export const selectAllEnvs = (state: AppState): Environment.AsObject[] =>
  selectAll(state.environments);
export const selectEnvIds = (state: AppState): EntityId[] =>
  selectIds(state.environments);
export const selectEnvById = (id: EntityId | undefined) => (
  state: AppState
): Environment.AsObject | undefined =>
  id ? selectById(state.environments, id) : undefined;
export const selectEnvEntities = (
  state: AppState
): Dictionary<Environment.AsObject> => selectEntities(state.environments);

export const fetchEnvironments = createAsyncThunk<Environment.AsObject[], void>(
  "environments/fetchList",
  async () => {
    const { environmentsList } = await envsApi.getEnvironments();
    return environmentsList;
  }
);

export const deleteEnvironment = createAsyncThunk<
  void,
  { environmentId: string }
>(`${MODULE_NAME}/delete`, async (props) => {
  await envsApi.deleteEnvironment(props);
});

export const environmentsSlice = createSlice({
  name: MODULE_NAME,
  initialState: environmentsAdapter.getInitialState(),
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchEnvironments.fulfilled, (state, action) => {
        environmentsAdapter.addMany(
          state,
          action.payload.filter((env) => env.deleted === false)
        );
      })
      .addCase(deleteEnvironment.fulfilled, (state, action) => {
        environmentsAdapter.removeOne(state, action.meta.arg.environmentId);
      });
  },
});

export { Environment } from "pipe/pkg/app/web/model/environment_pb";
