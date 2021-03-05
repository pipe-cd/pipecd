import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
  Dictionary,
  EntityId,
} from "@reduxjs/toolkit";
import { Environment } from "pipe/pkg/app/web/model/environment_pb";
import { AppState } from ".";
import * as envsApi from "../api/environments";

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

export const addEnvironment = createAsyncThunk<
  void,
  { name: string; desc: string }
>("environments/fetchList", async (props) => {
  await envsApi.AddEnvironment(props);
});

export const environmentsSlice = createSlice({
  name: "environments",
  initialState: environmentsAdapter.getInitialState(),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchEnvironments.fulfilled, (state, action) => {
      environmentsAdapter.addMany(state, action.payload);
    });
  },
});

export { Environment } from "pipe/pkg/app/web/model/environment_pb";
