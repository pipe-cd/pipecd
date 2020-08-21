import React, { FC, useEffect, useState } from "react";
import {
  makeStyles,
  Typography,
  Divider,
  TextField,
  MenuItem,
  Link,
  Button,
} from "@material-ui/core";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchTemplateList,
  DeploymentConfigTemplateLabel,
  DeploymentConfigTemplateLabelKey,
  DeploymentConfigTemplate,
  selectTemplateByAppId,
} from "../modules/deployment-configs";
import { AppState } from "../modules";

const useStyles = makeStyles((theme) => ({
  root: {
    width: 600,
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
}));

const TEXT = {
  TEMPLATE_FILENAME: ".pipe.yaml",
  CREATE_LINK: "Create template file",
};

interface Props {
  applicationId: string;
  onSkip: () => void;
}

export const DeploymentConfigForm: FC<Props> = ({ applicationId, onSkip }) => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [label, setLabel] = useState<DeploymentConfigTemplateLabel | undefined>(
    undefined
  );
  const template = useSelector<
    AppState,
    DeploymentConfigTemplate.AsObject | null
  >((state) => selectTemplateByAppId(state.deploymentConfigs));

  useEffect(() => {
    if (label !== undefined) {
      dispatch(fetchTemplateList({ labels: [label], applicationId }));
    }
  }, [dispatch, label, applicationId]);

  return (
    <div className={classes.root}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Generate deployment configuration`}</Typography>
      <Divider />
      <div className={classes.content}>
        <TextField
          fullWidth
          required
          select
          // disabled={disabled}
          label="Deployment Strategy"
          variant="outlined"
          margin="dense"
          onChange={(e) =>
            setLabel(
              (e.target.value as unknown) as DeploymentConfigTemplateLabel
            )
          }
          value={label}
          style={{ flex: 1 }}
        >
          {Object.keys(DeploymentConfigTemplateLabel).map((key) => (
            <MenuItem
              key={key}
              value={
                DeploymentConfigTemplateLabel[
                  key as DeploymentConfigTemplateLabelKey
                ]
              }
            >
              {key}
            </MenuItem>
          ))}
        </TextField>

        <Typography variant="subtitle1" className={classes.filename}>
          {TEXT.TEMPLATE_FILENAME}
        </Typography>
        <TextField
          multiline
          fullWidth
          variant="outlined"
          margin="dense"
          rows={30}
          rowsMax={30}
          value={
            template ? template.content : "Select a deployment strategy first"
          }
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
    </div>
  );
};
