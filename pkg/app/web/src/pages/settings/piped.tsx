import {
  Box,
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  Divider,
  IconButton,
  makeStyles,
  Menu,
  MenuItem,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Toolbar,
  Typography,
} from "@material-ui/core";
import {
  Add as AddIcon,
  Close as CloseIcon,
  MoreVert as MoreVertIcon,
  FilterList as FilterIcon,
} from "@material-ui/icons";
import dayjs from "dayjs";
import React, { FC, memo, useCallback, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddPipedDrawer } from "../../components/add-piped-drawer";
import { EditPipedDrawer } from "../../components/edit-piped-drawer";
import { PipedFilter, FilterValues } from "../../components/piped-filter";
import { UI_TEXT_FILTER, UI_TEXT_HIDE_FILTER } from "../../constants/ui-text";
import { AppState } from "../../modules";
import {
  clearRegisteredPipedInfo,
  disablePiped,
  enablePiped,
  fetchPipeds,
  Piped,
  recreatePipedKey,
  RegisteredPiped,
  selectAll,
} from "../../modules/pipeds";
import { AppDispatch } from "../../store";

const useStyles = makeStyles((theme) => ({
  item: {
    backgroundColor: theme.palette.background.paper,
  },
  disabledPipedsAccordion: {
    padding: 0,
  },
  disabledItemsSummary: {
    borderBottom: "1px solid rgba(0, 0, 0, .125)",
  },
  pipedsList: {
    flex: 1,
  },
  disabledItemsSecondaryHeader: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(3),
  },
  disabledItem: {
    opacity: 0.6,
  },
  toolbarSpacer: {
    flexGrow: 1,
  },
}));

const ITEM_HEIGHT = 48;

const usePipeds = (filterValues: FilterValues): Piped[] => {
  const pipeds = useSelector<AppState, Piped[]>((state) =>
    selectAll(state.pipeds)
  );

  if (filterValues.enabled) {
    return pipeds.filter((piped) => piped.disabled === false);
  }

  if (filterValues.enabled === false) {
    return pipeds.filter((piped) => piped.disabled);
  }

  return pipeds;
};

export const SettingsPipedPage: FC = memo(function SettingsPipedPage() {
  const classes = useStyles();
  const [openFilter, setOpenFilter] = useState(false);
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [actionTarget, setActionTarget] = useState<Piped | null>(null);
  const [editPipedId, setEditPipedId] = useState<string | null>(null);
  const [filterValues, setFilterValues] = useState<FilterValues>({
    enabled: true,
  });
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const isOpenMenu = Boolean(anchorEl);
  const dispatch = useDispatch<AppDispatch>();
  const pipeds = usePipeds(filterValues);

  const registeredPiped = useSelector<AppState, RegisteredPiped | null>(
    (state) => state.pipeds.registeredPiped
  );

  const handleMenuOpen = useCallback(
    (event: React.MouseEvent<HTMLButtonElement>, piped: Piped): void => {
      setActionTarget(piped);
      setAnchorEl(event.currentTarget);
    },
    []
  );

  const closeMenu = useCallback(() => {
    setAnchorEl(null);
    setTimeout(() => {
      setActionTarget(null);
    }, 200);
  }, []);

  const handleDisableClick = useCallback(() => {
    closeMenu();
    if (!actionTarget) {
      return;
    }

    const act = actionTarget.disabled ? enablePiped : disablePiped;

    dispatch(act({ pipedId: actionTarget.id })).then(() => {
      dispatch(fetchPipeds(true));
    });
  }, [dispatch, actionTarget, closeMenu]);

  const handleClose = useCallback(() => {
    setIsOpenForm(false);
  }, []);

  const handleClosePipedInfo = useCallback(() => {
    dispatch(clearRegisteredPipedInfo());
    dispatch(fetchPipeds(true));
  }, [dispatch]);

  const handleRecreate = useCallback(() => {
    if (actionTarget) {
      dispatch(recreatePipedKey({ pipedId: actionTarget.id }));
    }
    closeMenu();
  }, [dispatch, actionTarget, closeMenu]);

  const handleEdit = useCallback(() => {
    if (actionTarget) {
      setEditPipedId(actionTarget.id);
    }
    closeMenu();
  }, [actionTarget, closeMenu]);

  const handleEditClose = useCallback(() => {
    setEditPipedId(null);
  }, []);

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
        <div className={classes.toolbarSpacer} />
        <Button
          color="primary"
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>
      <Divider />

      <Box display="flex" height="100%">
        <TableContainer component={Paper} square>
          <Table aria-label="piped list" size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Version</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Started At</TableCell>
                <TableCell align="right" />
              </TableRow>
            </TableHead>
            <TableBody>
              {pipeds.map((piped) => (
                <TableRow key={`pipe-${piped.id}`}>
                  <TableCell>
                    <Typography variant="subtitle2">
                      {`${piped.name} (${piped.id.slice(0, 8)})`}
                    </Typography>
                  </TableCell>
                  <TableCell>{piped.version}</TableCell>
                  <TableCell>
                    <Typography variant="body2" color="textSecondary">
                      {piped.desc}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    {piped.startedAt === 0
                      ? "Not Yet Started"
                      : dayjs(piped.startedAt * 1000).fromNow()}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      edge="end"
                      aria-label="open menu"
                      onClick={(e) => handleMenuOpen(e, piped)}
                    >
                      <MoreVertIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {openFilter && (
          <PipedFilter values={filterValues} onChange={setFilterValues} />
        )}
      </Box>

      <Menu
        id="piped-menu"
        anchorEl={anchorEl}
        keepMounted
        open={isOpenMenu}
        onClose={() => closeMenu()}
        PaperProps={{
          style: {
            maxHeight: ITEM_HEIGHT * 4.5,
            width: "20ch",
          },
        }}
      >
        {actionTarget && actionTarget.disabled ? (
          <MenuItem onClick={handleDisableClick}>Enable</MenuItem>
        ) : (
          [
            <MenuItem key="piped-menu-edit" onClick={handleEdit}>
              Edit
            </MenuItem>,
            <MenuItem key="piped-menu-recreate" onClick={handleRecreate}>
              Recreate Key
            </MenuItem>,
            <MenuItem key="piped-menu-disable" onClick={handleDisableClick}>
              Disable
            </MenuItem>,
          ]
        )}
      </Menu>

      <AddPipedDrawer open={isOpenForm} onClose={handleClose} />
      <EditPipedDrawer pipedId={editPipedId} onClose={handleEditClose} />

      <Dialog open={Boolean(registeredPiped)}>
        <DialogTitle>Piped registered</DialogTitle>
        <DialogContent>
          <TextField
            label="id"
            variant="outlined"
            value={registeredPiped?.id || ""}
            fullWidth
            margin="dense"
          />
          <TextField
            label="secret key"
            variant="outlined"
            value={registeredPiped?.key || ""}
            fullWidth
            margin="dense"
          />
          <Box display="flex" justifyContent="flex-end" m={1} mt={2}>
            <Button color="primary" onClick={handleClosePipedInfo}>
              CLOSE
            </Button>
          </Box>
        </DialogContent>
      </Dialog>
    </>
  );
});
