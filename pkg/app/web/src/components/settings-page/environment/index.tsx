import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@material-ui/core";
import { unwrapResult } from "@reduxjs/toolkit";
import { FC, memo, useCallback, useEffect } from "react";
import { DELETE_ENVIRONMENT_SUCCESS } from "~/constants/toast-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  deleteEnvironment,
  fetchEnvironments,
  selectEnvIds,
} from "~/modules/environments";
import { addToast } from "~/modules/toasts";
import { DeleteEnvironmentDialog } from "./components/delete-confirm-dialog";
import { EnvironmentListItem } from "./components/environment-list-item";

export const SettingsEnvironmentPage: FC = memo(
  function SettingsEnvironmentPage() {
    const dispatch = useAppDispatch();
    const envIds = useAppSelector(selectEnvIds);

    const handleDelete = useCallback(
      async (environmentId: string) => {
        dispatch(deleteEnvironment({ environmentId }))
          .then(unwrapResult)
          .then(() => {
            dispatch(
              addToast({
                message: DELETE_ENVIRONMENT_SUCCESS,
                severity: "success",
              })
            );
          })
          .catch(() => null);
      },
      [dispatch]
    );

    useEffect(() => {
      dispatch(fetchEnvironments());
    }, [dispatch]);

    return (
      <>
        <TableContainer component={Paper} square>
          <Table aria-label="environment list" size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell colSpan={2}>Description</TableCell>
                <TableCell>ID</TableCell>
                <TableCell align="right" />
              </TableRow>
            </TableHead>

            <TableBody>
              {envIds.map((envId) => (
                <EnvironmentListItem
                  id={envId}
                  key={`env-list-item-${envId}`}
                />
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        <DeleteEnvironmentDialog onDelete={handleDelete} />
      </>
    );
  }
);
