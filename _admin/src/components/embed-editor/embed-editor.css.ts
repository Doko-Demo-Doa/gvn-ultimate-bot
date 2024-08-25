import { rem } from "@mantine/core";
import { style } from "@vanilla-extract/css";
import { vars } from "~/theme";

export const embedPreviewContainer = style({
  padding: `${rem(vars.spacing.sm)} ${rem(vars.spacing.md)}`,
});

export const groupWrapper = style({
  borderLeftWidth: `2px solid`,
  borderLeftStyle: "solid",
  paddingLeft: rem(vars.spacing.md),
});

export const segmented = style({
  maxWidth: rem(400),
});
