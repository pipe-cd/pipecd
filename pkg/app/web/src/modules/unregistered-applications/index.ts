import {
    createAsyncThunk,
} from "@reduxjs/toolkit";
import {
    ApplicationInfo,
} from "pipe/pkg/app/web/model/common_pb";
import * as applicationsAPI from "~/api/applications";

const MODULE_NAME = "unregistered-applications";

export const fetchUnregisteredApplications = createAsyncThunk<ApplicationInfo.AsObject[]>(
    `${MODULE_NAME}/fetchList`,
    async () => {
    const { applicationsList } = await applicationsAPI.getUnregisteredApplications({});
    return applicationsList as ApplicationInfo.AsObject[];
});

export { ApplicationInfo } from "pipe/pkg/app/web/model/common_pb";
