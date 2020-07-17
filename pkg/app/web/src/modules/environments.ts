import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { Environment as EnvironmentModel } from "pipe/pkg/app/web/model/environment_pb";
import * as envsApi from "../api/environments";

export type Environment = EnvironmentModel.AsObject;

export const environmentsAdapter = createEntityAdapter<Environment>({});

export const {
  selectById,
  selectAll,
  selectEntities,
} = environmentsAdapter.getSelectors();

export const fetchEnvironments = createAsyncThunk<Environment[], void>(
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
