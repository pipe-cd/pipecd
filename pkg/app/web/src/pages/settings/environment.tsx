import {
  makeStyles,
  Toolbar,
  Divider,
  Button,
  List,
  ListItemText,
  ListItem,
  Drawer,
} from "@material-ui/core";
import React, { FC, memo, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { AppState } from "../../modules";
import {
  selectAll as selectEnvsAll,
  Environment,
  addEnvironment,
  fetchEnvironments,
} from "../../modules/environments";
import AddIcon from "@material-ui/icons/Add";
import { AddEnvForm } from "../../components/add-env-form";
import { AppDispatch } from "../../store";
import { selectProjectName } from "../../modules/me";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
  },
  listItem: {
    backgroundColor: theme.palette.background.paper,
  },
}));

const TEXT = {
  NO_DESCRIPTION: "No description",
};

export const SettingsEnvironmentPage: FC = memo(
  function SettingsEnvironmentPage() {
    const classes = useStyles();
    const dispatch = useDispatch<AppDispatch>();
    const [isOpenForm, setIsOpenForm] = useState(false);
    const projectName = useSelector<AppState, string>((state) =>
      selectProjectName(state.me)
    );
    const envs = useSelector<AppState, Environment[]>((state) =>
      selectEnvsAll(state.environments)
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
            ADD
          </Button>
        </Toolbar>
        <Divider />

        <div className={classes.main}>
          <List disablePadding>
            {envs.map((env) => (
              <ListItem
                key={`env-${env.id}`}
                divider
                dense
                className={classes.listItem}
              >
                <ListItemText
                  primary={env.name}
                  secondary={env.desc || TEXT.NO_DESCRIPTION}
                />
              </ListItem>
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
