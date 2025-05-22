import { Box } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { FC } from "react";
import { findMetadataByKey } from "~/utils/find-metadata-by-key";
import { ResourceState } from "~~/model/application_live_state_pb";
import {
  CloseButton,
  InfoRowTitle,
  InfoRowValue,
  PanelTitle,
  PanelWrap,
} from "./styles";

export interface ResourceDetailProps {
  resource: ResourceState.AsObject;
  onClose: () => void;
}

export const ResourceDetail: FC<ResourceDetailProps> = ({
  resource,
  onClose,
}) => {
  return (
    <PanelWrap square>
      <CloseButton onClick={onClose} size="large">
        <CloseIcon />
      </CloseButton>
      <PanelTitle>{resource.name}</PanelTitle>
      <Box
        sx={{
          pt: 1,
          display: "flex",
          alignItems: "center",
        }}
      >
        <InfoRowTitle>Kind</InfoRowTitle>
        <InfoRowValue>
          {findMetadataByKey(resource.resourceMetadataMap, "Kind")}
        </InfoRowValue>
      </Box>
      <Box
        sx={{
          pt: 1,
          display: "flex",
          alignItems: "center",
        }}
      >
        <InfoRowTitle>Namespace</InfoRowTitle>
        <InfoRowValue>
          {findMetadataByKey(resource.resourceMetadataMap, "Namespace")}
        </InfoRowValue>
      </Box>
      <Box
        sx={{
          pt: 1,
          display: "flex",
          alignItems: "center",
        }}
      >
        <InfoRowTitle>Api Version</InfoRowTitle>
        <InfoRowValue>
          {findMetadataByKey(resource.resourceMetadataMap, "API Version")}
        </InfoRowValue>
      </Box>
      <Box
        sx={{
          pt: 1,
        }}
      >
        <InfoRowTitle>Health Description</InfoRowTitle>
        <InfoRowValue>{resource.healthDescription || "Empty"}</InfoRowValue>
      </Box>
    </PanelWrap>
  );
};
