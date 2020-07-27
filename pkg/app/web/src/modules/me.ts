import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { getMe } from "../api/me";
import { Role } from "pipe/pkg/app/web/model/role_pb";
import { useSelector } from "react-redux";
import { AppState } from ".";

interface Me {
  subject: string;
  avatarUrl: string;
  projectId: string;
  projectRole: Role.ProjectRole;
  isLogin: true;
}

export type MeState = Me | { isLogin: false } | null;

export const fetchMe = createAsyncThunk<Me>("me/fetch", async () => {
  const res = await getMe();
  return { ...res, isLogin: true };
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
        return { isLogin: false };
      });
  },
});

export const selectProjectName = (state: MeState): string => {
  if (state && state.isLogin) {
    return state.projectId;
  }

  return "";
};

export const useMe = (): MeState =>
  useSelector<AppState, MeState>((state) => state.me);
