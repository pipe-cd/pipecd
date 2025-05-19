import { Box } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { FC } from "react";
import {
  CloseButton,
  InfoRowTitle,
  InfoRowValue,
  PanelTitle,
  PanelWrap,
} from "./styles";

export interface LambdaResourceDetailProps {
  resource: {
    name: string;
    kind: string;
    healthDescription: string;
  };
  onClose: () => void;
}

export const LambdaResourceDetail: FC<LambdaResourceDetailProps> = ({
  resource,
  onClose,
}) => {
  return (
    <PanelWrap square>
      <CloseButton onClick={onClose} size="large">
        <CloseIcon />
      </CloseButton>
      <PanelTitle>{resource.name}</PanelTitle>
      <Box pt={1} display={"flex"} alignItems={"center"}>
        <InfoRowTitle>Kind</InfoRowTitle>
        <InfoRowValue>{resource.kind}</InfoRowValue>
      </Box>
      <Box pt={1}>
        <InfoRowTitle>Health Description</InfoRowTitle>
        <InfoRowValue>{resource.healthDescription || "Empty"}</InfoRowValue>
      </Box>
    </PanelWrap>
  );
};
