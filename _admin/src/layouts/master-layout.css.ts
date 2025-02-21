import { rem } from "@mantine/core";
import { style } from "@vanilla-extract/css";
import { vars } from "~/theme";

export const menuItem = style({
	padding: `${rem(vars.spacing.sm)} ${rem(vars.spacing.md)}`,
});

export const header = style({
	height: "100%",
	paddingLeft: rem(vars.spacing.md),
});
