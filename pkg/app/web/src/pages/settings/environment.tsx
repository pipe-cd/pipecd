import {
  Button,
  Divider,
  Drawer,
  List,
  makeStyles,
  Toolbar,
} from "@material-ui/core";
import { Add as AddIcon } from "@material-ui/icons";
import { EntityId } from "@reduxjs/toolkit";
import React, { FC, memo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddEnvForm } from "../../components/add-env-form";
import { EnvironmentListItem } from "../../components/environment-list-item";
import { UI_TEXT_ADD } from "../../constants/ui-text";
import { AppState } from "../../modules";
import {
  addEnvironment,
  fetchEnvironments,
  selectIds as selectEnvIds,
} from "../../modules/environments";
import { selectProjectName } from "../../modules/me";
import { AppDispatch } from "../../store";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
  },
  listItem: {
    backgroundColor: theme.palette.background.paper,
  },
}));

export const SettingsEnvironmentPage: FC = memo(
  function SettingsEnvironmentPage() {
    const classes = useStyles();
    const dispatch = useDispatch<AppDispatch>();
    const [isOpenForm, setIsOpenForm] = useState(false);
    const projectName = useSelector<AppState, string>((state) =>
      selectProjectName(state.me)
    );
    const envIds = useSelector<AppState, EntityId[]>((state) =>
      selectEnvIds(state.environments)
    );

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

        <div className={classes.main}>
          <List disablePadding>
            {envIds.map((envId) => (
              <EnvironmentListItem id={envId} key={`env-list-item-${envId}`} />
            ))}
          </List>
        </div>

        <Drawer anchor="right" open={isOpenForm} onClose={handleClose}>
          <AddEnvForm
            projectName={projectName}
            onClose={handleClose}
            onSubmit={handleSubmit}
          />
        </Drawer>
      </>
    );
  }
);
