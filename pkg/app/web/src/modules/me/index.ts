import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { Role } from "pipe/pkg/app/web/model/role_pb";
import { getMe } from "../../api/me";
import { useAppSelector } from "../../hooks/redux";

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

export const selectProjectName = (state: { me: MeState }): string => {
  if (state.me && state.me.isLogin) {
    return state.me.projectId;
  }

  return "";
};

export const useMe = (): MeState =>
  useAppSelector<MeState>((state) => state.me);

export { Role } from "pipe/pkg/app/web/model/role_pb";
