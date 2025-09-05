import { Box, Button, Divider, Drawer, Toolbar } from "@mui/material";
import { Add } from "@mui/icons-material";
import CloseIcon from "@mui/icons-material/Close";
import FilterIcon from "@mui/icons-material/FilterList";
import RefreshIcon from "@mui/icons-material/Refresh";
import LockOutlineIcon from "@mui/icons-material/LockOutline";
import { FC, useCallback, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import {
  UI_TEXT_ADD,
  UI_ENCRYPT_SECRET,
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_REFRESH,
} from "~/constants/ui-text";
import {
  arrayFormat,
  stringifySearchParams,
  useSearchParams,
} from "~/utils/search-params";
import AddApplicationDrawer from "./add-application-drawer";
import { ApplicationAddedView } from "./application-added-view";
import { ApplicationFilter } from "./application-filter";
import { ApplicationList } from "./application-list";
import EncryptSecretDrawer from "./encrypt-secret-drawer";
import { SpinnerIcon } from "~/styles/button";
import {
  ApplicationsFilterOptions,
  useGetApplications,
} from "~/queries/applications/use-get-applications";
import { getTypedValue, isString, isStringArray } from "~/utils/common";

export const ApplicationIndexPage: FC = () => {
  const navigate = useNavigate();
  const filterOptions = useSearchParams();
  const [openAddForm, setOpenAddForm] = useState(false);
  const [openFilter, setOpenFilter] = useState(true);
  const [openEncryptSecretDrawer, setOpenEncryptSecretDrawer] = useState(false);
  const [showCongratulation, setShowCongratulation] = useState(false);

  const searchValues: ApplicationsFilterOptions = useMemo(() => {
    return {
      ...filterOptions,
      page: undefined,
      activeStatus: getTypedValue(filterOptions, "activeStatus", isString),
      kind: getTypedValue(filterOptions, "kind", isString),
      syncStatus: getTypedValue(filterOptions, "syncStatus", isString),
      name: getTypedValue(filterOptions, "name", isString),
      pipedId: getTypedValue(filterOptions, "pipedId", isString),
      labels: getTypedValue(filterOptions, "labels", isStringArray),
    };
  }, [filterOptions]);

  const { data: applications, isLoading, refetch } = useGetApplications(
    searchValues
  );

  const currentPage =
    typeof filterOptions.page === "string"
      ? parseInt(filterOptions.page, 10)
      : 1;

  const updateURL = useCallback(
    (options: Record<string, string | number | boolean | undefined>) => {
      navigate(
        `${PAGE_PATH_APPLICATIONS}?${stringifySearchParams(
          { ...options },
          { arrayFormat: arrayFormat }
        )}`,
        { replace: true }
      );
    },
    [navigate]
  );

  const handleFilterChange = useCallback(
    (options) => {
      updateURL({ ...options, page: 1 });
    },
    [updateURL]
  );
  const handleFilterClear = useCallback(() => {
    updateURL({ page: currentPage });
  }, [updateURL, currentPage]);

  const handleCloseApplicationAddedView = (): void => {
    setShowCongratulation(false);
  };

  const handlePageChange = useCallback(
    (page: number) => {
      updateURL({ ...filterOptions, page });
    },
    [updateURL, filterOptions]
  );

  return (
    <>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<Add />}
          onClick={() => setOpenAddForm(true)}
        >
          {UI_TEXT_ADD}
        </Button>
        <Button
          color="primary"
          startIcon={<LockOutlineIcon />}
          onClick={() => setOpenEncryptSecretDrawer(true)}
          sx={{ ml: 1 }}
        >
          {UI_ENCRYPT_SECRET}
        </Button>
        <Box
          sx={{
            flex: 1,
          }}
        />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={() => refetch()}
          sx={{ position: "relative" }}
          disabled={isLoading}
        >
          {UI_TEXT_REFRESH}
          {isLoading && <SpinnerIcon />}
        </Button>
        <Button
          color="primary"
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>
      <Divider />
      <Box
        sx={{
          display: "flex",
          overflowY: "hidden",
          overflowX: "auto",
          flex: 1,
        }}
      >
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            flex: 1,
            p: 2,
          }}
        >
          <ApplicationList
            applications={applications || []}
            currentPage={currentPage}
            onPageChange={handlePageChange}
          />
        </Box>
        {openFilter && (
          <ApplicationFilter
            applications={applications || []}
            options={filterOptions}
            onChange={handleFilterChange}
            onClear={handleFilterClear}
          />
        )}
      </Box>
      <AddApplicationDrawer
        open={openAddForm}
        onClose={() => setOpenAddForm(false)}
        onAdded={() => {
          setOpenAddForm(false);
          setShowCongratulation(true);
        }}
      />

      <Drawer
        anchor="right"
        open={!!showCongratulation}
        onClose={(_, reason) => {
          if (reason === "backdropClick") return;
          handleCloseApplicationAddedView();
        }}
      >
        <ApplicationAddedView onClose={handleCloseApplicationAddedView} />
      </Drawer>
      <EncryptSecretDrawer
        open={openEncryptSecretDrawer}
        onClose={() => setOpenEncryptSecretDrawer(false)}
      />
    </>
  );
};
