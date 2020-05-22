import { createSlice } from "@reduxjs/toolkit";

type Deployments = {};

const initialState: Deployments = {};

export const deploymentsSlice = createSlice({
  name: "deployments",
  initialState,
  reducers: {},
});
