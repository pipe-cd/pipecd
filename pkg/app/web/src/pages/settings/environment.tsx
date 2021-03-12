import {
  Box,
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
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddEnvForm } from "../../components/add-env-form";
import { EnvironmentListItem } from "../../components/environment-list-item";
import { UI_TEXT_ADD } from "../../constants/ui-text";
import {
  addEnvironment,
  fetchEnvironments,
  selectEnvIds,
} from "../../modules/environments";
import { selectProjectName } from "../../modules/me";
import { AppDispatch } from "../../store";

export const SettingsEnvironmentPage: FC = memo(
  function SettingsEnvironmentPage() {
    const dispatch = useDispatch<AppDispatch>();
    const [isOpenForm, setIsOpenForm] = useState(false);
    const projectName = useSelector(selectProjectName);
    const envIds = useSelector(selectEnvIds);

    const handleClose = (): void => {
      setIsOpenForm(false);
    };

    const handleSubmit = (props: { name: string; desc: string }): void => {
      dispatch(addEnvironment(props)).finally(() => {
        setIsOpenForm(false);
        dispatch(fetchEnvironments());
      });
    };
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

        <Box display="flex" height="100%">
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
        </Box>

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
