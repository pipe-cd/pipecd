import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { Piped as PipedModel } from "pipe/pkg/app/web/model/piped_pb";
import { registerPiped } from "../api/piped";

export type Piped = Required<PipedModel.AsObject>;

export interface RegisteredPiped {
  id: string;
  key: string;
}

type Pipeds = {
  registeredPiped: RegisteredPiped | null;
};

const initialState: Pipeds = {
  registeredPiped: null,
};

export const addPiped = createAsyncThunk<RegisteredPiped, string>(
  "pipeds/add",
  async (desc) => {
    const res = await registerPiped({ desc });
    return res;
  }
);

export const pipedsSlice = createSlice({
  name: "pipeds",
  initialState,
  reducers: {
    clearRegisteredPipedInfo(state) {
      state.registeredPiped = null;
    },
  },
  extraReducers: (builder) => {
    builder.addCase(addPiped.fulfilled, (state, action) => {
      state.registeredPiped = action.payload;
    });
  },
});

export const { clearRegisteredPipedInfo } = pipedsSlice.actions;
