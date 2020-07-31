import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { AppState } from ".";
import { useSelector } from "react-redux";

export interface LoginState {
  projectName: string | null;
}

const initialState: LoginState = {
  projectName: null,
};

export const loginSlice = createSlice({
  name: "login",
  initialState,
  reducers: {
    setProjectName(state, action: PayloadAction<string>) {
      state.projectName = action.payload;
    },
    clearProjectName(state) {
      state.projectName = null;
    },
  },
});

export const { clearProjectName, setProjectName } = loginSlice.actions;

export const useProjectName = (): string | null => {
  return useSelector<AppState, string | null>(
    (state) => state.login.projectName
  );
};
