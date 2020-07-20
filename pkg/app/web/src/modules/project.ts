import { createSlice } from "@reduxjs/toolkit";

type Project = {
  name: string;
};

const initialState: Project = {
  // TODO: Use fetched project name
  name: "pipe-cd",
};

export const projectSlice = createSlice({
  name: "project",
  initialState,
  reducers: {},
});
