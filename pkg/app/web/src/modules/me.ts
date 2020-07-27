import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { getMe } from "../api/me";
import { Role } from "pipe/pkg/app/web/model/role_pb";

interface Me {
  subject: string;
  avatarUrl: string;
  projectId: string;
  projectRole: Role.ProjectRole;
}

type MeState = Me | null;

export const fetchMe = createAsyncThunk<Me>("me/fetch", async () => {
  const res = await getMe();
  return res;
});

export const meSlice = createSlice({
  name: "me",
  initialState: null as MeState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchMe.fulfilled, (_, action) => {
        return action.payload;
      })
      .addCase(fetchMe.rejected, () => {
        return null;
      });
  },
});

export const selectProjectName = (state: MeState): string => {
  if (state) {
    return state.projectId;
  }

  return "";
};
