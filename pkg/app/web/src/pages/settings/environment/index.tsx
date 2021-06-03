import {
  Button,
  Divider,
  Drawer,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
} from "@material-ui/core";
import { Add as AddIcon } from "@material-ui/icons";
import { unwrapResult } from "@reduxjs/toolkit";
import { useEffect } from "react";
import { FC, memo, useCallback, useState } from "react";
import { AddEnvForm } from "../../../components/add-env-form";
import { EnvironmentListItem } from "../../../components/environment-list-item";
import { DELETE_ENVIRONMENT_SUCCESS } from "../../../constants/toast-text";
import { UI_TEXT_ADD } from "../../../constants/ui-text";
import { useAppDispatch, useAppSelector } from "../../../hooks/redux";
import {
  addEnvironment,
  deleteEnvironment,
  fetchEnvironments,
  selectEnvIds,
} from "../../../modules/environments";
import { selectProjectName } from "../../../modules/me";
import { addToast } from "../../../modules/toasts";
import { DeleteEnvironmentDialog } from "./components/delete-confirm-dialog";

export const SettingsEnvironmentPage: FC = memo(
  function SettingsEnvironmentPage() {
    const dispatch = useAppDispatch();
    const [isOpenForm, setIsOpenForm] = useState(false);
    const projectName = useAppSelector(selectProjectName);
    const envIds = useAppSelector(selectEnvIds);

    const handleClose = (): void => {
      setIsOpenForm(false);
    };

    const handleSubmit = (props: { name: string; desc: string }): void => {
      dispatch(addEnvironment(props)).finally(() => {
        setIsOpenForm(false);
        dispatch(fetchEnvironments());
      });
    };

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
        <Toolbar variant="dense">
          <Button
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => setIsOpenForm(true)}
          >
            {UI_TEXT_ADD}
          </Button>
        </Toolbar>
        <Divider />

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

        <Drawer anchor="right" open={isOpenForm} onClose={handleClose}>
          <AddEnvForm
            projectName={projectName}
            onCancel={handleClose}
            onSubmit={handleSubmit}
          />
        </Drawer>
      </>
    );
  }
);
