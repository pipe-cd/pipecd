import {
  Menu,
  MenuItem,
  Paper,
  Table,
  IconButton,
  TableBody,
  TableCell,
  makeStyles,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@material-ui/core";
import { useProjectSettingStyles } from "~/styles/project-setting";
import { MoreVert as MenuIcon } from "@material-ui/icons";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchProject } from "~/modules/project";
import * as React from "react";
import { FC, memo, useCallback, useEffect } from "react";
import {
  RBACPolicy,
  RBAC_RESOURCE_TYPE_TEXT,
  RBAC_ACTION_TYPE_TEXT,
} from "~/modules/project";

const useStyles = makeStyles((theme) => ({
  title: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    paddingTop: theme.spacing(2),
  },
}));

const SUB_SECTION_TITLE = "Role";
const RESOURCE_ACTION_SEPARATOR = ";";
const RESOURCES_KEY = "resources";
const ACTIONS_KEY = "actions";

export const RoleTable: FC = memo(function RoleTable() {
  const classes = useStyles();
  const projectSettingClasses = useProjectSettingStyles();
  const rbacRoles = useAppSelector((state) => state.project.rbacRoles);
  const dispatch = useAppDispatch();
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );

  const formalizePoliciesFromPoliciesList = ({
    policiesList,
  }: {
    policiesList: Array<RBACPolicy>;
  }): Array<string> => {
    const policies: Array<string> = [];
    policiesList.map((policy) => {
      const resources: Array<string> = [];
      policy.resourcesList.map((resource) => {
        resources.push(RBAC_RESOURCE_TYPE_TEXT[resource.type]);
      });

      const actions: Array<string> = [];
      policy.actionsList.map((action) => {
        actions.push(RBAC_ACTION_TYPE_TEXT[action]);
      });

      const resource = RESOURCES_KEY + "=" + resources.join(",");
      const action = ACTIONS_KEY + "=" + actions.join(",");
      policies.push(resource + RESOURCE_ACTION_SEPARATOR + action);
    });

    return policies;
  };

  const handleOpenMenu = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      setAnchorEl(e.currentTarget);
    },
    [setAnchorEl]
  );

  const handleCloseMenu = useCallback(() => {
    setAnchorEl(null);
  }, [setAnchorEl]);

  useEffect(() => {
    dispatch(fetchProject());
  }, [dispatch]);

  return (
    <>
      <div className={classes.title}>
        <Typography
          variant="h6"
          className={projectSettingClasses.titleWithIcon}
        >
          {SUB_SECTION_TITLE}
        </Typography>
      </div>

      <TableContainer component={Paper} square>
        <Table size="small" stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell>Role</TableCell>
              <TableCell>Policies</TableCell>
              <TableCell align="right" />
            </TableRow>
          </TableHead>
          <TableBody>
            {rbacRoles.map((role) => (
              <TableRow key={role.name}>
                <TableCell>{role.name}</TableCell>
                <TableCell>
                  {formalizePoliciesFromPoliciesList({
                    policiesList: role.policiesList,
                  }).map((policy, i) => (
                    <p key={i}>{policy}</p>
                  ))}
                </TableCell>
                <TableCell align="right">
                  <IconButton
                    data-id={role.name}
                    onClick={handleOpenMenu}
                    disabled={role.isBuiltin}
                  >
                    <MenuIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Menu
        id="role-menu"
        open={Boolean(anchorEl)}
        anchorEl={anchorEl}
        onClose={handleCloseMenu}
      >
        <MenuItem>Edit</MenuItem>
        <MenuItem>Delete</MenuItem>
      </Menu>
    </>
  );
});
