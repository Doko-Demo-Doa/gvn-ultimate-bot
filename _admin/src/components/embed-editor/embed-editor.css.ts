import { rem } from "@mantine/core";
import { style } from "@vanilla-extract/css";
import { vars } from "~/theme";

export const embedPreviewContainer = style({
  padding: `${rem(vars.spacing.sm)} ${rem(vars.spacing.md)}`,
});

export const typeIndicator = style({
  width: "2px",
  height: "400px",
  backgroundColor: vars.colors.red[6],
});

export const segmented = style({
  maxWidth: rem(400),
});
