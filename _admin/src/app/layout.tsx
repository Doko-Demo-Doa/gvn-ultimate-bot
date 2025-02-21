import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import "@mantine/dropzone/styles.css";
import { ClerkProvider } from "@clerk/nextjs";
import {
  ColorSchemeScript,
  MantineProvider,
  mantineHtmlProps,
} from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import type { Metadata } from "next";
import type React from "react";
import { customApiClient } from "~/datasource/rest/api-client";
import Providers from "~/providers/master-provider";

import { theme } from "../theme";

export const metadata: Metadata = {
  title: "DautoVN Admin",
  description: "Dashboard for DautoVN",
};

customApiClient.init({
  baseUrl: process.env.NEXT_PUBLIC_BASE_API_URL || "",
});

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <ClerkProvider>
      <html lang="en" {...mantineHtmlProps}>
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
            <MantineProvider defaultColorScheme="dark" theme={theme}>
              {children}
              <Notifications />
            </MantineProvider>
          </Providers>
        </body>
      </html>
    </ClerkProvider>
  );
}
