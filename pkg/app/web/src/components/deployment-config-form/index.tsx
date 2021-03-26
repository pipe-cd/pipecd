import {
  Box,
  Button,
  Divider,
  IconButton,
  Link,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
} from "@material-ui/core";
import CopyIcon from "@material-ui/icons/FileCopyOutlined";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import copy from "copy-to-clipboard";
import React, { FC, useEffect, useState, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../../modules";
import {
  DeploymentConfigTemplate,
  fetchTemplateList,
  selectTemplatesByAppId,
} from "../../modules/deployment-configs";
import { addToast } from "../../modules/toasts";

const useStyles = makeStyles((theme) => ({
  title: {
    padding: theme.spacing(2),
  },
  filename: {
    marginTop: theme.spacing(2),
    fontFamily: theme.typography.fontFamilyMono,
    color: theme.palette.text.secondary,
  },
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
  },
}));

const TEXT = {
  TITLE: "Add the deployment configuration file",
  PLACEHOLDER:
    "# Fill the deployment configuration here. You can also choose one of the provided templates above to edit.",
  CONFIGURATION_FILENAME: ".pipe.yaml",
  CREATE_LINK: "Add to application configuration directory in Git",
};

export interface DeploymentConfigFormProps {
  onSkip: () => void;
}

export const DeploymentConfigForm: FC<DeploymentConfigFormProps> = memo(
  function DeploymentConfigForm({ onSkip }) {
    const classes = useStyles();
    const dispatch = useDispatch();
    const [templateIndex, setTemplateIndex] = useState(0);
    const [configValue, setConfigValue] = useState(TEXT.PLACEHOLDER);
    const templates = useSelector<
      AppState,
      DeploymentConfigTemplate.AsObject[]
    >((state) => selectTemplatesByAppId(state.deploymentConfigs) || []);

    const template = templates[templateIndex];

    const handleOnClickCopy = (): void => {
      copy(configValue);
      dispatch(addToast({ message: "Deployment config copied to clipboard" }));
    };

    const handleTemplateChange = (
      e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>
    ): void => {
      setTemplateIndex(parseInt(e.target.value, 10));
    };

    useEffect(() => {
      dispatch(fetchTemplateList({ labels: [] }));
    }, [dispatch]);

    useEffect(() => {
      if (template) {
        setConfigValue(template.content);
      }
    }, [setConfigValue, template]);

    return (
      <Box width={600} flex={1} display="flex" flexDirection="column">
        <Typography className={classes.title} variant="h6">
          {TEXT.TITLE}
        </Typography>

        <Divider />

        <Box p={2}>
          {templates.length === 0 ? null : (
            <TextField
              fullWidth
              required
              select
              label="Template"
              variant="outlined"
              margin="dense"
              onChange={handleTemplateChange}
              value={templateIndex}
              style={{ flex: 1 }}
              disabled={templates.length === 0}
            >
              {templates.map(({ name }, index) => (
                <MenuItem key={name} value={index}>
                  {name}
                </MenuItem>
              ))}
            </TextField>
          )}

          <Box
            display="flex"
            alignItems="flex-end"
            justifyContent="space-between"
          >
            <Typography variant="subtitle1" className={classes.filename}>
              {TEXT.CONFIGURATION_FILENAME}
            </Typography>
            <IconButton
              size="small"
              aria-label="Copy deployment config"
              onClick={handleOnClickCopy}
            >
              <CopyIcon fontSize="small" />
            </IconButton>
          </Box>
          <TextField
            multiline
            fullWidth
            variant="outlined"
            margin="dense"
            rows={30}
            rowsMax={30}
            value={configValue}
            onChange={(e) => setConfigValue(e.target.value)}
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

          <Box mt={1} textAlign="right">
            <Button onClick={onSkip} variant="outlined">
              SKIP
            </Button>
          </Box>
        </Box>
      </Box>
    );
  }
);
