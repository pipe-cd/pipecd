import {
  Button,
  CircularProgress,
  Divider,
  Link,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
} from "@material-ui/core";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import React, { FC, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  DeploymentConfigTemplate,
  fetchTemplateList,
  selectTemplatesByAppId,
} from "../modules/deployment-configs";

const useStyles = makeStyles((theme) => ({
  root: {
    width: 600,
    flex: 1,
    display: "flex",
    flexDirection: "column",
  },
  title: {
    padding: theme.spacing(2),
  },
  content: {
    padding: theme.spacing(2),
  },
  filename: {
    marginTop: theme.spacing(2),
    fontFamily: "Roboto Mono",
    color: theme.palette.text.secondary,
  },
  templateContent: {
    fontFamily: "Roboto Mono",
    fontSize: 14,
  },
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
  },
  actions: {
    marginTop: theme.spacing(1),
    textAlign: "right",
  },
  loading: {
    flex: 1,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
}));

const TEXT = {
  TITLE: "Add the deployment configuration file",
  PLACEHOLDER:
    "# Fill the deployment configuration here. You can also choose one of the provided templates above to edit.",
  CONFIGURATION_FILENAME: ".pipe.yaml",
  CREATE_LINK:
    "Add this deployment configuration file to application configuration directory in Git",
};

interface Props {
  applicationId: string;
  onSkip: () => void;
}

export const DeploymentConfigForm: FC<Props> = ({ applicationId, onSkip }) => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [templateIndex, setTemplateIndex] = useState(0);
  const templates = useSelector<
    AppState,
    DeploymentConfigTemplate.AsObject[] | null
  >((state) => selectTemplatesByAppId(state.deploymentConfigs));

  const template = templates && templates[templateIndex];

  useEffect(() => {
    dispatch(fetchTemplateList({ labels: [], applicationId }));
  }, [dispatch, applicationId]);

  return (
    <div className={classes.root}>
      <Typography className={classes.title} variant="h6">
        {TEXT.TITLE}
      </Typography>
      <Divider />

      {templates === null ? (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      ) : (
        <div className={classes.content}>
          <TextField
            fullWidth
            required
            select
            label="Template"
            variant="outlined"
            margin="dense"
            onChange={(e) => setTemplateIndex(parseInt(e.target.value, 10))}
            value={templateIndex}
            style={{ flex: 1 }}
          >
            {templates.map(({ name }, index) => (
              <MenuItem key={name} value={index}>
                {name}
              </MenuItem>
            ))}
          </TextField>

          <Typography variant="subtitle1" className={classes.filename}>
            {TEXT.CONFIGURATION_FILENAME}
          </Typography>
          <TextField
            multiline
            fullWidth
            variant="outlined"
            margin="dense"
            rows={30}
            rowsMax={30}
            value={template ? template.content : TEXT.PLACEHOLDER}
            InputProps={{
              className: classes.templateContent,
              margin: "dense",
            }}
          />

          {template && (
            <Link
              href={template.fileCreationUrl}
              target="_blank"
              rel="noreferrer"
            >
              {TEXT.CREATE_LINK}
              <OpenInNewIcon className={classes.linkIcon} />
            </Link>
          )}

          <div className={classes.actions}>
            <Button onClick={onSkip} variant="outlined">
              SKIP
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};
