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

export const rightArea = style({
  flexGrow: 1,
  maxWidth: rem(vars.breakpoints.sm),
});

export const mainEmbedEditorArea = style({
  width: "100%",
  display: "flex",
});

export const titleAndMainContentArea = style({
  display: "flex",
});

export const titleAndMainText = style({
  flexGrow: 1,
});

export const mainImageWrapper = style({
  overflow: "clip",
  borderRadius: rem(vars.radius.md),
  cursor: "pointer",
  position: "relative",
  border: `2px dashed ${vars.colors.gray[4]}`,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
});

export const mainImage = style({
  position: "absolute",
  transition: "filter linear 1s",
  filter: "brightness(0.6)",
});

export const uploadBoxWrapper = style({
  aspectRatio: "16 / 9",
  zIndex: 1,
});
