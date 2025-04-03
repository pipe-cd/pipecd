import { FC } from "react";
import { theme } from "~/theme";

type Props = {
  top: number;
  left: number;
  width: number;
  height: number;
  points: string;
};

const STROKE_WIDTH = 2;

const ResourceConnector: FC<Props> = ({
  top,
  left,
  width,
  height,
  points,
}: Props) => {
  return (
    <svg
      style={{
        position: "absolute",
        top,
        left,
      }}
      width={width}
      height={height}
    >
      <polyline
        points={points}
        strokeWidth={STROKE_WIDTH}
        stroke={theme.palette.divider}
        fill="transparent"
      />
    </svg>
  );
};

export default ResourceConnector;
