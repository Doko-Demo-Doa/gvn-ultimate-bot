import "@mantine/core/styles.css";
import React from "react";
import { MantineProvider, ColorSchemeScript } from "@mantine/core";
import { theme } from "../../theme";
import { Metadata } from "next";
import Providers from "~/providers/master-provider";

export const metadata: Metadata = {
  title: "DautoVN Admin",
  description: "Dashboard for DautoVN",
};

export default function RootLayout({ children }: { children: any }) {
  return (
    <html lang="en">
      <head>
        <ColorSchemeScript />
        <link rel="shortcut icon" href="/favicon.svg" />
        <meta
          name="viewport"
          content="minimum-scale=1, initial-scale=1, width=device-width, user-scalable=no"
        />
      </head>
      <body>
        <Providers>
          <MantineProvider theme={theme}>{children}</MantineProvider>
        </Providers>
      </body>
    </html>
  );
}
