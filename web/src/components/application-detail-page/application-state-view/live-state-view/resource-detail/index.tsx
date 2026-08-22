import { Box } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { FC } from "react";
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
        <InfoRowTitle>Resource Type</InfoRowTitle>
        <InfoRowValue>{resource.resourceType}</InfoRowValue>
      </Box>
      {resource.resourceMetadataMap.map(([key, value]) => (
        <Box
          key={key}
          sx={{
            pt: 1,
            display: "flex",
            alignItems: "center",
          }}
        >
          <InfoRowTitle>{key}</InfoRowTitle>
          <InfoRowValue>{value || "Empty"}</InfoRowValue>
        </Box>
      ))}
      {resource.healthDescription && (
        <Box
          sx={{
            pt: 1,
          }}
        >
          <InfoRowTitle>Health Description</InfoRowTitle>
          <InfoRowValue>{resource.healthDescription}</InfoRowValue>
        </Box>
      )}
    </PanelWrap>
  );
};
